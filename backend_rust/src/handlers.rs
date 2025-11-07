use axum::{
    extract::{Path, Query, State},
    http::{HeaderMap, StatusCode},
    response::Json,
};
use serde::Deserialize;
use std::collections::HashMap;
use std::sync::Arc;

use crate::models::{
    calculate_total_pages, CreateQuoteRequest, PaginatedQuotesResponse, QuoteResponse,
    UpdateQuoteRequest,
};
use crate::repository::QuoteRepository;

/// Получает IP адрес пользователя из запроса
fn get_user_ip(headers: &HeaderMap) -> String {
    // Проверяем X-Forwarded-For
    if let Some(forwarded) = headers.get("x-forwarded-for") {
        if let Ok(forwarded_str) = forwarded.to_str() {
            if let Some(first_ip) = forwarded_str.split(',').next() {
                return first_ip.trim().to_string();
            }
        }
    }

    // Проверяем X-Real-IP
    if let Some(real_ip) = headers.get("x-real-ip") {
        if let Ok(ip_str) = real_ip.to_str() {
            return ip_str.to_string();
        }
    }

    // По умолчанию возвращаем localhost
    "127.0.0.1".to_string()
}

/// Получает случайную цитату
pub async fn get_random(
    State(repo): State<Arc<QuoteRepository>>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let quote = repo.get_random().await.map_err(|_| StatusCode::NOT_FOUND)?;
    let user_ip = get_user_ip(&headers);
    let is_liked = repo.is_liked(&quote.id, &user_ip).await.unwrap_or(false);
    Ok(Json(quote.to_response(is_liked)))
}

/// Параметры запроса для получения всех цитат
#[derive(Debug, Deserialize)]
pub struct GetAllParams {
    page: Option<i32>,
    page_size: Option<i32>,
    search: Option<String>,
}

/// Получает все цитаты с пагинацией
pub async fn get_all(
    State(repo): State<Arc<QuoteRepository>>,
    Query(params): Query<GetAllParams>,
    headers: HeaderMap,
) -> Result<Json<PaginatedQuotesResponse>, StatusCode> {
    let page = params.page.unwrap_or(1).max(1);
    let page_size = params.page_size.unwrap_or(10).max(1).min(100);
    let search = params.search.as_deref();

    let (quotes, total) = repo
        .get_all(page, page_size, search)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let user_ip = get_user_ip(&headers);
    let mut responses = Vec::new();
    for quote in &quotes {
        let is_liked = repo.is_liked(&quote.id, &user_ip).await.unwrap_or(false);
        responses.push(quote.to_response(is_liked));
    }

    let total_pages = calculate_total_pages(total, page_size);

    Ok(Json(PaginatedQuotesResponse {
        quotes: responses,
        total,
        page,
        page_size,
        total_pages,
    }))
}

/// Получает цитату по ID
pub async fn get_by_id(
    State(repo): State<Arc<QuoteRepository>>,
    Path(id): Path<String>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let quote = repo.get_by_id(&id).await.map_err(|_| StatusCode::NOT_FOUND)?;
    let user_ip = get_user_ip(&headers);
    let is_liked = repo.is_liked(&quote.id, &user_ip).await.unwrap_or(false);
    Ok(Json(quote.to_response(is_liked)))
}

/// Создает новую цитату
pub async fn create(
    State(repo): State<Arc<QuoteRepository>>,
    Json(req): Json<CreateQuoteRequest>,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let quote = crate::models::Quote::from_request(req);
    repo.create(&quote)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    // Новая цитата не может быть лайкнута
    Ok(Json(quote.to_response(false)))
}

