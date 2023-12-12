package xrp_awsKms

import (
	"encoding/hex"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"peersyst/bridge-witness-go/internal/signer"
	awsKms "peersyst/bridge-witness-go/internal/signer/aws_kms"

	"github.com/rs/zerolog/log"
)

var _ signer.SignerProvider = &XrpAwsKmsSignerProvider{}

type XrpAwsKmsSignerProvider struct {
	kmsService *awsKms.AwsKmsService
}

func (service *XrpAwsKmsSignerProvider) SignTransaction(payload string, opts interface{}) string {
	// Decode payload transaction
	tx, err := transaction.UnmarshalTransaction(payload)
	if err != nil {
		log.Error().Msgf("Error decoding xrp transaction: '%s'", err)
		return ""
	}
	tx.SetSigningPubKey(service.kmsService.GetPublicKey())

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
	rBytes, sBytes, _ := service.kmsService.Sign(encodedHash)
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

func (service *XrpAwsKmsSignerProvider) SignMultiSigTransaction(payload string) string {
	// Decode payload transaction
	tx := signer.DecodeXrpTransaction(payload)
	if tx == nil {
		log.Error().Msgf("Error decoding xrp transaction")
		return ""
	}
	tx.SetSigningPubKey("")

	encodeForMultiSigning := signer.EncodeXrpTransactionForMultiSigning(&tx, service.GetAddress())
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
	rBytes, sBytes, _ := service.kmsService.Sign(encodedHash)
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

func (service *XrpAwsKmsSignerProvider) GetAddress() string {
	compressedPubKey := service.kmsService.GetCompressedPublicKey()
	if compressedPubKey == nil {
		return ""
	}

	// Xrp address encoding: https://xrpl.org/img/address-encoding.svg
	hashBytes := awsKms.HashRipemd160(awsKms.HashSha256(compressedPubKey))
	address := awsKms.EncodeAccountId(hashBytes)

	return address
}

func (service *XrpAwsKmsSignerProvider) GetPublicKey() string {
	return service.kmsService.GetPublicKey()
}

func (service *XrpAwsKmsSignerProvider) SignMessage(payload string) string {
	hexBytes, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding message (%+v) to sign: '%+v'", payload, err)
		return ""
	}

	messageHash := awsKms.HashSha512(hexBytes)[0:32]
	rBytes, sBytes, _ := service.kmsService.Sign(messageHash)
	if rBytes == nil || sBytes == nil {
		return ""
	}

	return awsKms.BigIntBytesToSignatureHex(rBytes, sBytes)
}

func NewXrpAwsKmsSignerProvider(kmsService *awsKms.AwsKmsService) *XrpAwsKmsSignerProvider {
	return &XrpAwsKmsSignerProvider{kmsService: kmsService}
}
