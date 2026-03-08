# Shadow — Privacy-Preserving Payments & Dark-Pool Trading

Shadow is a privacy-first payment and trading infrastructure built on **Ethereum Sepolia**, powered by **Chainlink CRE (Compute Runtime Environment)** for decentralized off-chain execution and **World ID** for Sybil-resistant identity verification. Users can send USDC to an encrypted recipient address and place confidential dark-pool trading orders — all without revealing sensitive data on-chain.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                        World App (Next.js)                          │
│   World ID Gate → Dashboard → Shielded Transfer / Dark Pool Trade  │
└──────────────────────────────┬──────────────────────────────────────┘
                               │  MiniKit.sendTransaction
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                    Shadow Contract (Solidity)                        │
│             Ethereum Sepolia: 0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De │
│                                                                     │
│  deposit(encryptedRecipient, amount)  → emit ShieldedDeposit        │
│  placeOrder(encryptedOrder, orderId)  → emit ShieldedOrder          │
│  onReport(metadata, report)           → _processReport → payout     │
│  ccipReceive(message)                 → cross-chain payout          │
└──────────────────────────────┬──────────────────────────────────────┘
                               │  on-chain events
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│              Chainlink CRE DON (Decentralized Oracle Network)       │
│           Workflow ID: 0063619ae2feb1f1704d31196f5ec25314ce37c2a5252d9dbf28f3b827fd5f99 │
│                                                                     │
│  Trigger 1: ShieldedDeposit                                         │
│    → RunInNodeMode: AES-256-CTR decrypt recipient (BFT consensus)   │
│    → GenerateReport (ECDSA signed by DON)                           │
│    → Shadow.onReport() → _executePayout → USDC transfer             │
│                                                                     │
│  Trigger 2: ShieldedOrder                                           │
│    → RunInNodeMode: AES-256-CTR decrypt order JSON (BFT consensus)  │
│    → Confidential HTTP POST /order → dark-pool matcher              │
│    → Confidential HTTP POST /match → settlement cycle               │
└──────────────────────────────┬──────────────────────────────────────┘
                               │  Confidential HTTP (secret-injected)
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│              SDP Dark-Pool Matcher (Rust / Axum)                    │
│                   https://shadowdeploy-1.onrender.com               │
│                                                                     │
│  POST /order  → idempotent order book insertion                     │
│  POST /match  → price-time priority matching + settlement           │
│    → SdpSettlement.settle() via Flashbots relay (MEV-protected)     │
│  GET  /health → liveness probe                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Repository Structure

```
shadow/
├── frontend/                   # Next.js 15 World Mini App
│   └── src/
│       ├── app/
│       │   ├── page.tsx                  # World ID verification gate
│       │   ├── transfer/page.tsx         # Shielded USDC transfer
│       │   ├── trading/page.tsx          # Markets + dark pool orders
│       │   ├── history/page.tsx          # Real on-chain tx history
│       │   ├── profile/page.tsx          # User profile
│       │   └── api/
│       │       ├── shielded-deposit/     # AES-CTR encrypts recipient
│       │       ├── shielded-order/       # AES-CTR encrypts order JSON
│       │       ├── transactions/         # Reads ShieldedDeposit/Payout events
│       │       ├── verify-proof/         # World ID proof verification
│       │       └── rp-signature/         # IDKit RP signature
│       ├── components/
│       │   ├── Dashboard/               # Balance + real tx list
│       │   ├── Pay/                     # Standard & shielded payment
│       │   ├── BottomNav/               # Home / Trade / History
│       │   └── TradingViewWidget/       # Embedded charts
│       ├── hooks/
│       │   └── useMultiChainBalance.ts  # USDC balance across 3 chains
│       └── abi/
│           └── shadow.abi.json          # Shadow contract ABI
│
├── workflow/                   # Chainlink CRE Workflow (Go → WASM)
│   ├── vault/
│   │   ├── main.go             # Workflow entrypoint, two triggers
│   │   ├── handler.go          # ShieldedDeposit handler
│   │   ├── matcher.go          # ShieldedOrder handler
│   │   ├── crypto.go           # AES-256-CTR decrypt + ABI encode
│   │   ├── event.go            # Log parsing
│   │   ├── types.go            # Config and payload types
│   │   ├── util.go             # Hex utilities
│   │   ├── workflow.yaml       # CRE workflow manifest
│   │   ├── config.staging.json
│   │   └── config.production.json
│   ├── contracts/evm/src/
│   │   ├── AuraVault.sol       # Shadow contract source
│   │   ├── abi/shadow.abi      # Compiled ABI
│   │   └── generated/shadow/   # Go bindings (abigen)
│   ├── project.yaml            # CRE project config
│   ├── secrets.yaml            # CRE secrets manifest
│   └── go.mod
│
└── sdp/                        # Dark-Pool Matching Service
    ├── ecloud/                 # Rust Axum HTTP server
    │   ├── Cargo.toml
    │   ├── Dockerfile
    │   └── src/
    │       ├── main.rs         # Server init + routing
    │       ├── handlers.rs     # POST /order, POST /match
    │       ├── settler.rs      # SdpSettlement on-chain calls
    │       └── types.rs        # Order, Fill, MatchResponse
    ├── sdp-ecloud/
    │   └── .env                # Wallet, RPC, contract config
    └── contracts/              # SdpSettlement Foundry project
```

