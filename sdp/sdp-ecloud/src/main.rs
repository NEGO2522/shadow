//! SDP ECloud — EigenCompute entry point.
//!
//! Exposes a minimal HTTP API for submitting orders and triggering settlement.
//! Runs inside an Intel TDX TEE on EigenCompute; the KMS injects `MNEMONIC`
//! at startup so no private key ever appears in environment config or logs.
//!
//! Endpoints:
//!   POST /order   — submit a LimitOrder to the dark pool order book
//!   POST /match   — run matching cycle: match → pre-screen → relay
//!   GET  /health  — liveness probe

use std::sync::Arc;
use tokio::sync::Mutex;

use axum::{
    body::Body,
    extract::State,
    http::{Method, Request, StatusCode},
    middleware::{self, Next},
    response::{Json, Response},
    routing::{get, post},
    Router,
};
use tower_http::cors::{Any, CorsLayer};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use sdp_kms::AppWallet;
use sdp_matching_engine::MatchEngine;
use sdp_pre_screener::{Screener, TxData};
use sdp_relayer::{encode_settlement_calldata, SettlementRelayer};
use sdp_shared::{LimitOrder, MatchResult, Side, SimResult};

// ─── Shared state ────────────────────────────────────────────────────────────

struct AppState {
    engine:   Mutex<MatchEngine>,
    screener: Screener,
    relayer:  SettlementRelayer,
}

// ─── Request / response types ────────────────────────────────────────────────

#[derive(Debug, Deserialize)]
struct OrderRequest {
    id:            Option<String>, // client-provided UUID for idempotency
    side:          String, // "buy" | "sell"
    price:         u64,
    quantity:      u64,
    trader_pubkey: String,
}

#[derive(Debug, Serialize)]
struct OrderResponse {
    id:      String,
    status:  String,
}

#[derive(Debug, Serialize)]
struct MatchResponse {
    fills:    usize,
    results:  Vec<FillInfo>,
}

#[derive(Debug, Serialize)]
struct FillInfo {
    buy_order_id:  String,
    sell_order_id: String,
    price:         u64,
    quantity:      u64,
    tx_hash:       Option<String>,
    screener:      String,
}

// ─── Auth middleware ──────────────────────────────────────────────────────────

async fn require_api_key(req: Request<Body>, next: Next) -> Result<Response, StatusCode> {
    let expected = std::env::var("MATCHER_API_KEY").unwrap_or_default();
    let token = req
        .headers()
        .get("Authorization")
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "))
        .unwrap_or("");
    if !expected.is_empty() && token != expected {
        return Err(StatusCode::UNAUTHORIZED);
    }
    Ok(next.run(req).await)
}

// ─── Handlers ────────────────────────────────────────────────────────────────

async fn health() -> &'static str {
    "ok"
}

async fn submit_order(
    State(state): State<Arc<AppState>>,
    Json(req): Json<OrderRequest>,
) -> Result<Json<OrderResponse>, StatusCode> {
    let side = match req.side.to_lowercase().as_str() {
        "buy"  => Side::Buy,
        "sell" => Side::Sell,
        _      => return Err(StatusCode::BAD_REQUEST),
    };

    let order_id = req.id
        .as_deref()
        .and_then(|s| Uuid::parse_str(s).ok())
        .unwrap_or_else(Uuid::new_v4);

    let order = LimitOrder {
        id:            order_id,
        side,
        price:         req.price,
        quantity:      req.quantity,
        timestamp:     now_ms(),
        trader_pubkey: req.trader_pubkey,
    };

    let id = order.id.to_string();
    state.engine.lock().await.add_order(order);

    Ok(Json(OrderResponse { id, status: "queued".into() }))
}

async fn run_match(
    State(state): State<Arc<AppState>>,
) -> Json<MatchResponse> {
    let fills: Vec<MatchResult> = state.engine.lock().await.execute_match();

    let mut results = Vec::with_capacity(fills.len());

    for fill in &fills {
        // Pre-screen the settlement tx before broadcasting.
        let tx = build_tx_data(fill);
        let sim = state.screener.simulate_settlement(tx.clone());

        let settlement_configured =
            std::env::var("SETTLEMENT_CONTRACT").map_or(false, |s| !s.is_empty());

        let (tx_hash, screener_status) = match sim {
            SimResult::Ok if settlement_configured => {
                match state.relayer.relay_match(fill.clone()).await {
                    Ok(hash) => (Some(format!("{hash:?}")), "ok".to_string()),
                    Err(e)   => {
                        eprintln!("relay error for fill {}: {e}", fill.buy_order_id);
                        (None, format!("relay error: {e}"))
                    }
                }
            }
            SimResult::Ok => {
                (None, "ok (relay skipped: SETTLEMENT_CONTRACT not set)".to_string())
            }
            SimResult::Abort(reason) => {
                (None, format!("aborted: {:?}", reason))
            }
        };

        results.push(FillInfo {
            buy_order_id:  fill.buy_order_id.to_string(),
            sell_order_id: fill.sell_order_id.to_string(),
            price:         fill.price,
            quantity:      fill.quantity,
            tx_hash,
            screener:      screener_status,
        });
    }

    let fill_count = results.len();
    if fill_count > 0 {
        println!("SDP ECloud | matched {fill_count} fill(s) this cycle");
    }
    Json(MatchResponse { fills: fill_count, results })
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

fn now_ms() -> u64 {
    std::time::SystemTime::now()
        .duration_since(std::time::UNIX_EPOCH)
        .unwrap_or_default()
        .as_millis() as u64
}

fn build_tx_data(result: &MatchResult) -> TxData {
    use ethers_core::types::{Address, U256};

    let to: Address = std::env::var("SETTLEMENT_CONTRACT")
        .ok()
        .and_then(|s| s.parse().ok())
        .unwrap_or(Address::zero());

    // Encode the real settle() calldata so the pre-screener simulates the
    // exact transaction that the relayer will broadcast.
    let calldata = encode_settlement_calldata(result);

    TxData {
        to,
        calldata,
        value: U256::zero(),
        gas_limit: 350_000,
    }
}

// ─── Main ────────────────────────────────────────────────────────────────────

#[tokio::main]
async fn main() {
    // Load .env in dev; in TEE, MNEMONIC is injected by EigenCompute KMS.
    let _ = dotenvy::dotenv();

    // Initialise wallet from KMS-provided mnemonic.
    let wallet = AppWallet::new_from_env();
    println!("SDP ECloud | wallet address: {}", wallet.address());

    // Relayer uses Flashbots Sepolia by default; override with RPC_URL env var.
    let rpc_url = std::env::var("RPC_URL")
        .unwrap_or_else(|_| sdp_relayer::FLASHBOTS_SEPOLIA_RPC.to_string());

    let state = Arc::new(AppState {
        engine:   Mutex::new(MatchEngine::new()),
        screener: Screener::new(),
        relayer:  SettlementRelayer::new(&rpc_url, wallet),
    });

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods([Method::GET, Method::POST])
        .allow_headers(Any);

    let protected = Router::new()
        .route("/order", post(submit_order))
        .route("/match", post(run_match))
        .layer(middleware::from_fn(require_api_key))
        .with_state(state.clone());

    let app = Router::new()
        .route("/health", get(health))
        .merge(protected)
        .layer(cors);

    // EigenCompute requires binding to 0.0.0.0.
    let addr = "0.0.0.0:3000";
    println!("SDP ECloud | listening on {addr}");

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
