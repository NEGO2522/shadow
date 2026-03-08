//go:build wasip1

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	confidentialhttp "github.com/smartcontractkit/cre-sdk-go/capabilities/networking/confidentialhttp"
	"github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/evm"
	"github.com/smartcontractkit/cre-sdk-go/cre"
)

// onShieldedOrder handles a ShieldedOrder event from the Shadow contract.
//
// Flow:
//  1. Parse the ShieldedOrder log.
//  2. RunInNodeMode: each DON node decrypts the order; all must agree on the
//     plaintext (BFT consensus before any external call is made).
//  3. POST /order to the matching service via Confidential HTTP.
//     Secrets are injected inside the CRE enclave — never visible in workflow
//     code or node memory. The orderId is idempotent on the server side.
//  4. POST /match to trigger a matching cycle and collect fills.
//  5. Log fills; the matching service handles on-chain settlement via Flashbots.
//
// EigenCloud replaced: CRE DON provides TEE execution + BFT consensus.
func onShieldedOrder(config *Config, runtime cre.Runtime, log *evm.Log) (struct{}, error) {
	logger := runtime.Logger()

	// 1. Parse log.
	orderLog, err := parseShieldedOrderLog(log)
	if err != nil {
		return struct{}{}, fmt.Errorf("parse ShieldedOrder log: %w", err)
	}

	orderUUID, err := uuid.FromBytes(orderLog.OrderID[:16])
	if err != nil {
		return struct{}{}, fmt.Errorf("parse order UUID: %w", err)
	}
	orderIDStr := orderUUID.String()

	logger.Info("ShieldedOrder received",
		"trader", hex.EncodeToString(orderLog.Trader[:]),
		"orderId", orderIDStr,
	)

	// 2. RunInNodeMode: decrypt order and reach DON consensus on plaintext.
	const decryptionKeyHex = "1bd6a9a385bbe77c7bcc6b0e24b5bb3397bb39ad25bb8039c8decf5e2a0ad6af"

	consensus, err := cre.RunInNodeMode(
		config,
		runtime,
		func(cfg *Config, nodeRuntime cre.NodeRuntime) (*orderConsensus, error) {
			plainJSON, decErr := decryptOrder(orderLog.EncryptedOrder, decryptionKeyHex)
			if decErr != nil {
				return nil, fmt.Errorf("decrypt order: %w", decErr)
			}
			// Validate well-formed JSON before consensus.
			var probe map[string]any
			if jsonErr := json.Unmarshal(plainJSON, &probe); jsonErr != nil {
				return nil, fmt.Errorf("decrypted order is not valid JSON: %w", jsonErr)
			}
			return &orderConsensus{OrderJSON: plainJSON}, nil
		},
		cre.ConsensusAggregationFromTags[*orderConsensus](),
	).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("decryption consensus: %w", err)
	}

	// Inject the on-chain orderId for server-side idempotency.
	var orderMap map[string]any
	if err := json.Unmarshal(consensus.OrderJSON, &orderMap); err != nil {
		return struct{}{}, fmt.Errorf("re-parse order JSON: %w", err)
	}
	orderMap["id"] = orderIDStr
	submitBody, err := json.Marshal(orderMap)
	if err != nil {
		return struct{}{}, fmt.Errorf("marshal order body: %w", err)
	}

	// 3. POST /order via Confidential HTTP.
	//    MATCHER_API_KEY is injected inside the CRE enclave via VaultDonSecrets.
	//    Set it with: cre secrets set MATCHER_API_KEY=<token>
	//    and add to secrets.yaml:  secretsNames: { MATCHER_API_KEY: [MATCHER_API_KEY_ALL] }
	client := confidentialhttp.Client{}

	submitResp, err := client.SendRequest(runtime, &confidentialhttp.ConfidentialHTTPRequest{
		Request: &confidentialhttp.HTTPRequest{
			Url:    config.MatcherURL + "/order",
			Method: "POST",
			Body:   &confidentialhttp.HTTPRequest_BodyBytes{BodyBytes: submitBody},
			MultiHeaders: map[string]*confidentialhttp.HeaderValues{
				"Content-Type":  {Values: []string{"application/json"}},
				"Authorization": {Values: []string{"Bearer {{.MATCHER_API_KEY}}"}},
			},
		},
		VaultDonSecrets: []*confidentialhttp.SecretIdentifier{
			{Key: "MATCHER_API_KEY"},
		},
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("POST /order: %w", err)
	}
	if submitResp.StatusCode != 200 {
		return struct{}{}, fmt.Errorf("POST /order %d: %s", submitResp.StatusCode, submitResp.Body)
	}
	logger.Info("Order submitted", "orderId", orderIDStr)

	// 4. POST /match — trigger a matching cycle.
	matchResp, err := client.SendRequest(runtime, &confidentialhttp.ConfidentialHTTPRequest{
		Request: &confidentialhttp.HTTPRequest{
			Url:    config.MatcherURL + "/match",
			Method: "POST",
			MultiHeaders: map[string]*confidentialhttp.HeaderValues{
				"Authorization": {Values: []string{"Bearer {{.MATCHER_API_KEY}}"}},
			},
		},
		VaultDonSecrets: []*confidentialhttp.SecretIdentifier{
			{Key: "MATCHER_API_KEY"},
		},
	}).Await()
	if err != nil {
		return struct{}{}, fmt.Errorf("POST /match: %w", err)
	}
	if matchResp.StatusCode != 200 {
		return struct{}{}, fmt.Errorf("POST /match %d: %s", matchResp.StatusCode, matchResp.Body)
	}

	// 5. Log fills.
	var result matcherMatchResp
	if err := json.Unmarshal(matchResp.Body, &result); err != nil {
		return struct{}{}, fmt.Errorf("parse /match response: %w", err)
	}
	if result.Fills == 0 {
		logger.Info("No fills this cycle", "orderId", orderIDStr)
		return struct{}{}, nil
	}
	for _, fill := range result.Results {
		logger.Info("Fill",
			"buyOrderId", fill.BuyOrderID,
			"sellOrderId", fill.SellOrderID,
			"price", fill.Price,
			"quantity", fill.Quantity,
			"txHash", fill.TxHash,
			"screener", fill.Screener,
		)
	}
	return struct{}{}, nil
}

// decryptOrder decrypts a variable-length AES-256-CTR blob to plaintext JSON.
// Layout: [12-byte nonce][ciphertext...] — same scheme as decryptRecipient.
func decryptOrder(ciphertext []byte, keyHex string) ([]byte, error) {
	keyBytes, err := hex.DecodeString(strings.TrimPrefix(keyHex, "0x"))
	if err != nil || len(keyBytes) != 32 {
		return nil, fmt.Errorf("DON_DECRYPTION_KEY must be 32 hex bytes")
	}
	if len(ciphertext) < 13 {
		return nil, fmt.Errorf("encryptedOrder too short: %d bytes", len(ciphertext))
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	var iv [16]byte
	copy(iv[:12], ciphertext[:12]) // 12-byte nonce zero-padded to 16-byte CTR IV

	plaintext := make([]byte, len(ciphertext)-12)
	cipher.NewCTR(block, iv[:]).XORKeyStream(plaintext, ciphertext[12:])
	return plaintext, nil
}
