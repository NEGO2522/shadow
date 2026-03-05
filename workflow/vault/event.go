//go:build wasip1

package main

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
)

// parseShieldedDepositLog decodes a ShieldedDeposit evm.Log.
// Topics[1]: indexed sender (32 bytes, right-aligned)
// Data:      abi.encode(bytes32 encryptedRecipient, uint256 amount) = 64 bytes
func parseShieldedDepositLog(log *evm.Log) (*ShieldedDepositLog, error) {
	if len(log.Topics) < 2 {
		return nil, fmt.Errorf("expected ≥2 topics, got %d", len(log.Topics))
	}
	if len(log.Data) < 64 {
		return nil, fmt.Errorf("expected ≥64 bytes data, got %d", len(log.Data))
	}

	var d ShieldedDepositLog
	senderTopic := log.Topics[1]
	copy(d.Sender[:], senderTopic[len(senderTopic)-20:])
	copy(d.EncryptedRecipient[:], log.Data[0:32])
	d.Amount = new(big.Int).SetBytes(log.Data[32:64])
	return &d, nil
}
