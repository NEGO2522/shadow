//go:build wasip1

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

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
