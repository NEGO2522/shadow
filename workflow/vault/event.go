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

// parseShieldedOrderLog decodes a ShieldedOrder evm.Log.
//
// Solidity: event ShieldedOrder(address indexed trader, bytes encryptedOrder, bytes32 indexed orderId)
//
//   - Topics[1]: indexed trader address (32 bytes, right-aligned)
//   - Topics[2]: indexed orderId (bytes32)
//   - Data:      ABI-encoded `bytes encryptedOrder`
//     [0:32]  = offset (0x20)
//     [32:64] = byte length of encryptedOrder
//     [64:]   = encryptedOrder data (padded to 32-byte boundary)
func parseShieldedOrderLog(log *evm.Log) (*ShieldedOrderLog, error) {
	if len(log.Topics) < 3 {
		return nil, fmt.Errorf("ShieldedOrder: expected ≥3 topics, got %d", len(log.Topics))
	}
	if len(log.Data) < 64 {
		return nil, fmt.Errorf("ShieldedOrder: data too short: %d bytes", len(log.Data))
	}

	var d ShieldedOrderLog

	// trader from topics[1] (right-aligned address)
	t := log.Topics[1]
	copy(d.Trader[:], t[len(t)-20:])

	// orderId from topics[2]
	copy(d.OrderID[:], log.Topics[2])

	// encryptedOrder: length at [32:64], data at [64:64+length]
	orderLen := new(big.Int).SetBytes(log.Data[32:64]).Uint64()
	if uint64(len(log.Data)) < 64+orderLen {
		return nil, fmt.Errorf("ShieldedOrder: data truncated (need %d, have %d)", 64+orderLen, len(log.Data))
	}
	d.EncryptedOrder = make([]byte, orderLen)
	copy(d.EncryptedOrder, log.Data[64:64+orderLen])
	return &d, nil
}
