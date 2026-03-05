# AuraVault — Shielded Pay-by-Human Workflow

A Chainlink CRE workflow that enables **privacy-preserving, humanity-verified USDC payments** across chains. Users deposit USDC with an encrypted recipient address; a Decentralised Oracle Network (DON) verifies the depositor is human via World ID inside a TEE, decrypts the recipient, and routes the payout — either locally on Sepolia or cross-chain via CCIP.

---

## Directory Layout

```
workflow/
├── vault/
│   └── main.go               # CRE WASM handler (compiled to wasip1)
├── contracts/
│   └── evm/
│       └── src/
│           └── AuraVault.sol # Solidity vault + CCIP receiver
├── project.yaml              # CRE CLI target settings (staging / production)
├── secrets.yaml              # Secret names consumed from DON secret store
├── go.mod / go.sum           # Go module (cre-sdk-go v1.2.0)
└── .env                      # Local dev environment variables (not committed)
```

---

## How It Works

```
User
 │  deposit(encryptedRecipient, amount)
 ▼
AuraVault.sol  ──► emit ShieldedDeposit(sender, encryptedRecipient, amount)
                            │
                            │  (Chainlink CRE log trigger)
                            ▼
                     DON (5-of-9 nodes)
                            │
              ┌─────────────┴──────────────┐
              │   Each node (TEE boundary) │
              │  1. POST /verify → WorldID │  humanity check
              │  2. AES-GCM decrypt blob   │  reveals recipient + destChain
              └─────────────┬──────────────┘
                            │  ConsensusIdenticalAggregation
                            ▼
               destChain == Sepolia?
              ┌──── yes ────┤──── no ────────────────────┐
              ▼             │                            ▼
  forwarder.onReport()      │          ccipClient.SendMessage()
  AuraVault.onReport()      │          AuraVault._ccipReceive()
  usdc.transfer(recipient)  │          usdc.transfer(recipient)
  emit ShieldedPayout       │          emit ShieldedPayout
```

---

## Components

### `AuraVault.sol`

| Function | Visibility | Purpose |
|---|---|---|
| `deposit(bytes32, uint256)` | `external` | Accepts USDC + encrypted recipient; emits `ShieldedDeposit` |
| `onReport(bytes)` | `external` | Called by CRE forwarder for **same-chain** payouts; ABI-decodes `(address, uint256)` |
| `_ccipReceive(Any2EVMMessage)` | `internal override` | Called by CCIP router for **cross-chain** payouts; decodes recipient from `message.data` |

Events: `ShieldedDeposit(address indexed sender, bytes32 encryptedRecipient, uint256 amount)` and `ShieldedPayout(address indexed recipient, uint256 amount)`.

The `Client` library in the file is a minimal local stub. In production, replace it with:
```solidity
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
```

---

### `vault/main.go`

Compiled to WASM (`GOARCH=wasm GOOS=wasip1`) and deployed to the CRE DON.

**Flow inside `onShieldedDeposit`:**

1. **Parse log** — decodes `ShieldedDeposit` topics and data into `ShieldedDepositLog`.
2. **Fetch secrets** — pulls `WORLD_ID_APP_ID` and `DON_DECRYPTION_KEY` from the DON secret store.
3. **Confidential node mode** (`crehttp.SendRequest` + `ConsensusAggregationFromTags`) — each node independently:
   - POSTs to `https://developer.world.org/api/v4/verify/<appID>` with the nullifier hash, ZK proof, signal, and action.
   - AES-256-GCM decrypts the 32-byte encrypted blob using `DON_DECRYPTION_KEY`.
   - Extracts `[20-byte recipient address][8-byte big-endian destChainSelector]` from the plaintext.
4. **Generate DON report** — `runtime.GenerateReport` signs the ABI-encoded `(address, uint256)`.
5. **Route payout:**
   - `destChainSelector == 0 || localChain` → `evm.WriteReport` → `AuraVault.onReport`
   - otherwise → `ccip.SendMessage` with USDC token amount + encoded recipient in `data`

**Key types:**

```go
type payoutPayload struct {
    EncodedReport     []byte `consensus_aggregation:"identical"`
    DestChainSelector uint64 `consensus_aggregation:"identical"`
}

type routingResult struct {
    Recipient         [20]byte
    DestChainSelector uint64
}
```

---

## Secrets

Declared in `secrets.yaml` and stored in the Chainlink DON secret store:

| Secret | Description |
|---|---|
| `WORLD_ID_APP_ID` | World ID app ID used as the path parameter for `/verify/<appID>` |
| `DON_DECRYPTION_KEY` | 32-byte AES-256 key (hex), used to decrypt the encrypted recipient blob inside the TEE |

---

## Encrypted Recipient Format

The `encryptedRecipient` field passed to `deposit()` is a 32-byte AES-256-GCM blob:

```
[ 12-byte nonce ][ ciphertext + GCM tag ]
                          │
                   plaintext (28 bytes):
                   [ 20-byte recipient address ][ 8-byte big-endian destChainSelector ]
```

- `destChainSelector == 0` or the local Sepolia selector → same-chain payout via `onReport`.
- Any other selector → cross-chain payout via CCIP to the destination AuraVault.

> **Production note:** Replace the AES-GCM scheme with ECIES using the DON's public key so the depositor can encrypt without knowing the decryption key.

---

## Deployment

### 1. Deploy the contract

```bash
forge create contracts/evm/src/AuraVault.sol:AuraVault \
  --constructor-args <USDC_ADDRESS> <CRE_FORWARDER_ADDRESS> <CCIP_ROUTER_ADDRESS> \
  --rpc-url $SEPOLIA_RPC_URL \
  --private-key $DEPLOYER_KEY
```

**Deployed on Ethereum Sepolia:** `0x029e9b43A5dA7740440982432AAbD2820fD6Efe9`

### 2. Build the WASM handler

```bash
cd workflow
GOARCH=wasm GOOS=wasip1 go build -o vault.wasm ./vault
```

### 3. Deploy to CRE

```bash
cre workflow deploy --target staging-settings
```

### 4. Register secrets

```bash
cre secrets set WORLD_ID_APP_ID=<your_app_id>
cre secrets set DON_DECRYPTION_KEY=<hex_key>
```

---

## Dependencies

| Package | Version | Role |
|---|---|---|
| `cre-sdk-go` | v1.2.0 | Core workflow runtime, secrets, consensus |
| `cre-sdk-go/capabilities/blockchain/evm` | v1.0.0-beta.5 | Log trigger, `WriteReport` |
| `cre-sdk-go/capabilities/blockchain/ccip` | — | Cross-chain message sending |
| `cre-sdk-go/capabilities/networking/http` | v1.3.0 | Confidential HTTP (World ID API) |
