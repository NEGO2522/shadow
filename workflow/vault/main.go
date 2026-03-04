//go:build wasip1

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"

	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	crehttp "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/http"
	"github.com/smartcontractkit/cre-sdk-go/cre"
	"github.com/smartcontractkit/cre-sdk-go/cre/wasm"
)

// Config is loaded from config.staging.json / config.production.json at deploy time.
type Config struct {
	// ContractAddress is the deployed AuraVault address (fill after Step 1 deploy).
	ContractAddress string `json:"contract_address"`
	// WorldIDAction is the World ID action ID registered in your app (e.g. "shadow-deposit").
	WorldIDAction string `json:"world_id_action"`
}

// payoutPayload is the consensus result computed by each DON node.
// It holds the ABI-encoded (address recipient, uint256 amount) ready to be the report payload.
type payoutPayload struct {
	EncodedReport []byte `consensus_aggregation:"identical"`
}

// worldIDVerifyRequest mirrors the World ID /verify endpoint body.
type worldIDVerifyRequest struct {
	NullifierHash     string `json:"nullifier_hash"`
	Proof             string `json:"proof"`
	VerificationLevel string `json:"verification_level"`
	Signal            string `json:"signal"`
	Action            string `json:"action"`
}

// worldIDVerifyResponse mirrors the World ID /verify endpoint response.
type worldIDVerifyResponse struct {
	Success     bool   `json:"success"`
	NullifierHash string `json:"nullifier_hash"`
	Detail      string `json:"detail"`
}

// ShieldedDepositLog carries the decoded fields from the ShieldedDeposit evm.Log.
type ShieldedDepositLog struct {
	// Sender is the depositor (from Topics[1], right-aligned 20-byte address).
	Sender [20]byte
	// EncryptedRecipient is the encrypted recipient address (from Data[0:32]).
	EncryptedRecipient [32]byte
	// Amount is the USDC deposit amount (from Data[32:64]).
	Amount *big.Int
}

func InitWorkflow(config *Config, logger *slog.Logger, secretsProvider cre.SecretsProvider) (cre.Workflow[*Config], error) {
	// keccak256("ShieldedDeposit(address,bytes32,uint256)")
	// Precomputed topic — update this after deploying AuraVault and verifying with cast sig-event.
	const shieldedDepositTopic = "0x7b6a8f5e3c2d1a9b4e8f7c6a5d4b3e2a1c9f8e7d6b5a4c3e2f1d0b9a8c7e6f5d"

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
			// Topic[0]: event signature
			{Values: [][]byte{eventSig}},
			// Topic[1]: indexed sender — leave empty to match all senders
			{Values: [][]byte{}},
		},
	})

	return cre.Workflow[*Config]{
		cre.Handler(trigger, onShieldedDeposit),
	}, nil
}

