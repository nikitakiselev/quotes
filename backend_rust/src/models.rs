use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize, Serializer};
use uuid::Uuid;

fn serialize_datetime<S>(dt: &DateTime<Utc>, serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    serializer.serialize_str(&dt.to_rfc3339())
}

/// Цитата в системе
#[derive(Debug, Clone, Serialize, Deserialize, sqlx::FromRow)]
pub struct Quote {
    pub id: String,
    pub text: String,
    pub author: String,
    #[sqlx(rename = "likes_count")]
    pub likes_count: i32,
    #[sqlx(rename = "created_at")]
    pub created_at: DateTime<Utc>,
    #[sqlx(rename = "updated_at")]
    pub updated_at: DateTime<Utc>,
}

/// Запрос на создание цитаты
#[derive(Debug, Deserialize)]
pub struct CreateQuoteRequest {
    pub text: String,
    pub author: String,
}

/// Запрос на обновление цитаты
#[derive(Debug, Deserialize)]
pub struct UpdateQuoteRequest {
    pub text: Option<String>,
    pub author: Option<String>,
}

/// Ответ API с цитатой
#[derive(Debug, Serialize)]
pub struct QuoteResponse {
    pub id: String,
    pub text: String,
    pub author: String,
    pub likes_count: i32,
    pub is_liked: bool,
    #[serde(serialize_with = "serialize_datetime")]
    pub created_at: DateTime<Utc>,
    #[serde(serialize_with = "serialize_datetime")]
    pub updated_at: DateTime<Utc>,
}

/// Ответ API с пагинацией
#[derive(Debug, Serialize)]
pub struct PaginatedQuotesResponse {
    pub quotes: Vec<QuoteResponse>,
    pub total: i64,
    pub page: i32,
    pub page_size: i32,
    pub total_pages: i32,
}

impl Quote {
    /// Преобразует Quote в QuoteResponse
    pub fn to_response(&self, is_liked: bool) -> QuoteResponse {
        QuoteResponse {
            id: self.id.clone(),
            text: self.text.clone(),
            author: self.author.clone(),
            likes_count: self.likes_count,
            is_liked,
            created_at: self.created_at,
            updated_at: self.updated_at,
        }
    }

    /// Создает новую цитату из запроса
    pub fn from_request(req: CreateQuoteRequest) -> Self {
        let now = Utc::now();
        Self {
            id: Uuid::new_v4().to_string(),
            text: req.text,
            author: req.author,
            likes_count: 0,
            created_at: now,
            updated_at: now,
        }
    }
}

/// Вычисляет общее количество страниц
pub fn calculate_total_pages(total: i64, page_size: i32) -> i32 {
    ((total as f64) / (page_size as f64)).ceil() as i32
}

