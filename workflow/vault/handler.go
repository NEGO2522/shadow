//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"

	"workflow/contracts/evm/src/generated/shadow"
)

func onShieldedDeposit(config *Config, runtime cre.Runtime, log *evm.Log) (struct{}, error) {
	logger := runtime.Logger()

	// 1. Parse the ShieldedDeposit log
	depositLog, err := parseShieldedDepositLog(log)
	if err != nil {
		return struct{}{}, fmt.Errorf("parse log: %w", err)
	}
	logger.Info("ShieldedDeposit received",
		"sender", hex.EncodeToString(depositLog.Sender[:]),
		"amount", depositLog.Amount.String(),
	)

	const decryptionKeyHex = "1bd6a9a385bbe77c7bcc6b0e24b5bb3397bb39ad25bb8039c8decf5e2a0ad6af"

	// 2. Node mode: decrypt recipient + reach DON consensus on the payout payload.
	// World ID is verified at the frontend API layer before deposit is accepted.
	result, err := cre.RunInNodeMode(
		config,
		runtime,
		func(cfg *Config, nodeRuntime cre.NodeRuntime) (*payoutPayload, error) {
			recipient, decErr := decryptRecipient(depositLog.EncryptedRecipient[:], decryptionKeyHex)
			if decErr != nil {
				return nil, fmt.Errorf("decrypt: %w", decErr)
			}
			return &payoutPayload{
				EncodedReport: abiEncodeAddressUint256(recipient, depositLog.Amount),
			}, nil
		},
		cre.ConsensusAggregationFromTags[*payoutPayload](),
	).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("node mode: %w", err)
	}

	// 3. Generate signed DON report
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: result.EncodedReport,
		EncoderName:    "evm",
		SigningAlgo:    "ecdsa",
		HashingAlgo:    "keccak256",
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("generate report: %w", err)
	}

	// 4. Write report → Shadow.onReport → _processReport → _executePayout via generated binding
	evmClient := &evm.Client{ChainSelector: evm.EthereumTestnetSepolia}
	contractAddress := common.HexToAddress(config.ContractAddress)

	shadowContract, err := shadow.NewShadow(evmClient, contractAddress, nil)
	if err != nil {
		return struct{}{}, fmt.Errorf("create shadow binding: %w", err)
	}

	resp, err := shadowContract.WriteReport(runtime, report, nil).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("write report: %w", err)
	}

	logger.Info("ShieldedPayout submitted",
		"amount", depositLog.Amount.String(),
		"txHash", common.BytesToHash(resp.TxHash).Hex(),
	)
	return struct{}{}, nil
}