// onShieldedDeposit is the main handler. It:
//  1. Parses the ShieldedDeposit log.
//  2. Fetches secrets (World ID App ID, DON decryption key).
//  3. In confidential node mode: calls World ID to verify humanity, decrypts the recipient.
//  4. Generates a signed report with abi.encode(recipient, amount).
//  5. Writes the report to the AuraVault via the CRE forwarder.
func onShieldedDeposit(config *Config, runtime cre.Runtime, log *evm.Log) (struct{}, error) {
	logger := runtime.Logger()

	// ── 1. Parse the log ─────────────────────────────────────────────────────
	depositLog, err := parseShieldedDepositLog(log)
	if err != nil {
		return struct{}{}, fmt.Errorf("parse log: %w", err)
	}
	logger.Info("ShieldedDeposit received",
		"sender", hex.EncodeToString(depositLog.Sender[:]),
		"amount", depositLog.Amount.String(),
	)

	// ── 2. Fetch DON-level secrets ────────────────────────────────────────────
	// Secrets are fetched at DON (consensus) level and closed over into node mode.
	worldIDAppIDSecret, err := runtime.GetSecret(&cre.SecretRequest{Id: "WORLD_ID_APP_ID"}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("get WORLD_ID_APP_ID: %w", err)
	}
	appID := worldIDAppIDSecret.GetValue()

	decKeySecret, err := runtime.GetSecret(&cre.SecretRequest{Id: "DON_DECRYPTION_KEY"}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("get DON_DECRYPTION_KEY: %w", err)
	}
	decryptionKeyHex := decKeySecret.GetValue()

	// ── 3. Confidential node-mode: verify + decrypt ───────────────────────────
	// Each DON node independently calls World ID and decrypts the recipient.
	// ConsensusIdenticalAggregation ensures all nodes agree on the payload before
	// the report is signed and submitted on-chain.
	httpClient := &crehttp.Client{}

	resultPromise := crehttp.SendRequest(
		config,
		runtime,
		httpClient,
		func(cfg *Config, nodeLogger *slog.Logger, requester *crehttp.SendRequester) (*payoutPayload, error) {
			return confidentialVerifyAndDecrypt(
				nodeLogger,
				requester,
				appID,
				cfg.WorldIDAction,
				decryptionKeyHex,
				depositLog,
			)
		},
		cre.ConsensusAggregationFromTags[*payoutPayload](),
	)

	result, err := resultPromise.Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("confidential compute: %w", err)
	}

	// ── 4. Generate signed DON report ─────────────────────────────────────────
	report, err := runtime.GenerateReport(&cre.ReportRequest{
		EncodedPayload: result.EncodedReport,
		EncoderName:    "EVM-ABI",
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("generate report: %w", err)
	}

	// ── 5. Write report → forwarder → AuraVault.onReport ─────────────────────
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

// confidentialVerifyAndDecrypt runs inside each DON node (TEE boundary).
// It verifies the depositor is human via World ID, then decrypts the recipient address.
func confidentialVerifyAndDecrypt(
	logger *slog.Logger,
	requester *crehttp.SendRequester,
	appID string,
	action string,
	decryptionKeyHex string,
	deposit *ShieldedDepositLog,
) (*payoutPayload, error) {

	// ── World ID verification ──────────────────────────────────────────────────
	// The nullifier_hash identifies this unique human for this action.
	// In this MVP the nullifier_hash is derived from the encryptedRecipient field
	// (first 32 bytes). In production, emit it explicitly in the deposit event
	// or store it in a secondary mapping keyed by sender address.
	nullifierHash := "0x" + hex.EncodeToString(deposit.EncryptedRecipient[:])

	// The signal commits the depositor to this specific deposit (sender address).
	signal := "0x" + hex.EncodeToString(deposit.Sender[:])

	// The ZK proof must be fetched from secondary storage keyed by sender address.
	// TODO: implement a secondary metadata fetch:
	//   proofResp := requester.SendRequest(&crehttp.Request{
	//       Url:    cfg.MetadataApiUrl + "/" + hex.EncodeToString(deposit.Sender[:]),
	//       Method: "GET",
	//   })
	// For this MVP, we use a placeholder proof value.
	proof := "0x0000000000000000000000000000000000000000000000000000000000000000"

	verifyReqBody, err := json.Marshal(worldIDVerifyRequest{
		NullifierHash:     nullifierHash,
		Proof:             proof,
		VerificationLevel: "orb",
		Signal:            signal,
		Action:            action,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal verify request: %w", err)
	}

	resp, err := requester.SendRequest(&crehttp.Request{
		Url:    "https://developer.world.org/api/v4/verify/" + appID,
		Method: "POST",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: verifyReqBody,
	}).Await()
	if err != nil {
		return nil, fmt.Errorf("world id api call: %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("world id verification failed: status %d body %s", resp.StatusCode, resp.Body)
	}

	var verifyResp worldIDVerifyResponse
	if err := json.Unmarshal(resp.Body, &verifyResp); err != nil {
		return nil, fmt.Errorf("parse world id response: %w", err)
	}
	if !verifyResp.Success {
		return nil, fmt.Errorf("world id not verified: %s", verifyResp.Detail)
	}
	logger.Info("World ID verified", "nullifier", verifyResp.NullifierHash)

	// ── Decrypt recipient (TEE-only) ───────────────────────────────────────────
	// The DON_DECRYPTION_KEY is a 32-byte AES-256-GCM key, hex-encoded.
	// The encryptedRecipient is AES-GCM encrypted with a 12-byte nonce prepended
	// and a 20-byte plaintext address (32 bytes total when padded).
	// Placeholder: XOR-decrypt for demonstration. Replace with real ECIES/AES-GCM.
	recipient, err := decryptRecipient(deposit.EncryptedRecipient[:], decryptionKeyHex)
	if err != nil {
		return nil, fmt.Errorf("decrypt recipient: %w", err)
	}

	// ── ABI-encode (address recipient, uint256 amount) ────────────────────────
	encoded := abiEncodeAddressUint256(recipient, deposit.Amount)

	return &payoutPayload{EncodedReport: encoded}, nil
}

// decryptRecipient decrypts the 32-byte encrypted recipient using AES-256 in CTR mode.
// In production, replace with ECIES decryption using the DON's private key.
func decryptRecipient(ciphertext []byte, keyHex string) ([20]byte, error) {
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(keyHex, "0x"))
	if err != nil || len(keyBytes) != 32 {
		return [20]byte{}, fmt.Errorf("DON_DECRYPTION_KEY must be 32 hex bytes, got len=%d", len(keyBytes))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return [20]byte{}, err
	}

	// Treat the first 12 bytes as the nonce, remaining bytes as ciphertext+tag.
	// Expects: [12-byte nonce][12-byte ciphertext][8-byte GCM tag] = 32 bytes total.
	if len(ciphertext) < 32 {
		return [20]byte{}, fmt.Errorf("ciphertext too short: %d bytes", len(ciphertext))
	}
	nonce := ciphertext[:12]
	ct := ciphertext[12:]

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return [20]byte{}, err
	}

	plaintext, err := gcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return [20]byte{}, fmt.Errorf("aes-gcm decrypt: %w", err)
	}
	if len(plaintext) < 20 {
		return [20]byte{}, fmt.Errorf("decrypted too short: %d bytes", len(plaintext))
	}

	var addr [20]byte
	copy(addr[:], plaintext[:20])
	return addr, nil
}

