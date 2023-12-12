package xrp

import (
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"testing"
)

func TestXrp_findEarliestBlockInTxs(t *testing.T) {
	var transactions []xrpl.TransactionAndMetadata
	transactions = append(transactions, xrpl.TransactionAndMetadata{Transaction: transaction.TransactionStruct{LedgerSequence: 500}})
	transactions = append(transactions, xrpl.TransactionAndMetadata{Transaction: transaction.TransactionStruct{LedgerSequence: 100}})
	transactions = append(transactions, xrpl.TransactionAndMetadata{Transaction: transaction.TransactionStruct{LedgerSequence: 25}})

	earliest := findEarliestBlockInTxs(transactions, 50)
	expected := 2 // 50 > 25 so the position is 2
	if *earliest != expected {
		t.Errorf("expected %+v got %+v", expected, *earliest)
	}
}
