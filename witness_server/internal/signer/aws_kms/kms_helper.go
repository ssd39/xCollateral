package aws_kms

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ripemd160"
)

type EcdsaPublicKeyInfo struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.ObjectIdentifier
}

type EcdsaPublicKey struct {
	EcPublicKeyInfo EcdsaPublicKeyInfo
	PublicKey       asn1.BitString
}

func HashSha256(bytes []byte) []byte {
	hashSha256 := sha256.New()
	hashSha256.Write(bytes)
	return hashSha256.Sum(nil)
}

func HashSha512(bytes []byte) []byte {
	hashSha512 := sha512.New()
	hashSha512.Write(bytes)
	return hashSha512.Sum(nil)
}

func HashRipemd160(bytes []byte) []byte {
	hashRipemd160 := ripemd160.New()
	hashRipemd160.Write(bytes)
	return hashRipemd160.Sum(nil)
}

func GetEthereumSignature(expectedPublicKeyBytes []byte, txHash []byte, r []byte, s []byte) []byte {
	rsSignature := append(adjustSignatureLength(r), adjustSignatureLength(s)...)
	signature := append(rsSignature, []byte{0}...)

	recoveredPublicKeyBytes, err := crypto.Ecrecover(txHash, signature)
	if err != nil {
		log.Error().Msgf("Error recovering public key from signature: '%s'", err)
		return nil
	}

	if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
		signature = append(rsSignature, []byte{1}...)
		recoveredPublicKeyBytes, err = crypto.Ecrecover(txHash, signature)
		if err != nil {
			log.Error().Msgf("Error recovering public key from signature: '%s'", err)
			return nil
		}

		if hex.EncodeToString(recoveredPublicKeyBytes) != hex.EncodeToString(expectedPublicKeyBytes) {
			log.Error().Msg("Error reconstructing public key from signature")
			return nil
		}
	}

	return signature
}

func adjustSignatureLength(buffer []byte) []byte {
	buffer = bytes.TrimLeft(buffer, "\x00")
	for len(buffer) < 32 {
		zeroBuf := []byte{0}
		buffer = append(zeroBuf, buffer...)
	}
	return buffer
}

func BigIntBytesToSignatureHex(rBytes, sBytes []byte) string {
	var signature struct{ R, S *big.Int }
	signature.R = new(big.Int).SetBytes(rBytes)
	signature.S = new(big.Int).SetBytes(sBytes)

	sig, err := asn1.Marshal(signature)
	if err != nil {
		log.Error().Msgf("Error with marshal r and s signature: '%+v'", err)
		return ""
	}

	return hex.EncodeToString(sig)
}

func EvmAddressToXrplAccount(evmAddress string) string {
	evmBytes, err := hex.DecodeString(evmAddress[2:])
	if err != nil {
		log.Error().Msgf("Error EvmAddressToXrplAccount for address %+v: '%+v'", evmAddress, err)
		return ""
	}

	return EncodeAccountId(evmBytes)
}

func XrplAccountToEvmAddress(xrplAccount string) string {
	decoded := decodeAccountId(xrplAccount)
	if decoded == nil {
		return ""
	}
	evmAddress := "0x" + hex.EncodeToString(decoded)

	return evmAddress
}

var rippleAlphabet *base58.Alphabet = base58.NewAlphabet("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")

func EncodeAccountId(dataBytes []byte) string {
	hashWithVersion := append([]byte{0}, dataBytes...)
	checkSum := HashSha256(HashSha256(hashWithVersion))[:4]
	toEncodeBase58 := append(hashWithVersion, checkSum...)
	xrplAccount := base58.EncodeAlphabet(toEncodeBase58, rippleAlphabet)

	return xrplAccount
}

func decodeAccountId(xrplAccount string) []byte {
	decoded, err := base58.DecodeAlphabet(xrplAccount, rippleAlphabet)
	if err != nil {
		log.Error().Msgf("Error decoding xrplAccount: '%+v'", err)
		return nil
	}
	version := decoded[0 : len(decoded)-24]
	if []byte{0}[0] != version[0] {
		log.Error().Msgf("Decoded version must be 0 instead of: '%+v'", version[0])
		return nil
	}

	return decoded[len(decoded)-24 : len(decoded)-4]
}
