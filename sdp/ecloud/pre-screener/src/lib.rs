//! Pre-screener: sandboxed EVM simulation before settlement broadcast.
//!
//! # Feature flags
//!
//! | Feature | Screener behaviour |
//! |---------|-------------------|
//! | `dev` (default) | Mini pure-Rust EVM against an in-memory contract store |
//! | `tee` | `revm`-based simulation against a forked Sepolia state via Alchemy |
//!
//! ## `tee` mode dependencies (not yet in Cargo.toml — add before enabling):
//! ```toml
//! [features]
//! tee = ["revm", "alloy-provider", "alloy-network", "alloy-rpc-types", "tokio"]
//!
//! [dependencies]
//! revm            = { version = "14", default-features = false, features = ["std"] }
//! alloy-provider  = { version = "0.9", features = ["reqwest"] }
//! alloy-network   = "0.9"
//! alloy-rpc-types = "0.9"
//! ```
//!
//! ## `tee` simulation flow:
//! 1. Build an `alloy` HTTP provider against the Alchemy Sepolia RPC.
//! 2. Fetch the latest block number and create a `revm::db::AlloyDB` fork at that block.
//! 3. Run the `settle()` calldata through revm's `EVM::transact()`.
//! 4. Inspect the execution result for reverts, gas, and Uniswap price impact.

use ethers_core::types::{Address, Bytes, U256};
use sdp_shared::{AbortReason, SimResult};

/// Maximum gas a settlement tx may consume (500 000).
const MAX_GAS: u64 = 500_000;
/// Maximum allowed price-impact in basis points (50 bps = 0.5 %).
const MAX_SLIPPAGE_BPS: u64 = 50;

// ─── TxData ──────────────────────────────────────────────────────────────────

/// Transaction data submitted for pre-screening before broadcast.
#[derive(Debug, Clone)]
pub struct TxData {
    /// Target contract address.
    pub to: Address,
    /// ABI-encoded calldata.
    pub calldata: Bytes,
    /// ETH value (wei).
    pub value: U256,
    /// Gas budget for the simulated call.
    pub gas_limit: u64,
}

// ─── Mini-EVM ────────────────────────────────────────────────────────────────

/// Result of executing bytecode in the mini-EVM.
#[derive(Debug)]
enum EvmOutcome {
    /// Successful return with return data.
    Return(Vec<u8>),
    /// Explicit revert with optional reason bytes.
    Revert(Vec<u8>),
    /// Gas exhausted; `gas_used` = gas_limit.
    OutOfGas { gas_used: u64 },
}

/// Minimal stack-machine EVM interpreter supporting the opcodes needed for
/// dev-mode contract stubs.
///
/// Supported opcodes:
/// `STOP`(0x00), `MSTORE`(0x52), `JUMP`(0x56), `JUMPDEST`(0x5b),
/// `PUSH1`(0x60), `RETURN`(0xf3), `REVERT`(0xfd).
struct MiniEvm<'a> {
    code: &'a [u8],
    stack: Vec<[u8; 32]>,
    memory: Vec<u8>,
    pc: usize,
    gas_used: u64,
    gas_limit: u64,
}

impl<'a> MiniEvm<'a> {
    fn new(code: &'a [u8], gas_limit: u64) -> Self {
        Self {
            code,
            stack: Vec::new(),
            memory: Vec::new(),
            pc: 0,
            gas_used: 0,
            gas_limit,
        }
    }

    fn use_gas(&mut self, amount: u64) -> bool {
        self.gas_used += amount;
        self.gas_used <= self.gas_limit
    }

    fn push_u256(&mut self, value: [u8; 32]) {
        self.stack.push(value);
    }

    fn pop_u256(&mut self) -> Option<[u8; 32]> {
        self.stack.pop()
    }

    /// Read a big-endian u64 from the least-significant 8 bytes of a 32-byte slot.
    fn slot_to_u64(slot: &[u8; 32]) -> u64 {
        let mut buf = [0u8; 8];
        buf.copy_from_slice(&slot[24..32]);
        u64::from_be_bytes(buf)
    }

    fn mem_ensure(&mut self, offset: usize, size: usize) {
        let needed = offset + size;
        if self.memory.len() < needed {
            self.memory.resize(needed, 0);
        }
    }

