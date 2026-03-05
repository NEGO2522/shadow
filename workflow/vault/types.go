//go:build wasip1

package main

import "math/big"

// Config loaded from config.staging.json / config.production.json.
type Config struct {
	ContractAddress string `json:"contract_address"`
}

// payoutPayload is the consensus result agreed upon by all DON nodes.
type payoutPayload struct {
	EncodedReport []byte `consensus_aggregation:"identical"`
}

// ShieldedDepositLog holds the decoded fields from a ShieldedDeposit event.
type ShieldedDepositLog struct {
	Sender             [20]byte
	EncryptedRecipient [32]byte
	Amount             *big.Int
}
