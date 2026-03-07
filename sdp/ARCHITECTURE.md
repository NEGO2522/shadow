# Shadow Dark Pool — Architecture & Usage

## What it is

A **dark pool** trading system for ETH/USDC on Sepolia. Unlike public DEXes, orders are
never visible until after they match. The TEE (Trusted Execution Environment) is the
neutral party that sees both sides and matches them privately.

---

## System Overview

```
Browser (tin/ UI)
    │
    │  POST /order   ← user submits buy or sell
    │  POST /match   ← polls every 3s for fills
    │  GET  /health  ← liveness check
    ▼
sdp-ecloud (Axum, port 3000)
    │
    ├── MatchEngine          ← price-time priority order book (in-memory)
    ├── PreScreener          ← simulates the settlement tx before broadcast
    └── SettlementRelayer    ← signs + sends via Flashbots Sepolia
            │
            ▼
    SdpSettlement.sol        ← on-chain: calls Uniswap v3, emits Settled event
            │
            ▼
    Uniswap v3 (Sepolia)     ← actual token swap, MEV-protected via Flashbots
```

---

## How a Trade Flows End-to-End

```
1. Alice submits BUY  1.0 ETH @ $2841
   POST /order { side:"buy", price:284100, quantity:1000000, trader_pubkey:"0x..." }
                         ↑ cents          ↑ units of 10⁻⁶ ETH

2. Bob submits SELL 1.0 ETH @ $2841
   POST /order { side:"sell", price:284100, quantity:1000000, ... }

3. UI polls POST /match every 3s
   MatchEngine runs price-time priority crossing:
     best_bid (284100) >= best_ask (284100) → FILL
     exec_price = ask price (resting order)

4. PreScreener simulates the settle() calldata against the contract
   If slippage > 0.5% or gas > 500k → Abort, order NOT broadcast

5. SettlementRelayer broadcasts to Flashbots Sepolia:
   settle(buyId, sellId, amountIn=2_841_000_000, amountOutMin=995_000)
     amountIn     = 284100 × 1_000_000 / 100 = $2841 USDC (6-dec)
     amountOutMin = 1_000_000 × 995 / 1000   = 0.995 ETH (0.5% slippage)

6. Transaction lands atomically via Flashbots → not front-runnable
```

---

## Crate Structure

| Crate | Path | Role |
|---|---|---|
| `sdp-shared` | `ecloud/shared` | Shared types: `LimitOrder`, `MatchResult`, `SimResult`, `Side` |
| `sdp-matching-engine` | `ecloud/matching-engine` | Price-time priority BTreeMap order book; `execute_match()` returns fills |
| `sdp-pre-screener` | `ecloud/pre-screener` | Dev: mini pure-Rust EVM. TEE: revm + Alchemy fork (stub) |
| `sdp-kms` | `ecloud/kms` | Wraps `LocalWallet`; in TEE mode becomes hardware-backed signing |
| `sdp-relayer` | `ecloud/relayer` | Builds ABI-encoded `settle()` calldata, sends via Flashbots |
| `sdp-ecloud` | `sdp-ecloud/` | Axum binary wiring all crates together |
| `contracts/` | `contracts/` | Foundry project: `SdpSettlement.sol` + deploy script + tests |

---

## Price Encoding

Everything flows in two scales:

```
UI input    → sdpClient.ts multiplies ×100  → backend (cents + eth6)
Backend     → FillInfo divides ÷100         → UI display

price:    $2841.00  →  284100   (u64 cents)
quantity: 1.0 ETH   →  1000000  (u64 units of 10⁻⁶ ETH)

On-chain:
  amountIn     = 284100 × 1000000 / 100 = 2_841_000_000  (USDC, 6 dec)
  amountOutMin = 1000000 × 995 / 1000   =       995_000  (ETH, 6 dec)
```

> **Note:** WETH on-chain uses 18 decimals. Before mainnet, multiply `amountOutMin` by
> `10^12` to convert from ETH-6 to wei.

---

## Dark Pool Properties

