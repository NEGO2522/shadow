use ethers::{
    middleware::SignerMiddleware,
    providers::{Http, Middleware, Provider},
    signers::Signer as EthersSigner,
    types::{Address, Bytes, TransactionReceipt, TransactionRequest, TxHash, U256},
};

use sdp_kms::AppWallet;
use sdp_shared::MatchResult;

/// Flashbots Sepolia private RPC — transactions bypass the public mempool.
pub const FLASHBOTS_SEPOLIA_RPC: &str = "https://relay-sepolia.flashbots.net";

/// Sepolia chain ID.
const SEPOLIA_CHAIN_ID: u64 = 11_155_111;

/// Gas limit for the `settle()` call (Uniswap v3 single-hop swap + ERC-20 transfers).
const SETTLE_GAS_LIMIT: u64 = 350_000;

/// Relayer that builds, signs, and broadcasts settlement transactions via Flashbots Sepolia.
pub struct SettlementRelayer {
    rpc_url: String,
    wallet:  AppWallet,
}

impl SettlementRelayer {
    /// Create a new relayer.
    ///
    /// `rpc_url` — Flashbots Sepolia RPC (or a local Anvil node in tests).
    pub fn new(rpc_url: &str, wallet: AppWallet) -> Self {
        Self {
            rpc_url: rpc_url.to_string(),
            wallet,
        }
    }

    /// Build, sign, and broadcast a settlement transaction for `result`.
    ///
    /// Reads `SETTLEMENT_CONTRACT` from the environment at call time so that
    /// the address can be injected without restarting the process.
    ///
    /// Returns `Ok(tx_hash)` on success or an error string describing the failure.
    /// The caller should log the error and continue — this must never panic the server.
    pub async fn relay_match(&self, result: MatchResult) -> Result<TxHash, String> {
        let contract_addr = std::env::var("SETTLEMENT_CONTRACT")
            .map_err(|_| "SETTLEMENT_CONTRACT env var not set".to_string())?;
        let to: Address = contract_addr.parse()
            .map_err(|e| format!("SETTLEMENT_CONTRACT parse error: {e}"))?;

        let calldata = encode_settlement_calldata(&result);

        let tx = TransactionRequest::new()
            .to(to)
            .data(calldata)
            .value(U256::zero())
            .gas(SETTLE_GAS_LIMIT);

        let provider = Provider::<Http>::try_from(self.rpc_url.as_str())
            .map_err(|e| format!("RPC URL error: {e}"))?;
        let signer = self.wallet.inner().clone().with_chain_id(SEPOLIA_CHAIN_ID);
        let client = SignerMiddleware::new(provider, signer);

        let pending = client
            .send_transaction(tx, None)
            .await
            .map_err(|e| format!("broadcast failed: {e}"))?;

        Ok(pending.tx_hash())
    }

    /// Fetch the receipt for a previously broadcast transaction.
    /// Returns `None` if the transaction has not been mined yet.
    pub async fn get_receipt(&self, hash: TxHash) -> Option<TransactionReceipt> {
        let provider =
            Provider::<Http>::try_from(self.rpc_url.as_str()).expect("valid RPC URL");
        provider.get_transaction_receipt(hash).await.ok().flatten()
    }
}

// ─── Calldata encoding ───────────────────────────────────────────────────────

/// ABI-encode a `MatchResult` into calldata for `SdpSettlement.settle()`.
///
/// Function signature: `settle(bytes32,bytes32,uint256,uint256)`
/// Selector: keccak256("settle(bytes32,bytes32,uint256,uint256)")[..4] = 0x9e828403
///
/// Layout (132 bytes total):
///   [0..4]    selector
///   [4..36]   buyOrderId   — UUID 16 B, left-aligned in bytes32, zero-padded
///   [36..68]  sellOrderId  — same encoding
///   [68..100] amountIn     — USDC 6-decimal: price_cents × quantity ÷ 100
///   [100..132] amountOutMin — ETH 6-decimal: quantity × 995 ÷ 1000 (0.5 % slippage)
///
/// Price scale assumptions (matching the UI sdpClient.ts):
///   price    — u64 cents  (e.g. 284100 = $2841.00)
///   quantity — u64 in units of 10⁻⁶ ETH (e.g. 1_000_000 = 1.000000 ETH)
pub fn encode_settlement_calldata(result: &MatchResult) -> Bytes {
    // settle(bytes32,bytes32,uint256,uint256)
    const SELECTOR: [u8; 4] = [0x9e, 0x82, 0x84, 0x03];

    // amountIn: USDC with 6 decimals.
    //   price_cents × quantity_eth6 / 100 = USDC_6dec
    //   e.g. 284100 × 1_000_000 / 100 = 2_841_000_000 (= $2841 USDC)
    let amount_in: u128 = (result.price as u128)
        .saturating_mul(result.quantity as u128)
        / 100;

    // amountOutMin: ETH with 6 decimals, 0.5 % slippage applied.
    //   quantity_eth6 × 995 / 1000
    let amount_out_min: u128 = (result.quantity as u128)
        .saturating_mul(995)
        / 1000;

    let mut data: Vec<u8> = Vec::with_capacity(4 + 4 * 32);
    data.extend_from_slice(&SELECTOR);

    // buyOrderId as bytes32 — UUID is 16 bytes, left-aligned.
    let mut slot = [0u8; 32];
    slot[..16].copy_from_slice(result.buy_order_id.as_bytes());
    data.extend_from_slice(&slot);

    // sellOrderId as bytes32.
    let mut slot = [0u8; 32];
    slot[..16].copy_from_slice(result.sell_order_id.as_bytes());
    data.extend_from_slice(&slot);

    // amountIn as uint256 — u128, right-aligned in 32-byte slot.
    let mut slot = [0u8; 32];
    slot[16..32].copy_from_slice(&amount_in.to_be_bytes());
    data.extend_from_slice(&slot);

    // amountOutMin as uint256.
    let mut slot = [0u8; 32];
    slot[16..32].copy_from_slice(&amount_out_min.to_be_bytes());
    data.extend_from_slice(&slot);

    Bytes::from(data)
}

