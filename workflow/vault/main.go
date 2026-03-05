//go:build wasip1

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

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

	// 2. Fetch decryption key from DON secret store
	decKeySecret, err := runtime.GetSecret(&cre.SecretRequest{Id: "DON_DECRYPTION_KEY"}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("get DON_DECRYPTION_KEY: %w", err)
	}

	// 3. Node mode: decrypt recipient + reach DON consensus on the payout payload
	// TODO: add World ID HTTP verification here before production
	result, err := cre.RunInNodeMode(
		config,
		runtime,
		func(cfg *Config, nodeRuntime cre.NodeRuntime) (*payoutPayload, error) {
			recipient, err := decryptRecipient(depositLog.EncryptedRecipient[:], decKeySecret.GetValue())
			if err != nil {
				return nil, fmt.Errorf("decrypt: %w", err)
			}
			nodeRuntime.Logger().Info("World ID check: mock pass (wire up before production)")
			return &payoutPayload{
				EncodedReport: abiEncodeAddressUint256(recipient, depositLog.Amount),
			}, nil
		},
		cre.ConsensusAggregationFromTags[*payoutPayload](),
	).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("node mode: %w", err)
	}

	_ = result // used below — sentinel to suppress duplicate block
	result2, err := cre.RunInNodeMode(
		config,
		runtime,
		func(cfg *Config, nodeRuntime cre.NodeRuntime) (*payoutPayload, error) {
			recipient, err := decryptRecipient(depositLog.EncryptedRecipient[:], decKeySecret.GetValue())
			if err != nil {
				return nil, fmt.Errorf("decrypt: %w", err)
			}
			nodeRuntime.Logger().Info("World ID check: mock pass (wire up before production)")
			return &payoutPayload{
				EncodedReport: abiEncodeAddressUint256(recipient, depositLog.Amount),
			}, nil
		},
		cre.ConsensusAggregationFromTags[*payoutPayload](),
	)
	if err != nil {
		return struct{}{}, fmt.Errorf("node mode: %w", err)
	}

	// 4. Generate signed DON report
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: result.EncodedReport,
		EncoderName:    "EVM-ABI",
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("generate report: %w", err)
	}

	// 5. Write report → CRE forwarder → Shadow._processReport → payout
	contractAddrBytes, _ := hexToBytes20(config.ContractAddress)
	evmClient := &evm.Client{ChainSelector: evm.EthereumTestnetSepolia}
	_, err = evmClient.WriteReport(runtime, &evm.WriteCreReportRequest{
		Receiver: contractAddrBytes[:],
		Report:   report,
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("write report: %w", err)
	}

	logger.Info("ShieldedPayout submitted", "amount", depositLog.Amount.String())
	return struct{}{}, nil
}

// decryptRecipient decrypts the 32-byte AES-256-CTR blob.
// Layout: [12-byte nonce][20-byte ciphertext] — fits Solidity bytes32.
func decryptRecipient(ciphertext []byte, keyHex string) ([20]byte, error) {
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(keyHex, "0x"))
	if err != nil || len(keyBytes) != 32 {
		return [20]byte{}, fmt.Errorf("DON_DECRYPTION_KEY must be 32 hex bytes")
	}
	if len(ciphertext) < 32 {
		return [20]byte{}, fmt.Errorf("ciphertext too short: %d bytes", len(ciphertext))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return [20]byte{}, err
	}

	var iv [16]byte
	copy(iv[:12], ciphertext[:12]) // 12-byte nonce zero-padded to 16-byte CTR IV

	plaintext := make([]byte, 20)
	cipher.NewCTR(block, iv[:]).XORKeyStream(plaintext, ciphertext[12:32])

	var addr [20]byte
	copy(addr[:], plaintext)
	return addr, nil
}

// abiEncodeAddressUint256 encodes (address, uint256) per Solidity ABI spec.
func abiEncodeAddressUint256(addr [20]byte, amount *big.Int) []byte {
	encoded := make([]byte, 64)
	copy(encoded[12:32], addr[:])
	b := amount.Bytes()
	copy(encoded[64-len(b):64], b)
	return encoded
}

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

func hexToBytes20(hexAddr string) ([20]byte, error) {
	b, err := hex.DecodeString(strings.TrimPrefix(hexAddr, "0x"))
	if err != nil {
		return [20]byte{}, err
	}
	if len(b) != 20 {
		return [20]byte{}, fmt.Errorf("expected 20 bytes, got %d", len(b))
	}
	var addr [20]byte
	copy(addr[:], b)
	return addr, nil
}

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
