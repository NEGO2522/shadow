//go:build wasip1

package main

import "math/big"

// Config loaded from config.staging.json / config.production.json.
type Config struct {
	ContractAddress string `json:"contract_address"`
	MatcherURL      string `json:"matcher_url"` // base URL of the dark-pool matching service
}

// payoutPayload is the consensus result agreed upon by all DON nodes (ShieldedDeposit flow).
type payoutPayload struct {
	EncodedReport []byte `consensus_aggregation:"identical"`
}

// ShieldedDepositLog holds the decoded fields from a ShieldedDeposit event.
type ShieldedDepositLog struct {
	Sender             [20]byte
	EncryptedRecipient [32]byte
	Amount             *big.Int
}

// ── Dark-pool order matching ──────────────────────────────────────────────────

// ShieldedOrderLog holds fields decoded from a ShieldedOrder event.
//
// Solidity signature:
//
//	event ShieldedOrder(address indexed trader, bytes encryptedOrder, bytes32 indexed orderId);
//
// Compute the topic with: cast sig-event "ShieldedOrder(address,bytes,bytes32)"
type ShieldedOrderLog struct {
	Trader         [20]byte
	OrderID        [32]byte // bytes32; first 16 bytes are the UUID
	EncryptedOrder []byte   // AES-256-CTR: [12-byte nonce][ciphertext of JSON order]
}

// orderConsensus is the DON consensus result for a decrypted order (RunInNodeMode).
type orderConsensus struct {
	OrderJSON []byte `consensus_aggregation:"identical"`
}

// matchFill is one fill record returned by POST /match.
type matchFill struct {
	BuyOrderID  string `json:"buy_order_id"`
	SellOrderID string `json:"sell_order_id"`
	Price       uint64 `json:"price"`
	Quantity    uint64 `json:"quantity"`
	TxHash      string `json:"tx_hash"` // empty if relay not configured
	Screener    string `json:"screener"`
}

// matcherMatchResp is the full response from POST /match.
type matcherMatchResp struct {
	Fills   int         `json:"fills"`
	Results []matchFill `json:"results"`
}
