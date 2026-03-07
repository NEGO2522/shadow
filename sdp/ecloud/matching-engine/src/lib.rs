use std::collections::{BTreeMap, VecDeque};
use std::cmp::Reverse;

use uuid::Uuid;

use sdp_shared::{LimitOrder, MatchResult, Side};

/// Dummy attestation report — in dev mode returns fixed bytes.
/// In `tee` mode this would call into hardware (e.g. Intel TDX / AWS Nitro).
pub struct AttestationReport {
    pub raw: Vec<u8>,
}

pub struct MatchEngine {
    /// Buy side: keyed by Reverse(price) so BTreeMap iteration gives highest price first.
    bids: BTreeMap<Reverse<u64>, VecDeque<LimitOrder>>,
    /// Ask side: keyed by price ascending so lowest ask is first.
    asks: BTreeMap<u64, VecDeque<LimitOrder>>,
}

impl MatchEngine {
    pub fn new() -> Self {
        Self {
            bids: BTreeMap::new(),
            asks: BTreeMap::new(),
        }
    }

    /// Insert an order into the appropriate side of the book.
    pub fn add_order(&mut self, order: LimitOrder) {
        match order.side {
            Side::Buy => self
                .bids
                .entry(Reverse(order.price))
                .or_default()
                .push_back(order),
            Side::Sell => self
                .asks
                .entry(order.price)
                .or_default()
                .push_back(order),
        }
    }

    /// Remove an order by id from whichever side it lives on.
    /// Returns `true` if found and removed.
    pub fn cancel_order(&mut self, order_id: Uuid) -> bool {
        for queue in self.bids.values_mut() {
            if let Some(pos) = queue.iter().position(|o| o.id == order_id) {
                queue.remove(pos);
                return true;
            }
        }
        for queue in self.asks.values_mut() {
            if let Some(pos) = queue.iter().position(|o| o.id == order_id) {
                queue.remove(pos);
                return true;
            }
        }
        false
    }

    /// Price-Time priority crossing loop.
    ///
    /// Iterates while best bid >= best ask, filling FIFO within each price level.
    /// Partial fills are supported: whichever side is exhausted first is removed,
    /// the remainder stays in the book.
    pub fn execute_match(&mut self) -> Vec<MatchResult> {
        let mut results = Vec::new();
        let now_ms = std::time::SystemTime::now()
            .duration_since(std::time::UNIX_EPOCH)
            .unwrap_or_default()
            .as_millis() as u64;

        loop {
            // Peek at best bid and best ask.
            let best_bid_price = match self.bids.keys().next() {
                Some(Reverse(p)) => *p,
                None => break,
            };
            let best_ask_price = match self.asks.keys().next() {
                Some(p) => *p,
                None => break,
            };

            if best_bid_price < best_ask_price {
                break; // No crossing — done.
            }

            // Execution price = resting order's price (the ask, which arrived first if
            // we treat incoming as the bid). Use mid if simultaneous; here we use ask price.
            let exec_price = best_ask_price;

            // Pop front order from each side.
            let bid = self
                .bids
                .get_mut(&Reverse(best_bid_price))
                .and_then(|q| q.pop_front())
                .expect("bid queue non-empty");

            let ask = self
                .asks
                .get_mut(&best_ask_price)
                .and_then(|q| q.pop_front())
                .expect("ask queue non-empty");

            let fill_qty = bid.quantity.min(ask.quantity);

            results.push(MatchResult {
                buy_order_id: bid.id,
                sell_order_id: ask.id,
                price: exec_price,
                quantity: fill_qty,
                timestamp: now_ms,
            });

            // Return partially-filled remainder to front of queue.
            if bid.quantity > fill_qty {
                let remainder = LimitOrder { quantity: bid.quantity - fill_qty, ..bid };
                self.bids
                    .entry(Reverse(best_bid_price))
                    .or_default()
                    .push_front(remainder);
            } else {
                // Prune empty price level.
                if self.bids.get(&Reverse(best_bid_price)).map_or(true, |q| q.is_empty()) {
                    self.bids.remove(&Reverse(best_bid_price));
                }
            }

            if ask.quantity > fill_qty {
                let remainder = LimitOrder { quantity: ask.quantity - fill_qty, ..ask };
                self.asks
                    .entry(best_ask_price)
                    .or_default()
                    .push_front(remainder);
            } else {
                if self.asks.get(&best_ask_price).map_or(true, |q| q.is_empty()) {
                    self.asks.remove(&best_ask_price);
                }
            }
        }

        results
    }

    /// Return an attestation report for the current state of the engine.
    ///
    /// In `dev` mode this returns a dummy payload.
    /// In `tee` mode this would call hardware attestation APIs.
    pub fn attest(&self) -> AttestationReport {
        #[cfg(feature = "dev")]
        {
            AttestationReport {
                raw: b"DEV_ATTESTATION_STUB".to_vec(),
            }
        }
        #[cfg(not(feature = "dev"))]
        {
            // TODO: call Intel TDX / AWS Nitro attestation SDK
            todo!("hardware attestation not yet implemented")
        }
    }

