package evm

import (
	"crypto/ecdsa"
	"encoding/hex"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/signer"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

var _ signer.SignerProvider = &EvmLocalSignerProvider{}

type EvmLocalSignerProvider struct {
	privateKey *ecdsa.PrivateKey
}

func (p *EvmLocalSignerProvider) SignTransaction(payload string, opts interface{}) string {
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

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(parsedOpts.ChainId), p.privateKey)
	if err != nil {
		log.Error().Msgf("Error signing transaction: %v", err)
	}

	// Encode transaction
	rawSignedTxBytes, err2 := rlp.EncodeToBytes(signedTx)
	if err2 != nil {
		log.Error().Msgf("Error encoding transaction to bytes : '%s'", err2)
		return ""
	}

	return hex.EncodeToString(rawSignedTxBytes)
}

func (p *EvmLocalSignerProvider) SignMultiSigTransaction(payload string) string {
	return payload
}

func (p *EvmLocalSignerProvider) GetAddress() string {
	return crypto.PubkeyToAddress(*p.getEcdsaPublicKey()).Hex()
}

func (p *EvmLocalSignerProvider) getEcdsaPublicKey() *ecdsa.PublicKey {
	publicKey := p.privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	return publicKeyECDSA
}

func (p *EvmLocalSignerProvider) GetPublicKey() string {
	publicKeyBytes := crypto.CompressPubkey(p.getEcdsaPublicKey())
	return hex.EncodeToString(publicKeyBytes)
}

func (p *EvmLocalSignerProvider) SignMessage(payload string) string {
	messageHash, err := signer.HashEvmMessage(payload)
	if err != nil {
		log.Error().Msgf("Error hashing the payload %+v", err)
		return ""
	}

	signature, err := crypto.Sign(messageHash, p.privateKey)
	if err != nil {
		log.Error().Msgf("Error signing transaction %+v", err)
		return ""
	}
	return hex.EncodeToString(signer.NormalizeSignature(signature))
}

func NewEvmLocalSignerProvider(cfg config.LocalSigner) *EvmLocalSignerProvider {
	ecdsaPrivateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatal().Msgf("Error creating ECDSA private key %+v", err)
	}
	return &EvmLocalSignerProvider{privateKey: ecdsaPrivateKey}
}