---

## Smart Contracts

### Shadow Contract (`AuraVault.sol`)

**Deployed:** `0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De` (Ethereum Sepolia)

The Shadow contract is the on-chain privacy layer. It inherits from two base contracts:

- **`CCIPReceiver`** — Chainlink CCIP integration for cross-chain payouts. Accepts token transfers from other chains and routes them to the decrypted recipient.
- **`ReceiverTemplate`** — Chainlink CRE integration. Accepts signed DON reports via `onReport()` and validates the workflow ID, author, and signature before processing.

```solidity
// Deposit USDC with encrypted recipient
function deposit(bytes32 _encryptedRecipient, uint256 _amount) external

// Place a confidential dark-pool order
function placeOrder(bytes calldata _encryptedOrder, bytes32 _orderId) external

// Called by CRE forwarder — same-chain payout
function onReport(bytes metadata, bytes report) external   // → _processReport

// Called by CCIP router — cross-chain payout
function ccipReceive(Any2EVMMessage message) internal      // → _ccipReceive
```

**Events:**
| Event | Emitted When |
|-------|-------------|
| `ShieldedDeposit(address sender, bytes32 encryptedRecipient, uint256 amount)` | User deposits USDC with encrypted recipient |
| `ShieldedOrder(address trader, bytes encryptedOrder, bytes32 orderId)` | User places dark-pool order |
| `ShieldedPayout(address recipient, uint256 amount)` | CRE DON routes payout to decrypted recipient |

**Key design:** The `encryptedRecipient` is a `bytes32` containing a 12-byte AES-CTR nonce followed by 20 bytes of encrypted address ciphertext. No one observing the chain can determine who receives the funds.

---

### SdpSettlement Contract

**Deployed:** `0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71` (Ethereum Sepolia)

