use std::env;

/// Конфигурация приложения
#[derive(Debug, Clone)]
pub struct Config {
    pub db_host: String,
    pub db_port: String,
    pub db_user: String,
    pub db_password: String,
    pub db_name: String,
    pub db_ssl_mode: String,
    pub api_port: String,
    pub cors_origin: String,
}

impl Config {
    /// Загружает конфигурацию из переменных окружения
    pub fn load() -> Self {
        Self {
            db_host: get_env("DB_HOST", "localhost"),
            db_port: get_env("DB_PORT", "5432"),
            db_user: get_env("DB_USER", "quotes_user"),
            db_password: get_env("DB_PASSWORD", "quotes_password"),
            db_name: get_env("DB_NAME", "quotes_db"),
            db_ssl_mode: get_env("DB_SSLMODE", "disable"),
            api_port: get_env("API_PORT", "8081"),
            cors_origin: get_env("CORS_ORIGIN", "*"),
        }
    }

    /// Возвращает DSN строку для подключения к PostgreSQL
    pub fn database_url(&self) -> String {
        format!(
            "postgresql://{}:{}@{}:{}/{}?sslmode={}",
            self.db_user, self.db_password, self.db_host, self.db_port, self.db_name, self.db_ssl_mode
        )
    }
}

fn get_env(key: &str, default: &str) -> String {
    env::var(key).unwrap_or_else(|_| default.to_string())
}

