use chrono::Utc;
use sqlx::PgPool;
use uuid::Uuid;

use crate::models::Quote;

/// Репозиторий для работы с цитатами
pub struct QuoteRepository {
    pool: PgPool,
}

impl QuoteRepository {
    pub fn new(pool: PgPool) -> Self {
        Self { pool }
    }

    /// Возвращает случайную цитату
    pub async fn get_random(&self) -> anyhow::Result<Quote> {
        use tracing::warn;
        
        // Сначала проверяем, есть ли цитаты вообще
        let count: i64 = sqlx::query_scalar("SELECT COUNT(*) FROM quotes")
            .fetch_one(&self.pool)
            .await?;
        
        if count == 0 {
            warn!("No quotes found in database");
            return Err(anyhow::anyhow!("no quotes found"));
        }
        
        let quote = sqlx::query_as::<_, Quote>(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             ORDER BY RANDOM() 
             LIMIT 1"
        )
        .fetch_optional(&self.pool)
        .await?;

        quote.ok_or_else(|| {
            warn!("Query returned no rows despite count being {}", count);
            anyhow::anyhow!("no quotes found")
        })
    }

    /// Возвращает все цитаты с пагинацией и поиском
    pub async fn get_all(
        &self,
        page: i32,
        page_size: i32,
        search: Option<&str>,
    ) -> anyhow::Result<(Vec<Quote>, i64)> {
        let offset = (page - 1) * page_size;

        // Подсчет общего количества
        let total = if let Some(search_term) = search {
            let search_pattern = format!("%{}%", search_term);
            sqlx::query_scalar::<_, i64>(
                "SELECT COUNT(*) FROM quotes WHERE text ILIKE $1 OR author ILIKE $1"
            )
            .bind(&search_pattern)
            .fetch_one(&self.pool)
            .await?
        } else {
            sqlx::query_scalar::<_, i64>("SELECT COUNT(*) FROM quotes")
                .fetch_one(&self.pool)
                .await?
        };

        // Получение цитат
        let quotes = if let Some(search_term) = search {
            let search_pattern = format!("%{}%", search_term);
            sqlx::query_as::<_, Quote>(
                "SELECT id, text, author, likes_count, created_at, updated_at 
                 FROM quotes 
                 WHERE text ILIKE $1 OR author ILIKE $1 
                 ORDER BY created_at DESC 
                 LIMIT $2 OFFSET $3"
            )
            .bind(&search_pattern)
            .bind(page_size)
            .bind(offset)
            .fetch_all(&self.pool)
            .await?
        } else {
            sqlx::query_as::<_, Quote>(
                "SELECT id, text, author, likes_count, created_at, updated_at 
                 FROM quotes 
                 ORDER BY created_at DESC 
                 LIMIT $1 OFFSET $2"
            )
            .bind(page_size)
            .bind(offset)
            .fetch_all(&self.pool)
            .await?
        };

        Ok((quotes, total))
    }

    /// Возвращает цитату по ID
    pub async fn get_by_id(&self, id: &str) -> anyhow::Result<Quote> {
        let quote = sqlx::query_as::<_, Quote>(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             WHERE id = $1"
        )
        .bind(id)
        .fetch_optional(&self.pool)
        .await?;

        quote.ok_or_else(|| anyhow::anyhow!("quote not found"))
    }

    /// Создает новую цитату
    pub async fn create(&self, quote: &Quote) -> anyhow::Result<()> {
        sqlx::query(
            "INSERT INTO quotes (id, text, author, likes_count, created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, $6)"
        )
        .bind(&quote.id)
        .bind(&quote.text)
        .bind(&quote.author)
        .bind(quote.likes_count)
        .bind(quote.created_at)
        .bind(quote.updated_at)
        .execute(&self.pool)
        .await?;

        Ok(())
    }

    /// Обновляет существующую цитату
    pub async fn update(&self, id: &str, quote: &Quote) -> anyhow::Result<()> {
        let result = sqlx::query(
            "UPDATE quotes 
             SET text = $1, author = $2, updated_at = $3
             WHERE id = $4"
        )
        .bind(&quote.text)
        .bind(&quote.author)
        .bind(Utc::now())
        .bind(id)
        .execute(&self.pool)
        .await?;

        if result.rows_affected() == 0 {
            return Err(anyhow::anyhow!("quote not found"));
        }

        Ok(())
    }