    fn run(mut self) -> EvmOutcome {
        loop {
            if self.pc >= self.code.len() {
                return EvmOutcome::Return(vec![]);
            }

            let op = self.code[self.pc];
            self.pc += 1;

            match op {
                // STOP
                0x00 => return EvmOutcome::Return(vec![]),

                // MSTORE: pop offset, pop value; mem[offset..offset+32] = value (big-endian)
                0x52 => {
                    if !self.use_gas(3) {
                        return EvmOutcome::OutOfGas { gas_used: self.gas_limit };
                    }
                    let offset_slot = match self.pop_u256() {
                        Some(v) => v,
                        None => return EvmOutcome::Revert(b"stack underflow".to_vec()),
                    };
                    let value_slot = match self.pop_u256() {
                        Some(v) => v,
                        None => return EvmOutcome::Revert(b"stack underflow".to_vec()),
                    };
                    let offset = Self::slot_to_u64(&offset_slot) as usize;
                    self.mem_ensure(offset, 32);
                    self.memory[offset..offset + 32].copy_from_slice(&value_slot);
                }

                // JUMP: pop dest, jump to JUMPDEST
                0x56 => {
                    if !self.use_gas(8) {
                        return EvmOutcome::OutOfGas { gas_used: self.gas_limit };
                    }
                    let dest_slot = match self.pop_u256() {
                        Some(v) => v,
                        None => return EvmOutcome::Revert(b"stack underflow".to_vec()),
                    };
                    let dest = Self::slot_to_u64(&dest_slot) as usize;
                    if dest >= self.code.len() || self.code[dest] != 0x5b {
                        return EvmOutcome::Revert(b"bad jump dest".to_vec());
                    }
                    self.pc = dest;
                }

                // JUMPDEST: no-op marker
                0x5b => {
                    if !self.use_gas(1) {
                        return EvmOutcome::OutOfGas { gas_used: self.gas_limit };
                    }
                }

                // PUSH1: push next byte as 32-byte big-endian value
                0x60 => {
                    if !self.use_gas(3) {
                        return EvmOutcome::OutOfGas { gas_used: self.gas_limit };
                    }
                    if self.pc >= self.code.len() {
                        return EvmOutcome::Revert(b"truncated PUSH1".to_vec());
                    }
                    let byte = self.code[self.pc];
                    self.pc += 1;
                    let mut slot = [0u8; 32];
                    slot[31] = byte;
                    self.push_u256(slot);
                }

                // RETURN: pop offset, pop size; return mem[offset..offset+size]
                0xf3 => {
                    let offset_slot = match self.pop_u256() {
                        Some(v) => v,
                        None => return EvmOutcome::Return(vec![]),
                    };
                    let size_slot = match self.pop_u256() {
                        Some(v) => v,
                        None => return EvmOutcome::Return(vec![]),
                    };
                    let offset = Self::slot_to_u64(&offset_slot) as usize;
                    let size = Self::slot_to_u64(&size_slot) as usize;
                    self.mem_ensure(offset, size);
                    return EvmOutcome::Return(self.memory[offset..offset + size].to_vec());
                }

                // REVERT: pop offset, pop size; revert with mem[offset..offset+size]
                0xfd => {
                    let offset_slot = self.pop_u256().unwrap_or([0u8; 32]);
                    let size_slot = self.pop_u256().unwrap_or([0u8; 32]);
                    let offset = Self::slot_to_u64(&offset_slot) as usize;
                    let size = Self::slot_to_u64(&size_slot) as usize;
                    if size > 0 {
                        self.mem_ensure(offset, size);
                        return EvmOutcome::Revert(
                            self.memory[offset..offset + size].to_vec(),
                        );
                    }
                    return EvmOutcome::Revert(vec![]);
                }

                unknown => {
                    return EvmOutcome::Revert(
                        format!("unsupported opcode 0x{:02x}", unknown).into_bytes(),
                    );
                }
            }

            // Per-step gas guard for long-running loops.
            if self.gas_used > MAX_GAS + self.gas_limit {
                return EvmOutcome::OutOfGas { gas_used: self.gas_limit };
            }
        }
    }
}

// ─── Screener ────────────────────────────────────────────────────────────────

