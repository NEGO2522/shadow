# Shadow Dark Pool — Backend

A trustless dark-pool matching engine for ETH/USDC on Sepolia, running inside an Intel TDX Trusted Execution Environment (TEE) on [EigenCompute](https://www.eigencompute.com/). Orders are matched off-chain with full price-time priority, pre-screened via on-chain simulation, and settled through a Uniswap v3 swap broadcast over the Flashbots Sepolia relay.

---

## Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Repository Layout](#repository-layout)
- [Crate Reference](#crate-reference)
- [Smart Contract](#smart-contract)
- [API Reference](#api-reference)
- [Price & Quantity Encoding](#price--quantity-encoding)
- [Trade Lifecycle](#trade-lifecycle)
- [Environment Variables](#environment-variables)
- [Local Development](#local-development)
- [Docker](#docker)
- [EigenCompute Deployment](#eigencompute-deployment)
- [Contract Deployment](#contract-deployment)
- [Deployed Addresses (Sepolia)](#deployed-addresses-sepolia)
- [Security Model](#security-model)

---

## Overview

Shadow is a **dark pool**: orders are never visible on-chain until they settle. Traders submit signed limit orders over HTTPS to the TEE; the TEE matches them internally, simulates the resulting settlement transaction against a live Sepolia fork, and then broadcasts via Flashbots so miners cannot front-run or sandwich the swap.

Key properties:

| Property | Mechanism |
|---|---|
| Pre-trade privacy | Orders held in TEE memory only — never broadcast |
| MEV protection | Settlement sent via Flashbots Sepolia relay |
| Trustless access control | `settle()` reverts for any caller that is not the TEE wallet |
| Verifiable execution | EigenCompute Intel TDX attestation |
| Slippage guard | Pre-screener rejects fills with >0.5% slippage before relay |

---

## System Architecture

```
  Traders
    │  POST /order
    ▼
┌──────────────────────────────────────────────────────┐
│              sdp-ecloud  (Axum HTTP, port 3000)       │
│                                                        │
│  ┌──────────────┐   ┌─────────────────┐               │
│  │ MatchEngine  │   │   Screener      │               │
│  │ (BTreeMap    │──▶│ (mini-EVM dev / │               │
│  │  price-time) │   │  revm tee)      │               │
│  └──────────────┘   └────────┬────────┘               │
│                               │ SimResult::Ok          │
│                               ▼                        │
│                    ┌─────────────────────┐             │
│                    │ SettlementRelayer   │             │
│                    │ (encode calldata,   │             │
│                    │  sign, send bundle) │             │
│                    └──────────┬──────────┘             │
│         Intel TDX TEE         │                        │
└───────────────────────────────┼────────────────────────┘
                                │ Flashbots bundle
                                ▼
                    ┌─────────────────────┐
                    │  SdpSettlement.sol  │  Sepolia
                    │  (Uniswap v3 swap,  │
                    │   TEE-only access)  │
                    └─────────────────────┘
```

### Component Roles

**`sdp-ecloud`** — the binary. Owns shared state (`MatchEngine`, `Screener`, `SettlementRelayer`) and exposes three HTTP endpoints. Started by EigenCompute; `MNEMONIC` is injected by the KMS at boot so the private key never appears in config files or logs.

**`sdp-matching-engine`** — pure in-memory order book. Bids stored in a descending `BTreeMap` (highest price first), asks in ascending order. `execute_match()` runs a crossing loop and returns `Vec<MatchResult>` with full and partial fill support.

**`sdp-pre-screener`** — simulates each settlement transaction before broadcast. In dev mode this is a minimal pure-Rust bytecode interpreter. In tee mode (feature flag) it forks Sepolia via Alchemy and runs the real transaction through `revm`. Three abort guards: gas limit, contract revert, slippage >50 bps.

**`sdp-relayer`** — ABI-encodes the `settle(bytes32,bytes32,uint256,uint256)` calldata, signs the transaction with the TEE wallet, and sends a Flashbots bundle targeting blocks `n+1` through `n+10`.

**`sdp-kms`** — thin wallet abstraction. Dev: `MnemonicBuilder` from `MNEMONIC` env var at BIP-44 index `WALLET_INDEX`. TEE: hardware-backed key from EigenCompute KMS (same interface, different backend).

**`sdp-shared`** — shared types (`LimitOrder`, `MatchResult`, `Side`, `SimResult`, `AbortReason`) with no external dependencies.

---

## Repository Layout

```
sdp/
├── Dockerfile                   # Multi-stage build (build context = sdp/)
├── ARCHITECTURE.md              # Detailed architecture notes
│
├── ecloud/                      # Cargo workspace (library crates)
│   ├── Cargo.toml               # workspace root
│   ├── shared/                  # Core types (LimitOrder, MatchResult, …)
│   ├── matching-engine/         # Price-time priority order book
│   ├── pre-screener/            # Settlement tx simulation (mini-EVM / revm)
│   ├── relayer/                 # Flashbots bundle builder & broadcaster
│   └── kms/                     # BIP-39 wallet abstraction (dev + TEE)
│
├── sdp-ecloud/                  # Binary crate — Axum HTTP server
│   ├── Cargo.toml
│   ├── src/main.rs
│   ├── Dockerfile               # Alternate path (build context = sdp/)
│   ├── .env.example
│   └── README.md
│
├── contracts/                   # Foundry project
│   ├── foundry.toml
│   ├── src/SdpSettlement.sol    # Settlement contract
│   ├── script/Deploy.s.sol      # Deployment script
│   └── test/SdpSettlement.t.sol # Unit tests
│
└── docs/                        # Additional documentation
```

---

## Crate Reference

### `sdp-shared`

Core types shared across all crates.

```rust
pub enum Side { Buy, Sell }

pub struct LimitOrder {
    pub id:            Uuid,
    pub side:          Side,
    pub price:         u64,   // cents (e.g. 320000 = $3200.00)
    pub quantity:      u64,   // units of 10⁻⁶ ETH (1_000_000 = 1 ETH)
    pub timestamp:     u64,   // Unix ms — determines time priority
    pub trader_pubkey: String,
}

pub struct MatchResult {
    pub buy_order_id:  Uuid,
    pub sell_order_id: Uuid,
    pub price:         u64,
    pub quantity:      u64,
    pub timestamp:     u64,
}

pub enum SimResult { Ok, Abort(AbortReason) }
pub enum AbortReason {
    SlippageExceeded,
    ContractRevert(String),
    GasLimitExceeded,
}
```

### `sdp-matching-engine`

Price-time priority crossing engine.

```rust
let mut engine = MatchEngine::new();
engine.add_order(order);
let fills: Vec<MatchResult> = engine.execute_match();
engine.cancel_order(uuid);

// TEE: returns hardware attestation; dev: returns dummy bytes
let attestation: Vec<u8> = engine.attest();
```

**Matching rules:**
- A buy and sell cross when `buy.price >= sell.price`
- Fill price = the resting order's price (maker price)
- Partial fills are supported; the unfilled remainder stays on the book
- FIFO within each price level

### `sdp-pre-screener`

```rust
pub struct TxData {
    pub to:        Address,
    pub calldata:  Vec<u8>,
    pub value:     U256,
    pub gas_limit: u64,
}

let screener = Screener::new();
match screener.simulate_settlement(tx) {
    SimResult::Ok              => { /* safe to relay */ }
    SimResult::Abort(reason)   => { /* drop or log */ }
}
```

**Abort conditions:**
| Condition | Limit |
|---|---|
| Gas limit exceeded | 500 000 gas |
| Slippage | > 50 bps (0.5%) |
| Contract revert | any non-zero revert reason |

### `sdp-relayer`

```rust
pub const FLASHBOTS_SEPOLIA_RPC: &str = "https://relay-sepolia.flashbots.net";
pub const SEPOLIA_CHAIN_ID: u64 = 11_155_111;
pub const SETTLE_GAS_LIMIT: u64 = 350_000;

let relayer = SettlementRelayer::new(&rpc_url, wallet);
let tx_hash: H256 = relayer.relay_match(fill).await?;

// ABI encoding
pub fn encode_settlement_calldata(result: &MatchResult) -> Vec<u8>;
```

**`settle()` ABI encoding:**

```
selector  0x9e828403  keccak256("settle(bytes32,bytes32,uint256,uint256)")
slot 0    buyOrderId  (UUID → bytes32)
slot 1    sellOrderId (UUID → bytes32)
slot 2    amountIn    price_cents × quantity_eth6 ÷ 100  (USDC 6-dec)
slot 3    amountOutMin quantity_eth6 × 995 ÷ 1000        (0.5% slippage)
```

Total calldata: 4 + 4 × 32 = **132 bytes**.

### `sdp-kms`

```rust
let wallet = AppWallet::new_from_env(); // reads MNEMONIC + WALLET_INDEX
println!("{}", wallet.address());       // 0xC76F3E8e77cD...

let signed_tx: Bytes = wallet.sign_transaction(tx).await?;
let inner: LocalWallet = wallet.inner();
```

---

## Component Deep Dive

### What `sdp/` Is

The entire backend of the Shadow dark pool — everything from accepting orders off the internet to settling a real token swap on Sepolia. It is designed to run trustlessly inside an **Intel TDX Trusted Execution Environment** on EigenCompute, which is the core point: no one — not even you as the operator — can see the order book or tamper with the matching logic while it's running.

---

### `ecloud/shared` — The type layer

Nothing runs without this. Every other crate imports it. It defines the two core data types that flow through the whole pipeline:

```
LimitOrder  { id, side, price (cents), quantity (micro-ETH), timestamp, trader_pubkey }
MatchResult { buy_order_id, sell_order_id, price, quantity, timestamp }
SimResult   { Ok | Abort(SlippageExceeded | ContractRevert | GasLimitExceeded) }
```

Price is always **u64 cents** (`284100` = $2841.00). Quantity is always **u64 micro-ETH** (`1_000_000` = 1 ETH). These two scales never change across the entire system — the UI multiplies by 100 before sending, divides by 100 before displaying.

---

### `ecloud/matching-engine` — The order book

A pure in-memory order book using two `BTreeMap`s:
- **Bids**: keyed by `Reverse(price)` so the highest bid is always first
- **Asks**: keyed by `price` ascending so the lowest ask is always first

`execute_match()` runs a crossing loop: while `best_bid.price >= best_ask.price`, take both sides, fill as much quantity as both have, put remainders back. Returns a `Vec<MatchResult>`. Full and partial fills both work.

**The order book state is never exposed through the API.** `GET /orders` doesn't exist. That's the dark pool property — pre-trade opacity.

Also has an `attest()` method: in dev mode returns dummy bytes; in TEE mode returns an Intel TDX hardware attestation quote that cryptographically proves the binary running is the one built from your repo.

---

### `ecloud/pre-screener` — The simulation guard

Before any fill is broadcast on-chain, the pre-screener simulates the settlement transaction and can **kill it** if something would go wrong. Three abort conditions:

| Condition | Threshold |
|---|---|
| Gas usage | > 500 000 |
| Slippage | > 50 bps (0.5%) |
| Contract revert | any revert reason |

**Dev mode** — runs a tiny hand-rolled Rust bytecode interpreter: STOP, MSTORE, JUMP, JUMPDEST, PUSH1, RETURN, REVERT. Enough to catch obvious reverts without needing an actual RPC.

**TEE mode** (`--features tee`) — replaces the mini-EVM with `revm` forked against live Sepolia via Alchemy RPC. Every settlement is simulated against the exact current on-chain state before the relayer ever touches it.

---

### `ecloud/relayer` — The Flashbots broadcaster

Takes a `MatchResult`, produces 132 bytes of ABI-encoded calldata, signs it with the TEE wallet, and sends a **Flashbots bundle** targeting blocks `n+1` through `n+10`.

The ABI encoding for `settle(bytes32, bytes32, uint256, uint256)`:

```
selector    = 0x9e828403   (keccak256 of the function signature)
slot 0      = buyOrderId   (UUID → bytes32)
slot 1      = sellOrderId  (UUID → bytes32)
slot 2      = amountIn     = price_cents × quantity_eth6 ÷ 100  →  USDC (6 dec)
slot 3      = amountOutMin = quantity_eth6 × 995 ÷ 1000         →  0.5% slippage floor
```

**Why Flashbots**: the settlement tx never enters the public mempool. Searchers and validators cannot see it until the block is sealed — no front-running, no sandwiching.

---

### `ecloud/kms` — The wallet abstraction

Reads `MNEMONIC` from env and derives the signing key using BIP-44 at index `WALLET_INDEX`. The same interface works in both modes:

- **Dev** — `MnemonicBuilder` from `ethers`; key derived at runtime, lives in process memory
- **TEE** — same interface, but the EigenCompute KMS injects `MNEMONIC` into TDX sealed memory at boot; the host OS never sees it, you never see it, it's gone when the instance terminates

This is what makes the TEE wallet trustless: traders can verify (via attestation) that key management is happening inside the enclave and cannot be exfiltrated.

---

### `sdp-ecloud` — The binary that wires it all

An Axum HTTP server on `0.0.0.0:3000` (EigenCompute requires `0.0.0.0`). Shared state is an `Arc<AppState>` holding a `Mutex<MatchEngine>`, a `Screener`, and a `SettlementRelayer`.

**`POST /order`** — validates `side`, creates a `LimitOrder` with a new UUID and current timestamp, pushes to the engine. Returns the UUID.

**`POST /match`** — the main loop. Calls `execute_match()`, then for each fill: builds the `TxData`, calls `screener.simulate_settlement()`, and if `SimResult::Ok` and `SETTLEMENT_CONTRACT` is set, calls `relayer.relay_match()`. Returns every fill's outcome including `tx_hash` and screener status.

**`GET /health`** — `200 ok`. Used by EigenCompute as a liveness probe.

---

### `outside/` — Planned external components

Three stub directories — `api/`, `dashboard/`, `event-listener/` — for planned components outside the TEE boundary: a public REST API gateway, an operator monitoring dashboard, and an on-chain event listener for indexing `Settled` events. Not yet implemented.

---

### EigenCompute's Role

Without EigenCompute this is just a centralized server you have to trust. With it:

| Guarantee | How |
|---|---|
| Build is reproducible | Anyone can verify the deployed binary was compiled from a specific audited commit |
| Memory is sealed | Order book, private key, and MNEMONIC live in TDX-encrypted memory the host cannot read |
| Key injection is trustless | MNEMONIC flows from env file → EigenCompute KMS → sealed TEE memory, never appearing elsewhere |
| Execution is attestable | `ecloud compute build verify <digest>` cryptographically proves the TEE is running exactly what it claims |

```bash
# Verify a running deployment
ecloud compute build verify <build-digest-from-app-info>
```

---

## Smart Contract

**`SdpSettlement.sol`** — deployed on Sepolia at [`0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71`](https://sepolia.etherscan.io/address/0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71)

```solidity
function settle(
    bytes32 buyOrderId,
    bytes32 sellOrderId,
    uint256 amountIn,     // USDC (6 decimals)
    uint256 amountOutMin  // WETH (18 decimals, 0.5% slippage floor)
) external;
```

**What it does:**
1. Reverts with `NotTee()` if `msg.sender != teeWallet`
2. Pulls `amountIn` USDC from the TEE wallet via `transferFrom`
3. Approves the Uniswap v3 Router
4. Executes `exactInputSingle` (USDC → WETH, 0.05% fee tier)
5. Emits `Settled(buyOrderId, sellOrderId, amountIn, amountOut, block.number)`

**Immutables (Sepolia):**

| Name | Address |
|---|---|
| Uniswap v3 Router | `0xE592427A0AEce92De3Edee1F18E0157C05861564` |
| WETH | `0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14` |
| USDC | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` |
| TEE wallet | `0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d` |

---

## API Reference

All endpoints served on `0.0.0.0:3000`. CORS is open (`*`).

### `POST /order`

Submit a limit order to the dark pool.

**Request:**
```json
{
  "side":          "buy",
  "price":         320000,
  "quantity":      1000000,
  "trader_pubkey": "0xabc..."
}
```

| Field | Type | Description |
|---|---|---|
| `side` | `"buy"` \| `"sell"` | Order direction |
| `price` | `u64` | Limit price in **cents** (e.g. `320000` = $3,200.00) |
| `quantity` | `u64` | Amount in **micro-ETH** (1 000 000 = 1 ETH) |
| `trader_pubkey` | `string` | Trader's Ethereum address or public key |

**Response `200`:**
```json
{ "id": "550e8400-e29b-41d4-a716-446655440000", "status": "queued" }
```

---

### `POST /match`

Trigger one matching cycle. Crosses all currently crossable orders, pre-screens each fill, and relays those that pass.

**Request:** empty body

**Response `200`:**
```json
{
  "fills": 2,
  "results": [
    {
      "buy_order_id":  "...",
      "sell_order_id": "...",
      "price":         320000,
      "quantity":      500000,
      "tx_hash":       "0xabc...",
      "screener":      "ok"
    }
  ]
}
```

| `screener` value | Meaning |
|---|---|
| `"ok"` | Simulation passed, relay succeeded |
| `"ok (relay skipped: SETTLEMENT_CONTRACT not set)"` | Simulation-only mode |
| `"aborted: SlippageExceeded"` | Fill dropped — slippage > 0.5% |
| `"relay error: ..."` | Relayer RPC error |

---

### `GET /health`

Liveness probe. Returns `200 ok`.

---

## Price & Quantity Encoding

```
price    u64  cents          320000  → $3,200.00
quantity u64  micro-ETH      1000000 → 1.000000 ETH

amountIn (USDC 6-dec) = price_cents × quantity_eth6 ÷ 100
   e.g. 320000 × 500000 ÷ 100 = 1_600_000_000  (1600 USDC)

amountOutMin (WETH 18-dec) = quantity_eth6 × 995 ÷ 1000
   e.g. 500000 × 995 ÷ 1000 = 497_500 micro-ETH  (0.5% slippage)
```

The UI (`tin/`) multiplies UI prices by 100 before sending and divides by 100 when displaying.

---

## Trade Lifecycle

```
1. Trader  →  POST /order { side, price, quantity, trader_pubkey }
              Server assigns UUID, pushes to MatchEngine

2. Cron / UI  →  POST /match
              MatchEngine.execute_match() → Vec<MatchResult>

3. For each fill:
   a. build_tx_data() → TxData { to: SdpSettlement, calldata, gas_limit: 350k }
   b. Screener.simulate_settlement(tx) → SimResult
      • Dev:  mini-EVM bytecode run
      • TEE:  revm fork against live Sepolia
   c. If SimResult::Ok && SETTLEMENT_CONTRACT set:
        SettlementRelayer.relay_match(fill) → H256
        → encodes calldata (4 + 128 bytes)
        → signs with TEE wallet (KMS)
        → sends Flashbots bundle for blocks n+1 … n+10
   d. Response includes tx_hash + screener status

4. SdpSettlement.settle() executes on-chain:
   → verifies msg.sender == teeWallet
   → transfers USDC in
   → Uniswap v3 exactInputSingle USDC → WETH
   → emits Settled event
```

---

## Environment Variables

Create `sdp-ecloud/.env` (copy from `.env.example`):

```env
# BIP-39 mnemonic — in production, injected by EigenCompute KMS at boot
MNEMONIC="..."

# BIP-44 derivation index (0-based). Increment until wallet address matches funded TEE address.
WALLET_INDEX=3

# Deployed SdpSettlement contract (Sepolia). Leave empty for simulation-only mode.
SETTLEMENT_CONTRACT=0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71

# Alchemy Sepolia endpoint (also used by pre-screener in TEE mode for fork simulation).
RPC_URL=https://eth-sepolia.g.alchemy.com/v2/<key>
```

> **Note:** In production the EigenCompute KMS injects `MNEMONIC` automatically. You only need to set the other three in the EigenCompute environment config.

---

## Local Development

**Prerequisites:** Rust stable, Cargo

```bash
cd /path/to/sdp/sdp-ecloud
cp .env.example .env
# Edit .env with real values

# Run with Hardhat test mnemonic (simulation-only, no on-chain relay)
MNEMONIC="test test test test test test test test test test test junk" cargo run

# Server starts on http://localhost:3000
curl http://localhost:3000/health        # → ok

# Submit a buy order ($3200, 1 ETH)
curl -X POST http://localhost:3000/order \
  -H 'Content-Type: application/json' \
  -d '{"side":"buy","price":320000,"quantity":1000000,"trader_pubkey":"0xabc"}'

# Submit a matching sell order
curl -X POST http://localhost:3000/order \
  -H 'Content-Type: application/json' \
  -d '{"side":"sell","price":319000,"quantity":1000000,"trader_pubkey":"0xdef"}'

# Trigger matching
curl -X POST http://localhost:3000/match
```

**Run all tests:**
```bash
cd sdp/ecloud
cargo test --workspace
```

---

## Docker

Build context **must** be the `sdp/` parent directory:

```bash
# From repo root:
docker build \
  -f sdp/sdp-ecloud/Dockerfile \
  -t sdp-ecloud:latest \
  sdp/

# Run with env file
docker run --rm -p 3000:3000 \
  --env-file sdp/sdp-ecloud/.env \
  sdp-ecloud:latest
```

The runtime image is `debian:bookworm-slim` with CA certificates copied from the builder stage (no `apt-get` at runtime — required for EigenCompute restricted network environments).

---

## EigenCompute Deployment

`sdp-ecloud` is deployed as a verifiable build on EigenCompute (EigenLayer's TEE compute network). The build is reproducible from a specific git commit, and EigenCompute generates an Intel TDX attestation that anyone can verify.

### Prerequisites

```bash
npm install -g @layr-labs/ecloud-cli
ecloud auth login    # or export ECLOUD_PRIVATE_KEY=0x...
```

### First deploy

```bash
# Get the commit you want to deploy
COMMIT=$(git rev-parse HEAD)

ecloud compute app deploy \
  --name        sdp-ecloud \
  --verifiable \
  --repo        https://github.com/daiwikmh/shadow \
  --commit      $COMMIT \
  --build-context      sdp \
  --build-dockerfile   sdp-ecloud/Dockerfile \
  --env-file    sdp/sdp-ecloud/.env \
  --environment sepolia \
  --private-key 0x<your-deployer-key>
```

This triggers a verifiable build on EigenCompute:
1. EigenCompute clones the repo at `$COMMIT`
2. Builds the Docker image from `sdp/sdp-ecloud/Dockerfile` with context `sdp/`
3. Provisions a TDX instance and injects secrets from the env file via KMS
4. Returns a build digest you can verify with `ecloud compute build verify`

### Upgrade (new commit)

```bash
ecloud compute app upgrade \
  --verifiable \
  --repo      https://github.com/daiwikmh/shadow \
  --commit    $(git rev-parse HEAD) \
  --build-context    sdp \
  --build-dockerfile sdp-ecloud/Dockerfile \
  --env-file  sdp/sdp-ecloud/.env
```

### Useful commands

```bash
# List deployments
ecloud compute app list

# Detailed app info (instance IP, status, build digest)
ecloud compute app info

# Stream live logs
ecloud compute app logs --follow

# Verify the build provenance
ecloud compute build verify <build-digest>

# Manage environment secrets
ecloud compute environment --help

# Stop / start (preserves instance)
ecloud compute app stop
ecloud compute app start

# Terminate permanently
ecloud compute app terminate
```

---

## Contract Deployment

Requires [Foundry](https://book.getfoundry.sh/). Install OpenZeppelin + Uniswap remappings first:

```bash
cd sdp/contracts
forge install OpenZeppelin/openzeppelin-contracts
forge install Uniswap/v3-periphery
forge install Uniswap/v3-core
```

**Deploy to Sepolia:**

```bash
export RPC_URL=https://eth-sepolia.g.alchemy.com/v2/<key>
export PRIVATE_KEY=0x<deployer-key>
export TEE_WALLET=0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d
export ETHERSCAN_API_KEY=<key>

forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast \
  --verify
```

After deployment, set `SETTLEMENT_CONTRACT` in `sdp-ecloud/.env` and update the EigenCompute environment:

```bash
ecloud compute environment set SETTLEMENT_CONTRACT=<new-address>
```

**Run tests:**

```bash
# Unit tests (no fork required)
forge test

# Fork tests (requires Sepolia RPC)
forge test --fork-url $RPC_URL -vvv
```

---

## Deployed Addresses (Sepolia)

| Name | Address |
|---|---|
| SdpSettlement | [`0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71`](https://sepolia.etherscan.io/address/0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71) |
| TEE Wallet | [`0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d`](https://sepolia.etherscan.io/address/0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d) |
| Uniswap v3 Router | `0xE592427A0AEce92De3Edee1F18E0157C05861564` |
| WETH (Sepolia) | `0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14` |
| USDC (Sepolia) | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` |
| Deployment block | 2026-02-27 |

---

## Security Model

### Threat model

| Threat | Mitigation |
|---|---|
| Miner/validator sees order before settlement | Flashbots relay — tx is only revealed to the winning validator at inclusion time |
| Operator reads order book | Orders in TEE memory; Intel TDX memory encryption prevents host OS access |
| Operator replaces binary | EigenCompute verifiable build — anyone can verify the build digest matches the source commit |
| Slippage / bad price | Pre-screener rejects fills with >0.5% slippage before relay |
| Unauthorized `settle()` call | `SdpSettlement` reverts with `NotTee()` for any caller other than the TEE wallet |
| Private key exfiltration | Key never stored on disk; KMS injects `MNEMONIC` at boot into TEE secure memory only |
| Front-running after bundle reveal | Flashbots bundles are atomic; partial inclusion is not possible |

### Attestation

The `MatchEngine::attest()` method returns a TEE attestation quote in tee-feature mode. This allows counterparties to verify that:
1. The binary running is built from a specific audited commit
2. The environment is a genuine Intel TDX enclave
3. The matching logic has not been tampered with

Verify a running deployment:

```bash
ecloud compute build verify <build-digest-from-app-info>
```