Handles on-chain settlement of matched dark-pool orders through Uniswap V3. Only callable by the `teeWallet` (the matcher service's signing wallet).

```solidity
function settle(
    bytes32 buyOrderId,
    bytes32 sellOrderId,
    uint256 amountIn,
    uint256 amountOutMin
) external
```

Transactions are broadcast through the **Flashbots Sepolia relay** (`https://relay-sepolia.flashbots.net`) for MEV protection — matched trades cannot be front-run.

---

## Chainlink CRE Workflow

The CRE (Compute Runtime Environment) workflow is the core of Shadow's privacy model. It runs as a **WASM binary** deployed to a Chainlink DON — a decentralized network of oracle nodes that must reach Byzantine Fault Tolerant (BFT) consensus before any action is taken.

**Deployed Workflow:**
- ID: `0063619ae2feb1f1704d31196f5ec25314ce37c2a5252d9dbf28f3b827fd5f99`
- Binary: compiled Go → `wasip1` target
- Chain: Ethereum Testnet Sepolia

### Why CRE?

Traditional privacy solutions require a centralized server or a trusted execution environment (TEE) to hold decryption keys. CRE eliminates this by distributing the secret key across a DON: **no single node can decrypt data or forge a payout**. A supermajority of nodes must independently decrypt and agree on the same plaintext before any funds move.

### Trigger 1 — ShieldedDeposit Handler (`handler.go`)

Fires when `Shadow.deposit()` is called on-chain.

```
Event: ShieldedDeposit(sender, encryptedRecipient, amount)
Topic: 0xa528a5618e16311d917a25d1b6f7d83f001f50e5b4bee369286d16c83784e22a
```

**Execution flow:**

```
1. Parse log → extract encryptedRecipient (bytes32) + amount (uint256)

2. RunInNodeMode (per-node isolated execution):
   ├── AES-256-CTR decrypt encryptedRecipient
   │     Layout: [nonce:12][ciphertext:20] → 20-byte recipient address
   │     IV: nonce zero-padded to 16 bytes (Go: copy(iv[:12], nonce))
   └── Return payoutPayload{EncodedReport: abi.encode(address, uint256)}

3. BFT Consensus: ConsensusAggregationFromTags
   All nodes must agree on identical payoutPayload before proceeding

4. GenerateReport → ECDSA-signed DON report (keccak256 hashing)

5. shadow.WriteReport() → Shadow.onReport() → _processReport()
   → abi.decode(report, (address, uint256))
   → usdc.transfer(recipient, amount)
```

### Trigger 2 — ShieldedOrder Handler (`matcher.go`)

Fires when `Shadow.placeOrder()` is called on-chain.

```
Event: ShieldedOrder(trader, encryptedOrder, orderId)
Topic: 0xf6604c3bcd7473aa76b7d82d56a4453e9b64a148ca27eb7c07e9f698093b5af5
```

**Execution flow:**

```
1. Parse log → extract encryptedOrder (bytes) + orderId (bytes32)
   orderId: first 16 bytes = UUID, parsed with uuid.FromBytes()

2. RunInNodeMode (per-node isolated execution):
   ├── AES-256-CTR decrypt encryptedOrder
   │     Layout: [nonce:12][ciphertext:N] → JSON string
   └── Validate JSON + return orderConsensus{OrderJSON}

3. BFT Consensus: all nodes agree on identical decrypted order JSON

4. Inject orderId into order payload (server-side idempotency key)

5. Confidential HTTP POST /order
   ├── URL: https://shadowdeploy-1.onrender.com/order
   ├── Auth: Bearer {{.MATCHER_API_KEY}}  ← injected inside CRE enclave
   └── Body: order JSON with "id" field set to orderId UUID

6. Confidential HTTP POST /match
   └── Triggers matching cycle → fills → SdpSettlement.settle()
```

### Confidential HTTP

The workflow uses `confidentialhttp.Client` from the CRE SDK to make outbound HTTP calls where **secrets are injected inside the DON enclave** and never appear in workflow source code or node memory. The `MATCHER_API_KEY` secret is registered with:

```bash
cre secrets set MATCHER_API_KEY=<value> --target staging-settings
```

And referenced in the request as `{{.MATCHER_API_KEY}}` — resolved only at execution time within the secure enclave.

### Encryption Scheme

All sensitive data uses **AES-256-CTR** with a 12-byte random nonce:

| Field | Layout | Size |
|-------|--------|------|
| `encryptedRecipient` | `[nonce:12][address_ciphertext:20]` | 32 bytes (`bytes32`) |
| `encryptedOrder` | `[nonce:12][json_ciphertext:N]` | Variable (`bytes`) |

The 12-byte nonce is zero-padded to a 16-byte CTR IV (bytes 12–15 = `0x00`). This is consistent between the Go workflow (decryption) and the Node.js frontend API (encryption).

### Build & Deploy

```bash
cd workflow

# Build WASM binary
GOARCH=wasm GOOS=wasip1 go build -o vault/tmp.wasm ./vault/

# Dry run (no broadcast)
cre workflow simulate vault --target staging-settings

# Live deploy to DON
cre workflow deploy vault --target staging-settings
```

---

## Dark-Pool Matching Service (SDP)

**Live:** `https://shadowdeploy-1.onrender.com`

A Rust Axum HTTP server that maintains an in-memory order book and executes on-chain settlement when orders match.

### API

| Endpoint | Auth | Description |
|----------|------|-------------|
| `GET /health` | None | Liveness check |
| `POST /order` | Bearer token | Insert order (idempotent by `id`) |
| `POST /match` | Bearer token | Run matching cycle, settle fills |

### Matching Algorithm

Orders are matched using **price-time priority**:

1. Group orders by asset symbol
2. Sort buy orders descending by price, sell orders ascending by price
3. While `best_buy_price >= best_sell_price`: fill at midpoint price
4. Fill quantity = `min(buy_remaining, sell_remaining)`
5. Call `SdpSettlement.settle(buyOrderId, sellOrderId, amountIn, 0)` on-chain
6. Remove fully filled orders from the book

Settlement transactions go through the **Flashbots relay** for MEV protection.

### Order Format

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "side": "buy",
  "asset": "ETH",
  "price": 3500.00,
  "quantity": 1.5,
  "trader": "0xabc..."
}
```

---

## Frontend — World Mini App

Built with **Next.js 15**, integrated into the **World App** via MiniKit.

### Authentication Flow

1. User opens the app inside World App
2. **World ID IDKit** verification gate — proves unique human identity (Orb or Device level)
3. On success, `MiniKit.walletAddress` is stored in `localStorage`
4. All subsequent pages use the verified wallet address for balance + history

### Shielded Transfer Flow

```
User enters recipient address + amount
  → POST /api/shielded-deposit (server-side AES-256-CTR encryption)
  → Returns: { encryptedRecipient: bytes32, amountRaw: uint256 }
  → MiniKit.sendTransaction([
      usdc.approve(shadowAddress, amountRaw),
      shadow.deposit(encryptedRecipient, amountRaw)
    ])
  → Shadow contract emits ShieldedDeposit
  → CRE DON picks up event → decrypts → pays out
