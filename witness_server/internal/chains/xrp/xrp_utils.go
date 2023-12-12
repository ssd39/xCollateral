package xrp

import (
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"sort"
)

func findEarliestBlockInTxs(transactions []xrpl.TransactionAndMetadata, block uint64) *int {
	i := sort.Search(len(transactions), func(i int) bool { return transactions[i].Transaction.GetLedgerSequence() <= block })
	return &i
}