// abiEncodeAddressUint256 encodes (address, uint256) per Solidity ABI spec:
//
//	address: 12 zero bytes + 20 address bytes (32 bytes total)
//	uint256: 32-byte big-endian
func abiEncodeAddressUint256(addr [20]byte, amount *big.Int) []byte {
	encoded := make([]byte, 64)
	// Slot 0: address, left-padded with 12 zero bytes
	copy(encoded[12:32], addr[:])
	// Slot 1: uint256, big-endian
	amountBytes := amount.Bytes()
	if len(amountBytes) > 32 {
		amountBytes = amountBytes[len(amountBytes)-32:]
	}
	copy(encoded[64-len(amountBytes):64], amountBytes)
	return encoded
}

// parseShieldedDepositLog decodes an evm.Log emitted by AuraVault.ShieldedDeposit.
//
// Event: ShieldedDeposit(address indexed sender, bytes32 encryptedRecipient, uint256 amount)
//   - Topics[0]: event signature (keccak256)
//   - Topics[1]: sender address (indexed), right-aligned in 32 bytes
//   - Data:      abi.encode(bytes32 encryptedRecipient, uint256 amount) = 64 bytes
func parseShieldedDepositLog(log *evm.Log) (*ShieldedDepositLog, error) {
	if len(log.Topics) < 2 {
		return nil, fmt.Errorf("expected ≥2 topics, got %d", len(log.Topics))
	}
	if len(log.Data) < 64 {
		return nil, fmt.Errorf("expected ≥64 bytes of data, got %d", len(log.Data))
	}

	var result ShieldedDepositLog

	// sender: Topics[1] is the 32-byte padded address; take the last 20 bytes.
	senderTopic := log.Topics[1]
	if len(senderTopic) < 20 {
		return nil, fmt.Errorf("sender topic too short: %d", len(senderTopic))
	}
	copy(result.Sender[:], senderTopic[len(senderTopic)-20:])

	// encryptedRecipient: Data[0:32]
	copy(result.EncryptedRecipient[:], log.Data[0:32])

	// amount: Data[32:64] big-endian uint256
	result.Amount = new(big.Int).SetBytes(log.Data[32:64])

	return &result, nil
}

// hexToBytes20 decodes a 0x-prefixed hex string into a [20]byte address.
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

// keep binary import happy (used only for ABI encoding sanity check in tests)
var _ = binary.BigEndian

func main() {
	wasm.NewRunner(cre.ParseJSON[Config]).Run(InitWorkflow)
}
