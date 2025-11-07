mod config;
mod database;
mod handlers;
mod models;
mod repository;
mod router;

use std::env;
use tracing::info;

use config::Config;
use database::connect;
use repository::QuoteRepository;
use router::setup_router;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    // Инициализация логирования
    tracing_subscriber::fmt()
        .with_env_filter(tracing_subscriber::EnvFilter::from_default_env())
        .init();

    // Загрузка переменных окружения
    dotenv::dotenv().ok();

    // Инициализация конфигурации
    let cfg = Config::load();
    info!("Configuration loaded");

    // Подключение к базе данных
    let pool = connect(&cfg).await?;

    // Инициализация репозитория
    let repo = QuoteRepository::new(pool);

    // Настройка роутера
    let app = setup_router(repo, &cfg);

    // Запуск сервера
    let port = env::var("API_PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("0.0.0.0:{}", port);
    
    info!("Server starting on {}", addr);
    
    let listener = tokio::net::TcpListener::bind(&addr).await?;
    axum::serve(listener, app).await?;

    Ok(())
}