```

### Dark Pool Order Flow

```
User selects asset + side (buy/sell) + price + quantity
  → POST /api/shielded-order (server-side AES-256-CTR encryption)
  → Returns: { encryptedOrder: bytes, orderId: bytes32 }
  → MiniKit.sendTransaction([
      shadow.placeOrder(encryptedOrder, orderId)
    ])
  → Shadow contract emits ShieldedOrder
  → CRE DON picks up event → decrypts → relays to matcher
  → Matcher finds counterparty → settles on-chain
```

### Multi-Chain Balance

USDC balances are aggregated across three chains in real-time:

| Chain | USDC Address |
|-------|-------------|
| Ethereum Sepolia | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` |
| World Chain Sepolia | `0x79A02482A880bCE3F13e09Da970dC34db4CD24d1` |
| Arbitrum Sepolia | `0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d` |

---

## Key Addresses

| Component | Address |
|-----------|---------|
| Shadow Contract (Sepolia) | `0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De` |
| SdpSettlement Contract (Sepolia) | `0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71` |
| USDC (Sepolia) | `0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238` |
| Deployer Wallet | `0xC76F3E8e77cD40fAACf2C5F874774cAB1Ca9dB5d` |
| CRE Workflow ID | `0063619ae2feb1f1704d31196f5ec25314ce37c2a5252d9dbf28f3b827fd5f99` |
| SDP Matcher | `https://shadowdeploy-1.onrender.com` |

---

## Environment Setup

### Frontend (`frontend/.env`)

```env
NEXT_PUBLIC_APP_ID=app_69420e884d42172cf2ae55adb6939fcd
NEXT_PUBLIC_WLD_ACTION_NAME=payment-verification
NEXT_PUBLIC_WLD_RP_ID=rp_e8c2bc78ed3521e2
AUTH_SECRET=<nextauth_secret>
RP_SIGNING_KEY=<world_id_rp_signing_key>
NEXT_PUBLIC_SHADOW_CONTRACT=0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De
NEXT_PUBLIC_USDC_SEPOLIA=0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
NEXT_PUBLIC_SEPOLIA_RPC=https://eth-sepolia.g.alchemy.com/v2/<key>
DON_ENCRYPTION_KEY=<32_byte_hex_key_matching_workflow>
```

### Workflow (`workflow/secrets.yaml`)

```yaml
secretsNames:
  MATCHER_API_KEY:
    - MATCHER_API_KEY_ALL
```

```bash
cre secrets set MATCHER_API_KEY=<value> --target staging-settings
```

### SDP Matcher (`sdp/sdp-ecloud/.env`)

```env
MNEMONIC=<bip39_mnemonic>
WALLET_INDEX=3
SETTLEMENT_CONTRACT=0xB1F0214E2277c2843A9D2d90cCEAd664d19C9f71
RPC_URL=https://relay-sepolia.flashbots.net
MATCHER_API_KEY=<same_value_as_cre_secret>
PORT=8080
```

---

## How It All Connects

```
Privacy guarantee:
  On-chain:  only ciphertext is ever stored (encryptedRecipient, encryptedOrder)
  Off-chain: decryption happens inside CRE DON with BFT consensus
             — no single party sees the plaintext
  Settlement: DON signs a report → contract verifies DON signature → USDC moves

Sybil resistance:
  World ID orb verification → each human can only participate once
  Verified at the app layer before any deposit or order is accepted

MEV protection:
  Dark pool orders are invisible until matched
  Settlement transactions broadcast via Flashbots relay
  No front-running possible
```

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | Next.js 15, TypeScript, Tailwind CSS 4, MiniKit, viem |
| Identity | Worldcoin World ID (IDKit + MiniKit), NextAuth |
| Smart Contracts | Solidity 0.8.26, Foundry, OpenZeppelin, Chainlink CCIP |
| Off-chain Compute | Chainlink CRE SDK (Go → WASM), AES-256-CTR |
| Dark Pool | Rust, Axum, ethers-rs, Flashbots relay |
| Blockchain | Ethereum Sepolia, World Chain Sepolia, Arbitrum Sepolia |
