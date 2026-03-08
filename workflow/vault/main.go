//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

const shieldedDepositTopic = "0xa528a5618e16311d917a25d1b6f7d83f001f50e5b4bee369286d16c83784e22a"

// shieldedOrderTopic = keccak256("ShieldedOrder(address,bytes,bytes32)")
// Compute with: cast sig-event "ShieldedOrder(address,bytes,bytes32)"
// Replace this placeholder once the contract is deployed.
const shieldedOrderTopic = "0xf6604c3bcd7473aa76b7d82d56a4453e9b64a148ca27eb7c07e9f698093b5af5"

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	contractAddrBytes, err := hexToBytes20(config.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid contract_address: %w", err)
	}

	depositSig, err := hex.DecodeString(strings.TrimPrefix(shieldedDepositTopic, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid deposit event topic: %w", err)
	}

	orderSig, err := hex.DecodeString(strings.TrimPrefix(shieldedOrderTopic, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid order event topic: %w", err)
	}

	// Trigger 1: ShieldedDeposit — decrypt recipient and execute USDC payout.
	depositTrigger := evm.LogTrigger(evm.EthereumTestnetSepolia, &evm.FilterLogTriggerRequest{
		Addresses: [][]byte{contractAddrBytes[:]},
		Topics: []*evm.TopicValues{
			{Values: [][]byte{depositSig}},
			{Values: [][]byte{}}, // any sender
		},
	})

	// Trigger 2: ShieldedOrder — decrypt order and relay to dark-pool matcher.
	orderTrigger := evm.LogTrigger(evm.EthereumTestnetSepolia, &evm.FilterLogTriggerRequest{
		Addresses: [][]byte{contractAddrBytes[:]},
		Topics: []*evm.TopicValues{
			{Values: [][]byte{orderSig}},
			{Values: [][]byte{}}, // any trader
		},
	})

	return cre.Workflow[*Config]{
		cre.Handler(depositTrigger, onShieldedDeposit),
		cre.Handler(orderTrigger, onShieldedOrder),
	}, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
