package signer

import (
	"encoding/json"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"

	"github.com/rs/zerolog/log"
)

func EncodeXrpTransactionForSigning(transaction *transaction.Transaction) string {
	encodedTx := xrpl.GetXrplJs().EncodeForSigning(*transaction)
	return encodedTx
}

func EncodeXrpTransactionForMultiSigning(transaction *transaction.Transaction, signerAddress string) string {
	encodedTx := xrpl.GetXrplJs().EncodeForMultiSigning(*transaction, signerAddress)
	return encodedTx
}

func EncodeXrpTransaction(transaction *transaction.Transaction) string {
	result, err := json.Marshal(transaction)
	if err != nil {
		log.Error().Msgf("Error marshaling tx: '%+v'", err)
		return ""
	}
	encodedTx := xrpl.GetXrplJs().Encode(string(result))
	return encodedTx
}

func DecodeXrpTransaction(encodedTx string) transaction.Transaction {
	txString := xrpl.GetXrplJs().Decode(encodedTx)
	tx := transaction.TransactionStruct{}
	err := json.Unmarshal([]byte(txString), &tx)
	if err != nil {
		log.Error().Msgf("Error unmarshaling tx: '%+v'", err)
		return nil
	}
	return &tx
}
