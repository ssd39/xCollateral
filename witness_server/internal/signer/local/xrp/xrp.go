package xrp

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/hex"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"peersyst/bridge-witness-go/internal/signer"
	awsKms "peersyst/bridge-witness-go/internal/signer/aws_kms"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

var _ signer.SignerProvider = &XrpLocalSignerProvider{}

type XrpLocalSignerProvider struct {
	privateKey *ecdsa.PrivateKey
}

func (p *XrpLocalSignerProvider) SignTransaction(payload string, opts interface{}) string {
	tx, err := transaction.UnmarshalTransaction(payload)
	if err != nil {
		log.Error().Msgf("Error decoding xrp transaction: '%s'", err)
		return ""
	}
	tx.SetSigningPubKey(p.GetPublicKey())

	encodeForSigning := signer.EncodeXrpTransactionForSigning(&tx)
	if encodeForSigning == "" {
		log.Error().Msgf("Error encoding xrp transaction for signing")
		return ""
	}
	encodedBytes, err := hex.DecodeString(encodeForSigning)
	if err != nil {
		log.Error().Msgf("Error decoding encoded xrp transaction: '%s'", err)
		return ""
	}

	encodedHash := awsKms.HashSha512(encodedBytes)[0:32]
	rBytes, sBytes, _ := p.sign(encodedHash)
	if rBytes == nil || sBytes == nil {
		// Error logged on kms
		return ""
	}

	hexSignature := awsKms.BigIntBytesToSignatureHex(rBytes, sBytes)
	if hexSignature == "" {
		// Error logged on kms
		return ""
	}
	tx.SetTxnSignature(hexSignature)

	return signer.EncodeXrpTransaction(&tx)
}

func (p *XrpLocalSignerProvider) SignMultiSigTransaction(payload string) string {
	tx := signer.DecodeXrpTransaction(payload)
	if tx == nil {
		log.Error().Msgf("Error decoding xrp transaction")
		return ""
	}
	tx.SetSigningPubKey("")

	encodeForMultiSigning := signer.EncodeXrpTransactionForMultiSigning(&tx, p.GetAddress())
	if encodeForMultiSigning == "" {
		log.Error().Msgf("Error encoding xrp transaction for multi signing")
		return ""
	}
	encodedBytes, err := hex.DecodeString(encodeForMultiSigning)
	if err != nil {
		log.Error().Msgf("Error decoding encoded xrp transaction: '%s'", err)
		return ""
	}

	encodedHash := awsKms.HashSha512(encodedBytes)[0:32]
	rBytes, sBytes, _ := p.sign(encodedHash)
	if rBytes == nil || sBytes == nil {
		// Error logged on kms
		return ""
	}

	hexSignature := awsKms.BigIntBytesToSignatureHex(rBytes, sBytes)
	if hexSignature == "" {
		// Error logged on kms
		return ""
	}

	return hexSignature
}

func (p *XrpLocalSignerProvider) GetAddress() string {
	// Xrp address encoding: https://xrpl.org/img/address-encoding.svg
	hashBytes := awsKms.HashRipemd160(awsKms.HashSha256(p.getCompressedPublicKey()))
	address := awsKms.EncodeAccountId(hashBytes)

	return address
}

func (p *XrpLocalSignerProvider) getEcdsaPublicKey() *ecdsa.PublicKey {
	publicKey := p.privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return publicKeyECDSA
}

func (p *XrpLocalSignerProvider) getCompressedPublicKey() []byte {
	return crypto.CompressPubkey(p.getEcdsaPublicKey())
}

func (p *XrpLocalSignerProvider) GetPublicKey() string {
	return hex.EncodeToString(p.getCompressedPublicKey())
}

func (p *XrpLocalSignerProvider) sign(payload []byte) (r, s, signature []byte) {
	signature, err := p.privateKey.Sign(rand.Reader, []byte(payload), nil)
	if err != nil {
		log.Error().Msgf("Error signing message: '%s'", err)
		return nil, nil, nil
	}

	var parsedSig struct{ R, S *big.Int }
	if _, err := asn1.Unmarshal(signature, &parsedSig); err != nil {
		log.Error().Msgf("asn1.Unmarshal: %s", err)
		return nil, nil, nil
	}

	// left pad R and S with zeroes
	rBytes := parsedSig.R.Bytes()
	sBytes := parsedSig.S.Bytes()

	// To avoid replay attack with inverse signature, only one half of the curve is allowed as a valid signature
	secp256k1N := crypto.S256().Params().N
	secp256k1halfN := new(big.Int).Div(secp256k1N, big.NewInt(2))
	sBI := new(big.Int).SetBytes(sBytes)

	// sBI > secp256k1halfN
	if sBI.Cmp(secp256k1halfN) == 1 {
		sBytes = sBI.Sub(secp256k1N, sBI).Bytes()
	}

	return rBytes, sBytes, signature
}

func (p *XrpLocalSignerProvider) SignMessage(payload string) string {
	hexBytes, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding message (%+v) to sign: '%+v'", payload, err)
		return ""
	}

	messageHash := awsKms.HashSha512(hexBytes)[0:32]
	rBytes, sBytes, _ := p.sign(messageHash)
	if rBytes == nil || sBytes == nil {
		return ""
	}

	return awsKms.BigIntBytesToSignatureHex(rBytes, sBytes)
}

func NewXrpLocalSignerProvider(cfg config.LocalSigner) *XrpLocalSignerProvider {
	ecdsaPrivateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatal().Msgf("Error creating ECDSA private key %+v", err)
	}
	return &XrpLocalSignerProvider{privateKey: ecdsaPrivateKey}
}
