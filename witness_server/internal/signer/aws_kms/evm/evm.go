package evm

import (
	"encoding/hex"
	"peersyst/bridge-witness-go/internal/signer"
	awsKms "peersyst/bridge-witness-go/internal/signer/aws_kms"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

var _ signer.SignerProvider = &EvmAwsKmsSignerProvider{}

type EvmAwsKmsSignerProvider struct {
	kmsService *awsKms.AwsKmsService
}

func (service *EvmAwsKmsSignerProvider) SignTransaction(payload string, opts interface{}) string {
	// Parse chainId
	parsedOpts, ok := opts.(*signer.SignEvmTransactionOpts)
	if !ok {
		log.Error().Msg("Invalid nil chainId on sign evm transaction")
		return ""
	}

	// Decode transaction
	rawTxBytes, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding transaction: '%s'", err)
		return ""
	}
	tx := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		log.Error().Msgf("Error decoding transaction in rlp: '%s'", err)
		return ""
	}

	// Prepare transaction signing params
	pubkey := service.kmsService.GetUncompressedPublicKey()
	if pubkey == nil {
		// Error logged on kms
		return ""
	}
	pubKeyBytes := secp256k1.S256().Marshal(pubkey.X, pubkey.Y)
	txSigner := types.LatestSignerForChainID(parsedOpts.ChainId)
	txHashBytes := txSigner.Hash(tx).Bytes()

	// Sign transaction
	rBytes, sBytes, _ := service.kmsService.Sign(txHashBytes)
	if rBytes == nil || sBytes == nil {
		// Error logged on kms
		return ""
	}
	signature := awsKms.GetEthereumSignature(pubKeyBytes, txHashBytes, rBytes, sBytes)
	if signature == nil {
		return ""
	}

	signedTx, err := tx.WithSignature(txSigner, signature)
	if err != nil {
		log.Error().Msgf("Error signing evm transaction : '%s'", err)
		return ""

	}

	// Encode transaction
	rawSignedTxBytes, err2 := rlp.EncodeToBytes(signedTx)
	if err2 != nil {
		log.Error().Msgf("Error encoding transaction to bytes : '%s'", err2)
		return ""
	}

	return hex.EncodeToString(rawSignedTxBytes)
}

func (service *EvmAwsKmsSignerProvider) SignMultiSigTransaction(payload string) string {
	return payload
}

func (service *EvmAwsKmsSignerProvider) GetAddress() string {
	ecdsaPubKey := service.kmsService.GetUncompressedPublicKey()
	if ecdsaPubKey == nil {
		return ""
	}

	address := crypto.PubkeyToAddress(*ecdsaPubKey)
	return address.String()
}

func (service *EvmAwsKmsSignerProvider) GetPublicKey() string {
	return service.kmsService.GetPublicKey()
}

func (service *EvmAwsKmsSignerProvider) SignMessage(payload string) string {
	messageHash, err := signer.HashEvmMessage(payload)
	if err != nil {
		log.Error().Msgf("Error hashing the payload %+v", err)
		return ""
	}

	rBytes, sBytes, _ := service.kmsService.Sign(messageHash)
	if rBytes == nil || sBytes == nil {
		return ""
	}

	pubkey := service.kmsService.GetUncompressedPublicKey()
	if pubkey == nil {
		// Error logged on kms
		return ""
	}
	pubKeyBytes := secp256k1.S256().Marshal(pubkey.X, pubkey.Y)
	signature := awsKms.GetEthereumSignature(pubKeyBytes, messageHash, rBytes, sBytes)
	if signature == nil {
		return ""
	}
	if err != nil {
		log.Error().Msgf("Error transforming signature %+v", err)
		return ""
	}
	return hex.EncodeToString(signer.NormalizeSignature(signature))
}

func NewEvmAwsKmsSignerProvider(kmsService *awsKms.AwsKmsService) *EvmAwsKmsSignerProvider {
	return &EvmAwsKmsSignerProvider{kmsService: kmsService}
}