/// In-memory contract store used by the dev-mode screener.
type ContractStore = std::collections::HashMap<[u8; 20], Vec<u8>>;

/// Pre-screener that simulates settlement transactions before broadcast.
pub struct Screener {
    contracts: ContractStore,
}

impl Screener {
    /// Create a screener with an empty contract store.
    pub fn new() -> Self {
        Self {
            contracts: ContractStore::new(),
        }
    }

    /// Register a contract at `addr` with `bytecode` for dev-mode simulation.
    ///
    /// In TEE mode this method is unused; state comes from the forked Sepolia DB.
    pub fn with_contract(mut self, addr: Address, bytecode: Vec<u8>) -> Self {
        self.contracts.insert(addr.0, bytecode);
        self
    }

    /// Simulate `tx_data` and decide whether the settlement is safe.
    ///
    /// Checks (in order):
    /// 1. EVM execution — gas used, return data, revert reason.
    /// 2. Gas guard: `gas_used > 500_000` → `Abort(GasLimitExceeded)`.
    /// 3. Revert guard → `Abort(ContractRevert(reason))`.
    /// 4. Slippage guard: first 32 bytes of return data encode price-impact in bps;
    ///    if > 50 bps → `Abort(SlippageExceeded)`.
    /// 5. Otherwise → `SimResult::Ok`.
    pub fn simulate_settlement(&self, tx_data: TxData) -> SimResult {
        let bytecode = match self.contracts.get(&tx_data.to.0) {
            Some(code) => code.as_slice(),
            None => {
                // No code at address — treat as EOA transfer (always succeeds with 21 000 gas).
                return SimResult::Ok;
            }
        };

        let outcome = MiniEvm::new(bytecode, tx_data.gas_limit).run();

        match outcome {
            EvmOutcome::OutOfGas { gas_used } => {
                if gas_used > MAX_GAS {
                    SimResult::Abort(AbortReason::GasLimitExceeded)
                } else {
                    SimResult::Abort(AbortReason::ContractRevert("out of gas".to_string()))
                }
            }

            EvmOutcome::Revert(data) => {
                let reason = decode_revert_string(&data);
                SimResult::Abort(AbortReason::ContractRevert(reason))
            }

            EvmOutcome::Return(data) => {
                // Gas guard (checked on success path too).
                // We approximate gas_used from the outcome; in a real EVM this comes from
                // the execution result.  The mini-EVM tracks gas_used internally — we use
                // the gas_limit as an upper bound indicator here.
                // For the mini-EVM, OutOfGas already handles the gas > limit case; so any
                // Return that reaches here consumed ≤ gas_limit gas.
                if tx_data.gas_limit > MAX_GAS {
                    // Pre-screen: if requested gas_limit > cap, reject before running.
                    // (Infinite-loop contracts will hit OutOfGas above; this catches
                    //  overly greedy txs that would pass but violate the policy.)
                    return SimResult::Abort(AbortReason::GasLimitExceeded);
                }

                // Slippage guard.
                let slippage_bps = decode_slippage_bps(&data);
                if slippage_bps > MAX_SLIPPAGE_BPS {
                    return SimResult::Abort(AbortReason::SlippageExceeded);
                }

                SimResult::Ok
            }
        }
    }
}

impl Default for Screener {
    fn default() -> Self {
        Self::new()
    }
}

// ─── TEE mode: revm + Alchemy fork ───────────────────────────────────────────
//
// Compile with `--features tee` once the dependencies above are added.
// The `tee` path replaces the mini-EVM with a full revm execution against
// forked Sepolia state, giving accurate gas and slippage numbers.

/// Result of a `tee`-mode simulation (mirrors `SimResult` for the dev path).
#[cfg(feature = "tee")]
pub enum TeeSimResult {
    Ok { gas_used: u64 },
    Revert { reason: String },
    GasExceeded,
    SlippageExceeded,
}

