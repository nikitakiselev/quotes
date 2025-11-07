use axum::{
    http::{Method, StatusCode, Uri},
    response::IntoResponse,
    routing::{delete, get, post, put},
    Router,
};
use tower::ServiceBuilder;
use tower_http::{
    cors::{Any, CorsLayer},
    trace::TraceLayer,
};

use std::sync::Arc;

use crate::config::Config;
use crate::handlers;
use crate::repository::QuoteRepository;

/// Настраивает и возвращает роутер
pub fn setup_router(repo: QuoteRepository, cfg: &Config) -> Router {
    let repo = Arc::new(repo);
    // Настройка CORS
    let cors = if cfg.cors_origin == "*" {
        CorsLayer::new()
            .allow_origin(Any)
            .allow_methods([Method::GET, Method::POST, Method::PUT, Method::DELETE, Method::OPTIONS, Method::PATCH])
            .allow_headers(Any)
            .expose_headers(Any)
    } else {
        CorsLayer::new()
            .allow_origin(cfg.cors_origin.parse().unwrap())
            .allow_methods([Method::GET, Method::POST, Method::PUT, Method::DELETE, Method::OPTIONS, Method::PATCH])
            .allow_headers(Any)
            .expose_headers(Any)
            .allow_credentials(true)
    };

    // API routes
    let api_routes = Router::new()
        .route("/quotes/random", get(handlers::get_random))
        .route("/quotes/top/weekly", get(handlers::get_top_weekly))
        .route("/quotes/top/alltime", get(handlers::get_top_all_time))
        .route("/quotes/likes/reset", delete(handlers::reset_likes))
        .route("/quotes", get(handlers::get_all))
        .route("/quotes", post(handlers::create))
        .route("/quotes/:id", get(handlers::get_by_id))
        .route("/quotes/:id", put(handlers::update))
        .route("/quotes/:id", delete(handlers::delete))
        .route("/quotes/:id/like", put(handlers::like))
        .with_state(repo.clone());

    // Health check
    let health_route = Router::new().route("/health", get(handlers::health));

    // Объединяем все роуты
    Router::new()
        .nest("/api", api_routes)
        .merge(health_route)
        .layer(
            ServiceBuilder::new()
                .layer(TraceLayer::new_for_http())
                .layer(cors),
        )
        .fallback(handle_404)
}

/// Обработчик 404
async fn handle_404(uri: Uri) -> impl IntoResponse {
    // Если запрос к API, возвращаем JSON 404
    if uri.path().starts_with("/api") {
        return (
            StatusCode::NOT_FOUND,
            axum::Json(serde_json::json!({"error": "Not found"})),
        )
            .into_response();
    }

    // Для остальных запросов возвращаем простой 404
    (StatusCode::NOT_FOUND, "Not found").into_response()
}

