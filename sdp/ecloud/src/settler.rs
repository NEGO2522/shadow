use anyhow::{Context, Result};
use ethers::{
    contract::abigen,
    middleware::SignerMiddleware,
    providers::{Http, Provider},
    signers::{coins_bip39::English, LocalWallet, MnemonicBuilder, Signer},
    types::{Address, U256},
};
use std::{str::FromStr, sync::Arc};

abigen!(
    SdpSettlement,
    r#"[
        function settle(bytes32 buyOrderId, bytes32 sellOrderId, uint256 amountIn, uint256 amountOutMin) external
    ]"#
);

pub struct Settler {
    contract: SdpSettlement<SignerMiddleware<Provider<Http>, LocalWallet>>,
}

impl Settler {
    pub fn new() -> Result<Self> {
        let rpc_url = std::env::var("RPC_URL")
            .unwrap_or_else(|_| "https://relay-sepolia.flashbots.net".to_string());
        let mnemonic = std::env::var("MNEMONIC").context("MNEMONIC not set")?;
        let wallet_index: u32 = std::env::var("WALLET_INDEX")
            .unwrap_or_else(|_| "3".to_string())
            .parse()
            .unwrap_or(3);
        let contract_addr = std::env::var("SETTLEMENT_CONTRACT")
            .unwrap_or_else(|_| "0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71".to_string());

        let provider = Provider::<Http>::try_from(rpc_url.as_str())
            .context("Invalid RPC URL")?;

        let wallet = MnemonicBuilder::<English>::default()
            .phrase(mnemonic.as_str())
            .index(wallet_index)?
            .build()
            .context("Failed to build wallet from mnemonic")?
            .with_chain_id(11155111u64); // Sepolia

        tracing::info!("Settlement wallet: {:?}", wallet.address());

        let client = Arc::new(SignerMiddleware::new(provider, wallet));
        let address = Address::from_str(&contract_addr).context("Invalid settlement contract address")?;
        let contract = SdpSettlement::new(address, client);

        Ok(Self { contract })
    }

    /// Call SdpSettlement.settle(buyOrderId, sellOrderId, amountIn, 0).
    /// buyOrderId and sellOrderId are bytes32 derived from the order UUID.
    pub async fn settle(
        &self,
        buy_order_id: [u8; 32],
        sell_order_id: [u8; 32],
        amount_in: U256,
    ) -> Result<String> {
        let tx = self
            .contract
            .settle(buy_order_id, sell_order_id, amount_in, U256::zero())
            .send()
            .await
            .context("settle tx send failed")?;

        let receipt = tx
            .await
            .context("settle tx wait failed")?
            .context("settle tx: no receipt")?;

        let hash = format!("{:?}", receipt.transaction_hash);
        tracing::info!("Settle tx: {}", hash);
        Ok(hash)
    }
}
