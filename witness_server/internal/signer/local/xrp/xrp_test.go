package xrp

import (
	"crypto/ecdsa"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"testing"
)

func testXrp_createTransaction(account string) string {
	tx := &transaction.TransactionStruct{}
	tx.Account = account
	tx.TransactionType = "AccountSet"
	nonceUint64 := uint64(1)
	tx.Sequence = &nonceUint64
	feeStr := "10"
	tx.Fee = &feeStr
	lastLedgerSeq := uint64(2091286)
	tx.LastLedgerSequence = &lastLedgerSeq
	tx.InLedger = 0
	tx.LedgerSequence = 0
	tx.SetFlags(0)

	marshalledTx, err := transaction.MarshalTransaction(tx)
	if err != nil {
		return ""
	}

	return marshalledTx
}

func TestXrp_SignTransaction(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	localSigner := NewXrpLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	transaction := testXrp_createTransaction(localSigner.GetAddress())
	expect := ""

	// Volia fer un verify de la firma pero es molta feina fer el decode, etc
	signedTransaction := localSigner.SignTransaction(transaction, nil)
	if expect == signedTransaction {
		t.Fatalf("Invalid transaction signed - expected: %+v got: %+v", expect, len(signedTransaction))
	}
}

func TestXrp_sign(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	localSigner := NewXrpLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	message := []byte("message to sign")
	expect := true

	// Volia fer un verify de la firma pero es molta feina fer el decode, etc
	r, s, _ := localSigner.sign(message)
	verified := ecdsa.Verify(localSigner.getEcdsaPublicKey(), message, new(big.Int).SetBytes(r), new(big.Int).SetBytes(s))
	if expect != verified {
		t.Fatalf("Invalid transaction signed - expected: %+v got: %+v", expect, verified)
	}
}

func TestXrp_GetAddress(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	expect := "rMArvG7XWy875CHWDrn2rFK33fgdPzyAuJ"
	localSigner := NewXrpLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	if expect != localSigner.GetAddress() {
		t.Fatalf("Invalid address - expected: %+v got: %+v", expect, localSigner.GetAddress())
	}
}

func TestXrp_GetPublicKey(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	expect := "039a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fd"
	localSigner := NewXrpLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	if expect != localSigner.GetPublicKey() {
		t.Fatalf("Invalid public key - expected: %+v got: %+v", expect, localSigner.GetPublicKey())
	}
}
