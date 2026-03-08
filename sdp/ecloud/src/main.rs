mod handlers;
mod settler;
mod types;

use axum::{
    routing::{get, post},
    Router,
};
use dashmap::DashMap;
use std::sync::Arc;

use settler::Settler;
use types::Order;

pub struct AppState {
    pub orders: DashMap<String, Order>,
    pub settler: Arc<Settler>,
    pub api_key: String,
}

#[tokio::main]
async fn main() {
    // Load .env from sdp-ecloud directory if present
    let _ = dotenvy::from_path("../sdp-ecloud/.env");
    let _ = dotenvy::dotenv();

    tracing_subscriber::fmt()
        .with_env_filter(
            tracing_subscriber::EnvFilter::try_from_default_env()
                .unwrap_or_else(|_| "shadow_matcher=info".parse().unwrap()),
        )
        .init();

    let settler = Settler::new().expect("Failed to initialize settler");
    let api_key = std::env::var("MATCHER_API_KEY").unwrap_or_else(|_| {
        tracing::warn!("MATCHER_API_KEY not set — using empty string (insecure)");
        String::new()
    });
    let port = std::env::var("PORT").unwrap_or_else(|_| "8080".to_string());

    let state = Arc::new(AppState {
        orders: DashMap::new(),
        settler: Arc::new(settler),
        api_key,
    });

    let app = Router::new()
        .route("/health", get(handlers::health))
        .route("/order", post(handlers::post_order))
        .route("/match", post(handlers::post_match))
        .with_state(state);

    let addr = format!("0.0.0.0:{port}");
    tracing::info!("Shadow matcher listening on {addr}");
    let listener = tokio::net::TcpListener::bind(&addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