/// Обновляет существующую цитату
pub async fn update(
    State(repo): State<Arc<QuoteRepository>>,
    Path(id): Path<String>,
    Json(req): Json<UpdateQuoteRequest>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let mut quote = repo.get_by_id(&id).await.map_err(|_| StatusCode::NOT_FOUND)?;

    if let Some(text) = req.text {
        quote.text = text;
    }
    if let Some(author) = req.author {
        quote.author = author;
    }

    repo.update(&id, &quote)
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;

    let updated_quote = repo.get_by_id(&id).await.map_err(|_| StatusCode::NOT_FOUND)?;
    let user_ip = get_user_ip(&headers);
    let is_liked = repo
        .is_liked(&updated_quote.id, &user_ip)
        .await
        .unwrap_or(false);
    Ok(Json(updated_quote.to_response(is_liked)))
}

/// Удаляет цитату
pub async fn delete(
    State(repo): State<Arc<QuoteRepository>>,
    Path(id): Path<String>,
) -> Result<StatusCode, StatusCode> {
    repo.delete(&id)
        .await
        .map_err(|_| StatusCode::NOT_FOUND)?;
    Ok(StatusCode::NO_CONTENT)
}

/// Ставит лайк цитате
pub async fn like(
    State(repo): State<Arc<QuoteRepository>>,
    Path(id): Path<String>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, (StatusCode, Json<serde_json::Value>)> {
    let user_ip = get_user_ip(&headers);
    let user_agent = headers
        .get("user-agent")
        .and_then(|h| h.to_str().ok())
        .map(|s| s.to_string());

    if let Err(e) = repo.like(&id, &user_ip, user_agent.as_deref()).await {
        let error_msg = e.to_string();
        if error_msg.contains("already liked") || error_msg.contains("have already liked") {
            return Err((
                StatusCode::BAD_REQUEST,
                Json(serde_json::json!({"error": "Вы уже поставили лайк этой цитате"})),
            ));
        }
        if error_msg.contains("not found") {
            return Err((
                StatusCode::NOT_FOUND,
                Json(serde_json::json!({"error": "quote not found"})),
            ));
        }
        return Err((
            StatusCode::INTERNAL_SERVER_ERROR,
            Json(serde_json::json!({"error": "internal server error"})),
        ));
    }

    let quote = repo.get_by_id(&id).await.map_err(|_| {
        (
            StatusCode::NOT_FOUND,
            Json(serde_json::json!({"error": "quote not found"})),
        )
    })?;
    // После лайка пользователь точно лайкнул эту цитату
    Ok(Json(quote.to_response(true)))
}

/// Получает топ цитату за неделю
pub async fn get_top_weekly(
    State(repo): State<Arc<QuoteRepository>>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let quote = repo
        .get_top_weekly()
        .await
        .map_err(|_| StatusCode::NOT_FOUND)?;
    let user_ip = get_user_ip(&headers);
    let is_liked = repo.is_liked(&quote.id, &user_ip).await.unwrap_or(false);
    Ok(Json(quote.to_response(is_liked)))
}

/// Получает топ цитату за всё время
pub async fn get_top_all_time(
    State(repo): State<Arc<QuoteRepository>>,
    headers: HeaderMap,
) -> Result<Json<QuoteResponse>, StatusCode> {
    let quote = repo
        .get_top_all_time()
        .await
        .map_err(|_| StatusCode::NOT_FOUND)?;
    let user_ip = get_user_ip(&headers);
    let is_liked = repo.is_liked(&quote.id, &user_ip).await.unwrap_or(false);
    Ok(Json(quote.to_response(is_liked)))
}

/// Сбрасывает все лайки
pub async fn reset_likes(
    State(repo): State<Arc<QuoteRepository>>,
) -> Result<Json<HashMap<String, String>>, StatusCode> {
    repo.reset_likes()
        .await
        .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;
    let mut response = HashMap::new();
    response.insert("message".to_string(), "Все лайки успешно сброшены".to_string());
    Ok(Json(response))
}

/// Health check
pub async fn health() -> Json<HashMap<String, String>> {
    let mut response = HashMap::new();
    response.insert("status".to_string(), "ok".to_string());
    Json(response)
}