// ─── Tests ───────────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;
    use sdp_kms::AppWallet;
    use uuid::Uuid;

    const TEST_KEY: &str = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";

    fn test_match() -> MatchResult {
        MatchResult {
            buy_order_id:  Uuid::new_v4(),
            sell_order_id: Uuid::new_v4(),
            price:    284_100, // $2841.00
            quantity: 1_000_000, // 1.000000 ETH
            timestamp: 0,
        }
    }

    #[test]
    fn selector_is_correct() {
        let result = test_match();
        let calldata = encode_settlement_calldata(&result);
        // keccak256("settle(bytes32,bytes32,uint256,uint256)")[..4] = 9e828403
        assert_eq!(&calldata[..4], &[0x9e, 0x82, 0x84, 0x03]);
    }

    #[test]
    fn calldata_length() {
        let calldata = encode_settlement_calldata(&test_match());
        assert_eq!(calldata.len(), 132); // 4 + 4×32
    }

    #[test]
    fn amount_in_correct() {
        // price=284100 cents, quantity=1_000_000 eth6
        // amountIn = 284100 × 1_000_000 / 100 = 2_841_000_000 (= $2841 USDC 6-dec)
        let result = test_match();
        let calldata = encode_settlement_calldata(&result);
        // amountIn slot: bytes [68..100], value in last 16 bytes [84..100]
        let mut buf = [0u8; 16];
        buf.copy_from_slice(&calldata[84..100]);
        let amount_in = u128::from_be_bytes(buf);
        assert_eq!(amount_in, 2_841_000_000u128);
    }

    #[test]
    fn amount_out_min_is_995_permille_of_quantity() {
        // quantity=1_000_000 → amountOutMin = 1_000_000 × 995 / 1000 = 995_000
        let result = test_match();
        let calldata = encode_settlement_calldata(&result);
        // amountOutMin slot: bytes [100..132], value in last 16 bytes [116..132]
        let mut buf = [0u8; 16];
        buf.copy_from_slice(&calldata[116..132]);
        let amount_out_min = u128::from_be_bytes(buf);
        assert_eq!(amount_out_min, 995_000u128);
    }

    #[test]
    fn buy_order_id_encoded() {
        let result = test_match();
        let calldata = encode_settlement_calldata(&result);
        // buyOrderId slot: [4..36], UUID bytes in first 16 bytes.
        let uuid_bytes = result.buy_order_id.as_bytes();
        assert_eq!(&calldata[4..20], uuid_bytes.as_slice());
        // Remaining 16 bytes must be zero.
        assert_eq!(&calldata[20..36], &[0u8; 16]);
    }

    #[test]
    fn encode_calldata_price_slot() {
        let result = MatchResult {
            buy_order_id:  Uuid::new_v4(),
            sell_order_id: Uuid::new_v4(),
            price:    100,  // $1.00
            quantity: 1_000_000, // 1 ETH6
            timestamp: 0,
        };
        let calldata = encode_settlement_calldata(&result);
        // amountIn = 100 × 1_000_000 / 100 = 1_000_000 (= $1 USDC)
        let mut buf = [0u8; 16];
        buf.copy_from_slice(&calldata[84..100]);
        assert_eq!(u128::from_be_bytes(buf), 1_000_000u128);
    }

    #[tokio::test]
    async fn wallet_signs_transaction() {
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let tx = TransactionRequest::new().to(
            Address::zero(),
        );
        let sig = wallet.sign_transaction(tx).await;
        assert_ne!(sig.r, U256::zero());
    }

    #[test]
    fn relayer_constructs() {
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let _relayer = SettlementRelayer::new(FLASHBOTS_SEPOLIA_RPC, wallet);
    }

    /// Full relay smoke test — only runs when ANVIL_URL env var is set.
    #[tokio::test]
    async fn relay_match_smoke() {
        let anvil_url = match std::env::var("ANVIL_URL") {
            Ok(u) => u,
            Err(_) => return,
        };
        // Also requires SETTLEMENT_CONTRACT to be set.
        if std::env::var("SETTLEMENT_CONTRACT").is_err() {
            return;
        }
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let relayer = SettlementRelayer::new(&anvil_url, wallet);
        let result = relayer.relay_match(test_match()).await;
        assert!(result.is_ok(), "relay_match failed: {:?}", result.err());
    }
}
