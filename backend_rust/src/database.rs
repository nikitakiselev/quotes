use sqlx::postgres::{PgConnectOptions, PgSslMode};
use sqlx::PgPool;
use tracing::{error, info};

use crate::config::Config;

/// Подключается к базе данных PostgreSQL
pub async fn connect(cfg: &Config) -> anyhow::Result<PgPool> {
    info!("Connecting to database: {}@{}:{}", cfg.db_user, cfg.db_host, cfg.db_port);
    
    // Создаем опции подключения
    let ssl_mode = match cfg.db_ssl_mode.as_str() {
        "disable" => PgSslMode::Disable,
        "require" => PgSslMode::Require,
        _ => PgSslMode::Prefer,
    };
    
    let options = PgConnectOptions::new()
        .host(&cfg.db_host)
        .port(cfg.db_port.parse().unwrap_or(5432))
        .username(&cfg.db_user)
        .password(&cfg.db_password)
        .database(&cfg.db_name)
        .ssl_mode(ssl_mode);
    
    // Подключаемся к базе данных
    let pool = match PgPool::connect_with(options).await {
        Ok(pool) => pool,
        Err(e) => {
            error!("Failed to connect to database: {}", e);
            return Err(anyhow::anyhow!("Failed to connect to database: {}", e));
        }
    };
    
    // Проверяем соединение
    match sqlx::query("SELECT 1").execute(&pool).await {
        Ok(_) => {
            info!("Database connection established and verified");
        }
        Err(e) => {
            error!("Failed to verify database connection: {}", e);
            return Err(anyhow::anyhow!("Failed to verify database connection: {}", e));
        }
    }
    
    Ok(pool)
}