#[cfg(feature = "tee")]
impl Screener {
    /// Simulate a settlement transaction against forked Sepolia state.
    ///
    /// # Required environment variables
    /// - `RPC_URL` — Alchemy Sepolia HTTPS endpoint.
    ///
    /// # Steps
    /// 1. Create an alloy `ProviderBuilder` with the Alchemy RPC.
    /// 2. Wrap it in `revm::db::AlloyDB` at `BlockNumberOrTag::Latest`.
    /// 3. Configure `revm::EVM` with the `settle()` calldata as a call tx.
    /// 4. Run `evm.transact()` and inspect `ExecutionResult`.
    /// 5. Decode the `Settled` event log to extract `amountOut` for slippage check.
    ///
    /// # TODO: implementation (requires adding deps listed in module doc)
    ///
    /// ```rust,ignore
    /// use alloy_provider::ProviderBuilder;
    /// use revm::{
    ///     db::AlloyDB,
    ///     primitives::{ExecutionResult, TxEnv, TransactTo, U256, Bytes, Address},
    ///     EVM,
    /// };
    ///
    /// pub async fn simulate_settlement_tee(&self, tx_data: TxData) -> SimResult {
    ///     let rpc = std::env::var("RPC_URL").expect("RPC_URL required in tee mode");
    ///     let provider = ProviderBuilder::new().on_http(rpc.parse().unwrap());
    ///
    ///     let block = provider.get_block_number().await.expect("get block");
    ///     let db = AlloyDB::new(provider, block.into());
    ///     let mut cache_db = revm::db::CacheDB::new(db);
    ///
    ///     let mut evm = EVM::new();
    ///     evm.database(&mut cache_db);
    ///     evm.env.tx = TxEnv {
    ///         caller:   Address::ZERO, // TEE wallet; fund via env or state override
    ///         transact_to: TransactTo::Call(tx_data.to.0.into()),
    ///         data:     tx_data.calldata.0.clone().into(),
    ///         gas_limit: tx_data.gas_limit,
    ///         value:    U256::ZERO,
    ///         ..Default::default()
    ///     };
    ///
    ///     match evm.transact().map(|r| r.result) {
    ///         Ok(ExecutionResult::Success { gas_used, .. }) => {
    ///             if gas_used > MAX_GAS { SimResult::Abort(AbortReason::GasLimitExceeded) }
    ///             else { SimResult::Ok }
    ///         }
    ///         Ok(ExecutionResult::Revert { output, .. }) => {
    ///             SimResult::Abort(AbortReason::ContractRevert(
    ///                 decode_revert_string(&output)
    ///             ))
    ///         }
    ///         Ok(ExecutionResult::Halt { .. }) => {
    ///             SimResult::Abort(AbortReason::GasLimitExceeded)
    ///         }
    ///         Err(e) => SimResult::Abort(AbortReason::ContractRevert(e.to_string())),
    ///     }
    /// }
    /// ```
    pub async fn simulate_settlement_tee(&self, tx_data: TxData) -> SimResult {
        // Placeholder until deps are added.
        // Falls back to dev-mode simulation so the process stays functional.
        self.simulate_settlement(tx_data)
    }
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

/// Decode a uint256 from the first 32 bytes of return data as a u64 slippage value (bps).
fn decode_slippage_bps(data: &[u8]) -> u64 {
    if data.len() < 32 {
        return 0;
    }
    let mut buf = [0u8; 8];
    buf.copy_from_slice(&data[24..32]);
    u64::from_be_bytes(buf)
}

/// Attempt to extract a human-readable revert string.
/// Supports Solidity `Error(string)` ABI encoding (selector 0x08c379a0).
fn decode_revert_string(data: &[u8]) -> String {
    const ERROR_SELECTOR: &[u8] = &[0x08, 0xc3, 0x79, 0xa0];
    if data.starts_with(ERROR_SELECTOR) && data.len() >= 68 {
        if let Ok(s) = std::str::from_utf8(&data[68..]) {
            return s.trim_end_matches('\0').to_string();
        }
    }
    if data.is_empty() {
        return "revert".to_string();
    }
    format!("revert ({}B)", data.len())
}

// ─── Tests ───────────────────────────────────────────────────────────────────

#[cfg(test)]
mod tests {
    use super::*;

    fn addr(byte: u8) -> Address {
        Address::from([byte; 20])
    }

    fn make_tx(to: Address, gas_limit: u64) -> TxData {
        TxData {
            to,
            calldata: Bytes::new(),
            value: U256::zero(),
            gas_limit,
        }
    }