| Property | How it's enforced |
|---|---|
| No visible order book | MatchEngine state is never exposed via API |
| Anonymized fills | `/match` returns UUIDs, not trader addresses |
| MEV protection | Flashbots relay bypasses the public mempool |
| TEE attestation | Only the TEE wallet can call `settle()` on-chain (`teeWallet` immutable) |

---

## Running Locally (No On-Chain Settlement)

**Terminal 1 — backend**

```bash
cd /home/daiwi/trade/sdp/sdp-ecloud
# Comment out SETTLEMENT_CONTRACT in .env to skip relay
cargo run
# Listening on 0.0.0.0:3000
```

**Terminal 2 — UI**

```bash
cd /home/daiwi/trade/tin
npm run dev
# http://localhost:3001
```

With `SETTLEMENT_CONTRACT` empty the server still matches and pre-screens — it just skips
the on-chain broadcast. The `screener` field in the `/match` response will read
`"ok (relay skipped: SETTLEMENT_CONTRACT not set)"`.

---

## Running With Real On-Chain Settlement

The contract is already deployed. Just start the backend with `.env` intact:

```bash
cd /home/daiwi/trade/sdp/sdp-ecloud
cargo run
```

The `.env` has `SETTLEMENT_CONTRACT` set to the live address. The TEE wallet must hold
Sepolia USDC and ETH for each settlement call.

### Re-deploying

```bash
cd /home/daiwi/trade/sdp/contracts

forge install OpenZeppelin/openzeppelin-contracts
forge install Uniswap/v3-periphery
forge install Uniswap/v3-core

export PRIVATE_KEY=0x...
export TEE_WALLET=0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d
export RPC_URL=https://eth-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW

forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $RPC_URL \
  --private-key $PRIVATE_KEY \
  --broadcast
```

---

## API Reference

### `POST /order`

Submit a limit order to the dark pool.

```json
{
  "side":          "buy",        // "buy" | "sell"
  "price":         284100,       // u64 cents  ($2841.00)
  "quantity":      1000000,      // u64 ETH-6  (1.000000 ETH)
  "trader_pubkey": "0xabcd..."   // hex-encoded trader identity
}
```

**Response**

```json
{ "id": "550e8400-...", "status": "queued" }
```

---

### `POST /match`

Trigger a matching cycle. Returns all new fills since the last cycle.

```json
{
  "fills": 1,
  "results": [
    {
      "buy_order_id":  "550e8400-...",
      "sell_order_id": "6ba7b810-...",
      "price":         284100,
      "quantity":      1000000,
      "tx_hash":       "0xdeadbeef...",   // null if relay skipped
      "screener":      "ok"
    }
  ]
}
```

---

### `GET /health`

Returns `200 ok` when the process is live.

---

## Deployed Contracts (Sepolia)

| Contract | Address |
|---|---|
| **SdpSettlement** | `0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71` |
| TEE wallet | `0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d` |
| Uniswap v3 Router | `0xE592427A0AEce92De3Edee1F18E0157C05861564` |
| WETH | `0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14` |
| USDC | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` |
| Chainlink ETH/USD | `0x694AA1769357215DE4FAC081bf1f309aDC325306` |

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `MNEMONIC` | Yes | BIP-39 mnemonic; index 0 = TEE relayer wallet |
| `SETTLEMENT_CONTRACT` | No | Deployed `SdpSettlement` address; leave empty to skip relay |
| `RPC_URL` | No | Alchemy Sepolia RPC; defaults to Flashbots Sepolia relay |

---

## TEE Feature (Pre-Screener)

The `tee` feature flag on `sdp-pre-screener` gates a full `revm` simulation against
forked Sepolia state. When enabled, `simulate_settlement_tee()` replaces the mini-EVM
with an exact replay of what would happen on-chain.

To activate, add to `ecloud/pre-screener/Cargo.toml`:

```toml
[features]
tee = ["revm", "alloy-provider", "alloy-network", "alloy-rpc-types"]

[dependencies]
revm            = { version = "14", default-features = false, features = ["std"] }
alloy-provider  = { version = "0.9", features = ["reqwest"] }
alloy-network   = "0.9"
alloy-rpc-types = "0.9"
```

Then build with `cargo run --features tee`.

The full implementation pseudocode lives in `ecloud/pre-screener/src/lib.rs`.
