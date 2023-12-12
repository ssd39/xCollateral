package signer

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

func HashEvmMessage(payload string) ([]byte, error) {
	message, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding payload %+v", err)
		return []byte{}, err
	}
	messagePrefix := []byte("\x19Ethereum Signed Message:\n")
	messageLength := []byte(fmt.Sprintf("%v", len(message)))

	var hashData []byte
	hashData = append(hashData, messagePrefix...)
	hashData = append(hashData, messageLength...)
	hashData = append(hashData, message...)

	return crypto.Keccak256(hashData), nil
}

func NormalizeSignature(signature []byte) []byte {
	if signature[len(signature)-1] == 0 {
		signature[len(signature)-1] = 0x1b
	} else if signature[len(signature)-1] == 1 {
		signature[len(signature)-1] = 0x1c
	}
	return signature
}
