package signer

import "math/big"

type SignerProvider interface {
	SignTransaction(payload string, opts interface{}) string
	SignMultiSigTransaction(payload string) string
	SignMessage(payload string) string
	GetAddress() string
	GetPublicKey() string
}

type SignEvmTransactionOpts struct {
	ChainId *big.Int
}