    /// Returns 32 bytes of zeros → slippage = 0 bps → SimResult::Ok.
    /// Bytecode: PUSH1 0, PUSH1 0, MSTORE, PUSH1 32, PUSH1 0, RETURN
    const RETURN_ZERO_BPS: &[u8] =
        &[0x60, 0x00, 0x60, 0x00, 0x52, 0x60, 0x20, 0x60, 0x00, 0xf3];

    /// Returns 51 in the last byte of a 32-byte word → 51 bps > 50 → SlippageExceeded.
    /// Bytecode: PUSH1 51, PUSH1 0, MSTORE, PUSH1 32, PUSH1 0, RETURN
    const RETURN_51_BPS: &[u8] =
        &[0x60, 0x33, 0x60, 0x00, 0x52, 0x60, 0x20, 0x60, 0x00, 0xf3];

    /// Always reverts with empty data.
    /// Bytecode: PUSH1 0, PUSH1 0, REVERT
    const ALWAYS_REVERT: &[u8] = &[0x60, 0x00, 0x60, 0x00, 0xfd];

    /// Infinite JUMP loop — burns all gas.
    /// Bytecode: JUMPDEST, PUSH1 0, JUMP
    const INFINITE_LOOP: &[u8] = &[0x5b, 0x60, 0x00, 0x56];

    fn screener_with(addr: Address, code: &[u8]) -> Screener {
        Screener::new().with_contract(addr, code.to_vec())
    }

    #[test]
    fn sim_ok() {
        let a = addr(0xaa);
        let s = screener_with(a, RETURN_ZERO_BPS);
        assert!(matches!(s.simulate_settlement(make_tx(a, 100_000)), SimResult::Ok));
    }

    #[test]
    fn sim_slippage_exceeded() {
        let a = addr(0xbb);
        let s = screener_with(a, RETURN_51_BPS);
        assert!(matches!(
            s.simulate_settlement(make_tx(a, 100_000)),
            SimResult::Abort(AbortReason::SlippageExceeded)
        ));
    }

    #[test]
    fn sim_contract_revert() {
        let a = addr(0xcc);
        let s = screener_with(a, ALWAYS_REVERT);
        assert!(matches!(
            s.simulate_settlement(make_tx(a, 100_000)),
            SimResult::Abort(AbortReason::ContractRevert(_))
        ));
    }

    #[test]
    fn sim_gas_limit_exceeded_by_loop() {
        let a = addr(0xdd);
        let s = screener_with(a, INFINITE_LOOP);
        // gas_limit 600_000 > MAX_GAS 500_000 → loop burns gas, OutOfGas fires
        assert!(matches!(
            s.simulate_settlement(make_tx(a, 600_000)),
            SimResult::Abort(AbortReason::GasLimitExceeded)
        ));
    }

    #[test]
    fn sim_gas_limit_exceeded_by_policy() {
        let a = addr(0xee);
        // Contract succeeds but requested gas_limit > policy cap
        let s = screener_with(a, RETURN_ZERO_BPS);
        assert!(matches!(
            s.simulate_settlement(make_tx(a, 600_000)),
            SimResult::Abort(AbortReason::GasLimitExceeded)
        ));
    }

    #[test]
    fn eoa_transfer_always_ok() {
        let s = Screener::new(); // no contracts registered
        assert!(matches!(
            s.simulate_settlement(make_tx(addr(0xff), 21_000)),
            SimResult::Ok
        ));
    }

    #[test]
    fn decode_slippage_zero() {
        assert_eq!(decode_slippage_bps(&[0u8; 32]), 0);
    }

    #[test]
    fn decode_slippage_value() {
        let mut d = [0u8; 32];
        d[31] = 33;
        assert_eq!(decode_slippage_bps(&d), 33);
    }

    #[test]
    fn mini_evm_stop_returns_empty() {
        let code = &[0x00u8]; // STOP
        let outcome = MiniEvm::new(code, 100_000).run();
        assert!(matches!(outcome, EvmOutcome::Return(d) if d.is_empty()));
    }

    #[test]
    fn mini_evm_jump_loop_exhausts_gas() {
        let outcome = MiniEvm::new(INFINITE_LOOP, 600_000).run();
        assert!(matches!(outcome, EvmOutcome::OutOfGas { .. }));
    }
}
