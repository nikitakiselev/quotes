mod config;
mod database;
mod handlers;
mod models;
mod repository;
mod router;

use std::env;
use tracing::{error, info};
use tracing_subscriber;

use config::Config;
use database::connect;
use repository::QuoteRepository;
use router::setup_router;

#[tokio::main]
async fn main() {
    // Выводим в stderr для немедленного отображения
    eprintln!("Starting quotes backend (Rust)...");
    std::io::Write::flush(&mut std::io::stderr()).ok();
    
    // Инициализация логирования
    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| tracing_subscriber::EnvFilter::new("info")),
        )
        .init();

    eprintln!("Logger initialized");
    info!("Starting quotes backend (Rust)");

    // Загрузка переменных окружения
    dotenv::dotenv().ok();
    info!("Environment variables loaded");

    // Инициализация конфигурации
    let cfg = Config::load();
    info!(
        "Configuration loaded: DB={}@{}:{}, API_PORT={}",
        cfg.db_user, cfg.db_host, cfg.db_port, cfg.api_port
    );

    // Подключение к базе данных
    info!("Connecting to database...");
    let pool = match connect(&cfg).await {
        Ok(pool) => {
            info!("Database connection established");
            pool
        }
        Err(e) => {
            error!("Failed to connect to database: {}", e);
            eprintln!("ERROR: Failed to connect to database: {}", e);
            std::process::exit(1);
        }
    };

    // Инициализация репозитория
    let repo = QuoteRepository::new(pool);
    info!("Repository initialized");

    // Настройка роутера
    let app = setup_router(repo, &cfg);
    info!("Router configured");

    // Запуск сервера
    let port = env::var("API_PORT").unwrap_or_else(|_| "8081".to_string());
    let addr = format!("0.0.0.0:{}", port);

    info!("Server starting on {}", addr);

    let listener = match tokio::net::TcpListener::bind(&addr).await {
        Ok(listener) => {
            info!("TCP listener bound to {}", addr);
            listener
        }
        Err(e) => {
            error!("Failed to bind to {}: {}", addr, e);
            eprintln!("ERROR: Failed to bind to {}: {}", addr, e);
            std::process::exit(1);
        }
    };

    info!("Server is ready to accept connections on {}", addr);

    if let Err(e) = axum::serve(listener, app).await {
        error!("Server error: {}", e);
        eprintln!("ERROR: Server error: {}", e);
        std::process::exit(1);
    }
}