    /// Удаляет цитату
    pub async fn delete(&self, id: &str) -> anyhow::Result<()> {
        let result = sqlx::query("DELETE FROM quotes WHERE id = $1")
            .bind(id)
            .execute(&self.pool)
            .await?;

        if result.rows_affected() == 0 {
            return Err(anyhow::anyhow!("quote not found"));
        }

        Ok(())
    }

    /// Ставит лайк цитате
    pub async fn like(&self, id: &str, user_ip: &str, user_agent: Option<&str>) -> anyhow::Result<()> {
        // Начинаем транзакцию
        let mut tx = self.pool.begin().await?;

        // Проверяем, не лайкал ли уже этот пользователь эту цитату
        let existing_like: Option<String> = sqlx::query_scalar(
            "SELECT id FROM likes 
             WHERE quote_id = $1 AND user_ip = $2
             FOR UPDATE"
        )
        .bind(id)
        .bind(user_ip)
        .fetch_optional(&mut *tx)
        .await?;

        if existing_like.is_some() {
            return Err(anyhow::anyhow!("you have already liked this quote"));
        }

        // Увеличиваем счетчик лайков
        let result = sqlx::query(
            "UPDATE quotes 
             SET likes_count = likes_count + 1, updated_at = $1
             WHERE id = $2"
        )
        .bind(Utc::now())
        .bind(id)
        .execute(&mut *tx)
        .await?;

        if result.rows_affected() == 0 {
            return Err(anyhow::anyhow!("quote not found"));
        }

        // Сохраняем информацию о лайке
        let like_id = Uuid::new_v4().to_string();
        sqlx::query(
            "INSERT INTO likes (id, quote_id, user_ip, user_agent, created_at)
             VALUES ($1, $2, $3, $4, $5)
             ON CONFLICT (quote_id, user_ip) DO NOTHING"
        )
        .bind(&like_id)
        .bind(id)
        .bind(user_ip)
        .bind(user_agent)
        .bind(Utc::now())
        .execute(&mut *tx)
        .await?;

        tx.commit().await?;
        Ok(())
    }

    /// Проверяет, лайкнул ли пользователь цитату
    pub async fn is_liked(&self, quote_id: &str, user_ip: &str) -> anyhow::Result<bool> {
        let count: i64 = sqlx::query_scalar(
            "SELECT COUNT(*) FROM likes WHERE quote_id = $1 AND user_ip = $2"
        )
        .bind(quote_id)
        .bind(user_ip)
        .fetch_one(&self.pool)
        .await?;

        Ok(count > 0)
    }

    /// Возвращает топ цитату за неделю
    pub async fn get_top_weekly(&self) -> anyhow::Result<Quote> {
        let quote = sqlx::query_as::<_, Quote>(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             WHERE created_at >= NOW() - INTERVAL '7 days'
             ORDER BY likes_count DESC, created_at DESC
             LIMIT 1"
        )
        .fetch_optional(&self.pool)
        .await?;

        quote.ok_or_else(|| anyhow::anyhow!("no quotes found for the last week"))
    }

    /// Возвращает топ цитату за всё время
    pub async fn get_top_all_time(&self) -> anyhow::Result<Quote> {
        let quote = sqlx::query_as::<_, Quote>(
            "SELECT id, text, author, likes_count, created_at, updated_at 
             FROM quotes 
             ORDER BY likes_count DESC, created_at DESC
             LIMIT 1"
        )
        .fetch_optional(&self.pool)
        .await?;

        quote.ok_or_else(|| anyhow::anyhow!("no quotes found"))
    }

    /// Сбрасывает все лайки
    pub async fn reset_likes(&self) -> anyhow::Result<()> {
        let mut tx = self.pool.begin().await?;

        // Обнуляем счетчики лайков
        sqlx::query("UPDATE quotes SET likes_count = 0, updated_at = $1")
            .bind(Utc::now())
            .execute(&mut *tx)
            .await?;

        // Удаляем все записи о лайках
        sqlx::query("DELETE FROM likes")
            .execute(&mut *tx)
            .await?;

        tx.commit().await?;
        Ok(())
    }
}

