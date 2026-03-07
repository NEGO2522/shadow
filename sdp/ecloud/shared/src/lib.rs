use serde::{Deserialize, Serialize};
use uuid::Uuid;

#[derive(Debug, Clone, Copy, PartialEq, Eq, Serialize, Deserialize)]
pub enum Side {
    Buy,
    Sell,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct LimitOrder {
    pub id: Uuid,
    pub side: Side,
    /// Price in smallest unit (e.g. basis points or wei-equivalent ticks)
    pub price: u64,
    pub quantity: u64,
    /// Unix timestamp in milliseconds
    pub timestamp: u64,
    /// Hex-encoded compressed public key of submitting trader
    pub trader_pubkey: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct MatchResult {
    pub buy_order_id: Uuid,
    pub sell_order_id: Uuid,
    /// Execution price agreed upon (typically the resting order's price)
    pub price: u64,
    pub quantity: u64,
    pub timestamp: u64,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum AbortReason {
    SlippageExceeded,
    ContractRevert(String),
    GasLimitExceeded,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub enum SimResult {
    Ok,
    Abort(AbortReason),
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn limit_order_roundtrip() {
        let order = LimitOrder {
            id: Uuid::new_v4(),
            side: Side::Buy,
            price: 1000,
            quantity: 5,
            timestamp: 0,
            trader_pubkey: "0xdeadbeef".to_string(),
        };
        let json = serde_json::to_string(&order).expect("serialize");
        let decoded: LimitOrder = serde_json::from_str(&json).expect("deserialize");
        assert_eq!(order.id, decoded.id);
        assert_eq!(order.price, decoded.price);
    }

    #[test]
    fn sim_result_variants() {
        let ok = SimResult::Ok;
        let abort = SimResult::Abort(AbortReason::SlippageExceeded);
        assert!(matches!(ok, SimResult::Ok));
        assert!(matches!(abort, SimResult::Abort(AbortReason::SlippageExceeded)));
    }
}
