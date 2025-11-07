use sqlx::{PgPool, Pool, Postgres};
use tracing::info;

use crate::config::Config;

/// Подключается к базе данных PostgreSQL
pub async fn connect(cfg: &Config) -> anyhow::Result<PgPool> {
    let database_url = cfg.database_url();
    
    info!("Connecting to database: {}@{}:{}", cfg.db_user, cfg.db_host, cfg.db_port);
    
    let pool = PgPool::connect(&database_url).await?;
    
    // Проверяем соединение
    sqlx::query("SELECT 1").execute(&pool).await?;
    
    info!("Database connection established");
    
    Ok(pool)
}

