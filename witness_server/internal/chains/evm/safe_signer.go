package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"peersyst/bridge-witness-go/internal/signer"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type SafeSigner struct {
	safe   Safe
	signer signer.SignerProvider
}

type SafeSignature struct {
	signer    string
	publicKey string
	data      string
}

type SafeTransactionJSON struct {
	To             common.Address `json:"to,omitempty"`
	Value          string         `json:"value,omitempty"`
	Data           string         `json:"data,omitempty"`
	Operation      uint8          `json:"operation,omitempty"`
	SafeTxGas      string         `json:"safeTxGas,omitempty"`
	BaseGas        string         `json:"baseGas,omitempty"`
	GasPrice       string         `json:"gasPrice,omitempty"`
	GasToken       common.Address `json:"gasToken,omitempty"`
	RefundReceiver common.Address `json:"refundReceiver,omitempty"`
	Nonce          int64          `json:"nonce,omitempty"`
}

type SafeTransaction struct {
	To             common.Address
	Value          *big.Int
	Data           []byte
	Operation      uint8
	SafeTxGas      *big.Int
	BaseGas        *big.Int
	GasPrice       *big.Int
	GasToken       common.Address
	RefundReceiver common.Address
	Nonce          *big.Int
}

func (ss *SafeSigner) SignTransactionHash(hash string) (*SafeSignature, error) {
	signature := ss.signer.SignMessage(hash)
	hexSignature, err := hex.DecodeString(signature)
	if err != nil {
		return nil, err
	}
	hexSignature[len(hexSignature)-1] = hexSignature[len(hexSignature)-1] + 4

	return &SafeSignature{
		signer:    ss.signer.GetAddress(),
		publicKey: ss.signer.GetPublicKey(),
		data:      "0x" + hex.EncodeToString(hexSignature),
	}, nil
}

func (ss *SafeSigner) GetTransactionHash(safeTransaction SafeTransaction) ([32]byte, error) {
	callOpts := bind.CallOpts{Pending: true, From: common.HexToAddress(ss.signer.GetAddress()), Context: context.Background()}
	return ss.safe.GetTransactionHash(
		&callOpts,
		safeTransaction.To,
		safeTransaction.Value,
		safeTransaction.Data,
		safeTransaction.Operation,
		safeTransaction.SafeTxGas,
		safeTransaction.BaseGas,
		safeTransaction.GasPrice,
		safeTransaction.GasToken,
		safeTransaction.RefundReceiver,
		safeTransaction.Nonce,
	)
}

func parseSafeTransactionInJSON(jsonTx string) (*SafeTransaction, error) {
	var tx SafeTransactionJSON
	err := json.Unmarshal([]byte(jsonTx), &tx)
	if err != nil {
		return nil, err
	}

	value, ok := new(big.Int).SetString(tx.Value, 0)
	if !ok {
		return nil, errors.New("error parsing value as big int")
	}
	safeTxGas, ok := new(big.Int).SetString(tx.SafeTxGas, 0)
	if !ok {
		return nil, errors.New("error parsing safeTxGas as big int")
	}
	baseGas, ok := new(big.Int).SetString(tx.BaseGas, 0)
	if !ok {
		return nil, errors.New("error parsing baseGas as big int")
	}
	gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 0)
	if !ok {
		return nil, errors.New("error parsing gasPrice as big int")
	}

	dataParsed := tx.Data
	if strings.HasPrefix(tx.Data, "0x") {
		dataParsed = tx.Data[2:]
	}
	data, err := hex.DecodeString(dataParsed)
	if err != nil {
		return nil, err
	}

	safeTx := SafeTransaction{
		To:             tx.To,
		Value:          value,
		Data:           data,
		Operation:      tx.Operation,
		SafeTxGas:      safeTxGas,
		BaseGas:        baseGas,
		GasPrice:       gasPrice,
		GasToken:       tx.GasToken,
		RefundReceiver: tx.RefundReceiver,
		Nonce:          big.NewInt(tx.Nonce),
	}

	return &safeTx, nil
}
