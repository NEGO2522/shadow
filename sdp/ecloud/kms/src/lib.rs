use async_trait::async_trait;
use ethers_core::types::{transaction::eip2718::TypedTransaction, Signature, TransactionRequest};
use ethers_signers::{coins_bip39::English, LocalWallet, MnemonicBuilder, Signer as EthersSigner};

/// Trait implemented by any key-management backend that can sign Ethereum transactions.
/// The relayer depends on this trait rather than the concrete `AppWallet`, enabling
/// a drop-in swap to hardware-backed signing in TEE mode.
#[async_trait]
pub trait Signer: Send + Sync {
    async fn sign_tx(&self, tx: TransactionRequest) -> Signature;
    fn address(&self) -> ethers_core::types::Address;
}

/// Wraps a local ECDSA private key wallet.
///
/// - `dev` mode  : key is read from the `APP_WALLET_KEY` environment variable (plain hex).
/// - `tee` mode  : key would be unsealed from hardware KMS (Intel TDX, AWS Nitro, etc.).
#[derive(Clone)]
pub struct AppWallet {
    wallet: LocalWallet,
}

impl AppWallet {
    /// Construct from a hex-encoded private key (with or without `0x` prefix).
    pub fn from_key(hex_key: &str) -> Result<Self, ethers_signers::WalletError> {
        let clean = hex_key
            .strip_prefix("0x")
            .or_else(|| hex_key.strip_prefix("0X"))
            .unwrap_or(hex_key);
        let wallet: LocalWallet = clean.parse()?;
        Ok(Self { wallet })
    }

    /// Read private key from `APP_WALLET_KEY` env var.
    ///
    /// # Panics
    /// Panics if the env var is missing or the key is malformed.
    /// Read the EigenCompute KMS-injected `MNEMONIC` env var and derive wallet index 0.
    pub fn new_from_env() -> Self {
        let mnemonic = std::env::var("MNEMONIC")
            .expect("MNEMONIC injected by EigenCompute KMS");
        // WALLET_INDEX selects which BIP-44 account to use (default 0).
        // Set this to match whichever derived address is funded on-chain.
        let index: u32 = std::env::var("WALLET_INDEX")
            .ok()
            .and_then(|s| s.parse().ok())
            .unwrap_or(0);
        let wallet = MnemonicBuilder::<English>::default()
            .phrase(mnemonic.as_str())
            .index(index)
            .expect("valid index")
            .build()
            .expect("valid mnemonic");
        Self { wallet }
    }

    /// Sign an Ethereum transaction request.
    pub async fn sign_transaction(&self, tx: TransactionRequest) -> Signature {
        let typed: TypedTransaction = tx.into();
        self.wallet
            .sign_transaction(&typed)
            .await
            .expect("signing must not fail for a valid key")
    }

    /// Return the wallet's Ethereum address.
    pub fn address(&self) -> ethers_core::types::Address {
        self.wallet.address()
    }

    /// Return the inner `LocalWallet` for use with ethers middleware.
    pub fn inner(&self) -> &LocalWallet {
        &self.wallet
    }
}

#[async_trait]
impl Signer for AppWallet {
    async fn sign_tx(&self, tx: TransactionRequest) -> Signature {
        self.sign_transaction(tx).await
    }

    fn address(&self) -> ethers_core::types::Address {
        self.address()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // Anvil account 0 private key (well-known test key, never use in production).
    const TEST_KEY: &str = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80";

    #[test]
    fn from_key_parses_and_gives_correct_address() {
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let expected: ethers_core::types::Address =
            "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
                .parse()
                .unwrap();
        assert_eq!(wallet.address(), expected);
    }

    #[test]
    fn from_key_works_without_0x_prefix() {
        let plain = &TEST_KEY[2..]; // strip "0x"
        let wallet = AppWallet::from_key(plain).unwrap();
        let expected: ethers_core::types::Address =
            "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
                .parse()
                .unwrap();
        assert_eq!(wallet.address(), expected);
    }

    #[tokio::test]
    async fn sign_transaction_returns_valid_signature() {
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let tx = TransactionRequest::new();
        let sig = wallet.sign_transaction(tx).await;
        // A valid ECDSA signature has non-zero r and s.
        assert_ne!(sig.r, ethers_core::types::U256::zero());
        assert_ne!(sig.s, ethers_core::types::U256::zero());
    }

    #[tokio::test]
    async fn signer_trait_sign_tx() {
        let wallet = AppWallet::from_key(TEST_KEY).unwrap();
        let sig = wallet.sign_tx(TransactionRequest::new()).await;
        assert_ne!(sig.r, ethers_core::types::U256::zero());
    }
}
