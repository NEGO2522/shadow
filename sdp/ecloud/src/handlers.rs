use axum::{
    extract::State,
    http::{HeaderMap, StatusCode},
    response::IntoResponse,
    Json,
};
use std::sync::Arc;

use crate::{
    types::{Fill, MatchResponse, Order, Side},
    AppState,
};

fn check_auth(headers: &HeaderMap, expected: &str) -> bool {
    headers
        .get("authorization")
        .and_then(|v| v.to_str().ok())
        .and_then(|v| v.strip_prefix("Bearer "))
        .map(|t| t == expected)
        .unwrap_or(false)
}

/// POST /order — idempotent by order.id
pub async fn post_order(
    State(state): State<Arc<AppState>>,
    headers: HeaderMap,
    Json(order): Json<Order>,
) -> impl IntoResponse {
    if !check_auth(&headers, &state.api_key) {
        return (StatusCode::UNAUTHORIZED, "Unauthorized").into_response();
    }
    if order.id.is_empty() {
        return (StatusCode::BAD_REQUEST, "Missing order id").into_response();
    }

    // Idempotent insert
    state.orders.entry(order.id.clone()).or_insert(order.clone());
    tracing::info!("Order accepted: {} {:?} {} @ {}", order.id, order.side, order.asset, order.price);
    StatusCode::OK.into_response()
}

/// POST /match — run a matching cycle across all assets
pub async fn post_match(
    State(state): State<Arc<AppState>>,
    headers: HeaderMap,
) -> impl IntoResponse {
    if !check_auth(&headers, &state.api_key) {
        return (StatusCode::UNAUTHORIZED, Json(serde_json::json!({"error":"Unauthorized"}))).into_response();
    }

    // Collect all orders grouped by asset
    let mut by_asset: std::collections::HashMap<String, (Vec<Order>, Vec<Order>)> =
        std::collections::HashMap::new();

    for entry in state.orders.iter() {
        let o = entry.value().clone();
        let (buys, sells) = by_asset.entry(o.asset.clone()).or_default();
        match o.side {
            Side::Buy => buys.push(o),
            Side::Sell => sells.push(o),
        }
    }

    let mut fills: Vec<Fill> = Vec::new();
    let mut to_remove: Vec<String> = Vec::new();

    for (_asset, (mut buys, mut sells)) in by_asset {
        // Sort: buys descending by price, sells ascending
        buys.sort_by(|a, b| b.price.partial_cmp(&a.price).unwrap());
        sells.sort_by(|a, b| a.price.partial_cmp(&b.price).unwrap());

        let mut bi = 0;
        let mut si = 0;
        let mut buy_remaining: Vec<f64> = buys.iter().map(|o| o.quantity).collect();
        let mut sell_remaining: Vec<f64> = sells.iter().map(|o| o.quantity).collect();

        while bi < buys.len() && si < sells.len() {
            if buys[bi].price < sells[si].price {
                break; // No more matches possible
            }

            let fill_price = (buys[bi].price + sells[si].price) / 2.0;
            let fill_qty = buy_remaining[bi].min(sell_remaining[si]);

            if fill_qty <= 0.0 {
                if buy_remaining[bi] <= 0.0 { bi += 1; }
                if si < sells.len() && sell_remaining[si] <= 0.0 { si += 1; }
                continue;
            }

            // Attempt on-chain settlement
            let buy_id_bytes = uuid_to_bytes32(&buys[bi].id);
            let sell_id_bytes = uuid_to_bytes32(&sells[si].id);
            let amount_in = ethers::types::U256::from(
                (fill_qty * fill_price * 1_000_000.0) as u64, // USDC 6 decimals
            );

            let tx_hash = match state.settler.settle(buy_id_bytes, sell_id_bytes, amount_in).await {
                Ok(h) => h,
                Err(e) => {
                    tracing::error!("Settlement failed: {e}");
                    break;
                }
            };

            fills.push(Fill {
                buy_order_id: buys[bi].id.clone(),
                sell_order_id: sells[si].id.clone(),
                asset: buys[bi].asset.clone(),
                price: fill_price,
                quantity: fill_qty,
                tx_hash,
            });

            buy_remaining[bi] -= fill_qty;
            sell_remaining[si] -= fill_qty;

            if buy_remaining[bi] <= 0.0 {
                to_remove.push(buys[bi].id.clone());
                bi += 1;
            }
            if si < sells.len() && sell_remaining[si] <= 0.0 {
                to_remove.push(sells[si].id.clone());
                si += 1;
            }
        }
    }

    // Remove fully filled orders
    for id in &to_remove {
        state.orders.remove(id);
    }

    let resp = MatchResponse {
        fills: fills.len(),
        results: fills,
    };
    (StatusCode::OK, Json(resp)).into_response()
}

/// GET /health
pub async fn health() -> impl IntoResponse {
    StatusCode::OK
}

/// Convert UUID string to bytes32 (first 16 bytes = UUID, last 16 = zeros)
fn uuid_to_bytes32(uuid_str: &str) -> [u8; 32] {
    let mut out = [0u8; 32];
    // Remove hyphens and parse hex
    let hex = uuid_str.replace('-', "");
    if let Ok(bytes) = hex::decode(&hex) {
        let len = bytes.len().min(16);
        out[..len].copy_from_slice(&bytes[..len]);
    }
    out
}
