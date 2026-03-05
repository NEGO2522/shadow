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

// keccak256("ShieldedDeposit(address,bytes32,uint256)")
const shieldedDepositTopic = "0xa528a5618e16311d917a25d1b6f7d83f001f50e5b4bee369286d16c83784e22a"

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	contractAddrBytes, err := hexToBytes20(config.ContractAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid contract_address: %w", err)
	}

	eventSig, err := hex.DecodeString(strings.TrimPrefix(shieldedDepositTopic, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid event topic: %w", err)
	}

	trigger := evm.LogTrigger(evm.EthereumTestnetSepolia, &evm.FilterLogTriggerRequest{
		Addresses: [][]byte{contractAddrBytes[:]},
		Topics: []*evm.TopicValues{
			{Values: [][]byte{eventSig}},
			{Values: [][]byte{}}, // any sender
		},
	})

	return cre.Workflow[*Config]{
		cre.Handler(trigger, onShieldedDeposit),
	}, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
