//go:build wasip1

package main

import (
	"encoding/hex"
	"fmt"
	"strings"
)

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
