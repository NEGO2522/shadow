use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "lowercase")]
pub enum Side {
    Buy,
    Sell,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Order {
    pub id: String,
    pub side: Side,
    pub asset: String,
    /// Price in USDC (e.g. 3500.00)
    pub price: f64,
    /// Quantity of the asset (e.g. 1.5)
    pub quantity: f64,
    /// On-chain trader address (optional, injected by workflow)
    #[serde(default)]
    pub trader: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Fill {
    pub buy_order_id: String,
    pub sell_order_id: String,
    pub asset: String,
    pub price: f64,
    pub quantity: f64,
    pub tx_hash: String,
}

#[derive(Debug, Serialize)]
pub struct MatchResponse {
    pub fills: usize,
    pub results: Vec<Fill>,
}