    /// Number of active bids (summed across all price levels).
    pub fn bid_count(&self) -> usize {
        self.bids.values().map(|q| q.len()).sum()
    }

    /// Number of active asks (summed across all price levels).
    pub fn ask_count(&self) -> usize {
        self.asks.values().map(|q| q.len()).sum()
    }
}

impl Default for MatchEngine {
    fn default() -> Self {
        Self::new()
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use sdp_shared::Side;

    fn make_order(side: Side, price: u64, quantity: u64, ts: u64) -> LimitOrder {
        LimitOrder {
            id: Uuid::new_v4(),
            side,
            price,
            quantity,
            timestamp: ts,
            trader_pubkey: "0xtest".to_string(),
        }
    }

    #[test]
    fn no_cross_no_fill() {
        let mut engine = MatchEngine::new();
        engine.add_order(make_order(Side::Buy, 100, 10, 1));
        engine.add_order(make_order(Side::Sell, 110, 10, 2));
        let fills = engine.execute_match();
        assert!(fills.is_empty());
    }

    #[test]
    fn simple_full_fill() {
        let mut engine = MatchEngine::new();
        engine.add_order(make_order(Side::Buy, 100, 10, 1));
        engine.add_order(make_order(Side::Sell, 100, 10, 2));
        let fills = engine.execute_match();
        assert_eq!(fills.len(), 1);
        assert_eq!(fills[0].quantity, 10);
        assert_eq!(fills[0].price, 100);
        assert_eq!(engine.bid_count(), 0);
        assert_eq!(engine.ask_count(), 0);
    }

    #[test]
    fn partial_fill_bid_remainder() {
        let mut engine = MatchEngine::new();
        // Bid wants 15, ask has only 10
        engine.add_order(make_order(Side::Buy, 100, 15, 1));
        engine.add_order(make_order(Side::Sell, 100, 10, 2));
        let fills = engine.execute_match();
        assert_eq!(fills.len(), 1);
        assert_eq!(fills[0].quantity, 10);
        // 5 units remain on bid side
        assert_eq!(engine.bid_count(), 1);
        assert_eq!(engine.ask_count(), 0);
    }

    #[test]
    fn partial_fill_ask_remainder() {
        let mut engine = MatchEngine::new();
        engine.add_order(make_order(Side::Buy, 100, 5, 1));
        engine.add_order(make_order(Side::Sell, 100, 20, 2));
        let fills = engine.execute_match();
        assert_eq!(fills.len(), 1);
        assert_eq!(fills[0].quantity, 5);
        assert_eq!(engine.bid_count(), 0);
        assert_eq!(engine.ask_count(), 1);
    }

    #[test]
    fn fifo_price_time_priority() {
        let mut engine = MatchEngine::new();
        // Two bids at same price — earlier timestamp should fill first
        let bid1 = make_order(Side::Buy, 100, 5, 1);
        let bid2 = make_order(Side::Buy, 100, 5, 2);
        let bid1_id = bid1.id;
        engine.add_order(bid1);
        engine.add_order(bid2);
        engine.add_order(make_order(Side::Sell, 100, 5, 3));

        let fills = engine.execute_match();
        assert_eq!(fills.len(), 1);
        assert_eq!(fills[0].buy_order_id, bid1_id, "first-in bid should fill first");
        assert_eq!(engine.bid_count(), 1, "second bid should remain");
    }

    #[test]
    fn cancel_removes_order() {
        let mut engine = MatchEngine::new();
        let order = make_order(Side::Buy, 100, 10, 1);
        let id = order.id;
        engine.add_order(order);
        assert!(engine.cancel_order(id));
        assert_eq!(engine.bid_count(), 0);
        // Cancelling again returns false
        assert!(!engine.cancel_order(id));
    }

    #[test]
    fn multi_level_crossing() {
        let mut engine = MatchEngine::new();
        // Best bid 105, second 100
        engine.add_order(make_order(Side::Buy, 105, 3, 1));
        engine.add_order(make_order(Side::Buy, 100, 7, 2));
        // Asks at 100 and 103
        engine.add_order(make_order(Side::Sell, 100, 3, 3));
        engine.add_order(make_order(Side::Sell, 103, 3, 4));

        let fills = engine.execute_match();
        // Bid@105 crosses ask@100 (fill 3), then bid@105 exhausted
        // Bid@100 crosses ask@103? No — 100 < 103, so no more crosses
        assert_eq!(fills.len(), 1);
        assert_eq!(fills[0].quantity, 3);
    }

    #[test]
    fn attest_returns_bytes() {
        let engine = MatchEngine::new();
        let report = engine.attest();
        assert!(!report.raw.is_empty());
    }
}
