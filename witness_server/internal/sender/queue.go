package sender

import (
	"peersyst/bridge-witness-go/internal/chains"
	"time"
)

type TransactionData struct {
	Id          uint64
	Transaction string
	Block       uint64
}

type CreateAccountQueueItem struct {
	Provider        chains.ChainProvider
	TransactionData TransactionData
	GasFactor       uint
	Nonce           uint // Priority
}

type BroadcastTransactionQueueItem struct {
	Provider        chains.ChainProvider
	TransactionData TransactionData
	GasFactor       uint
	Nonce           uint // Priority
	Index           int  // Index in the heap
}

type BroadcastPriorityQueue []*BroadcastTransactionQueueItem

type TransactionStatusQueueItem struct {
	BroadcastTransactionQueueItem
	Hash      string
	ExpiresAt time.Time
}

var (
	CreateAccountQueue           chan *CreateAccountQueueItem
	BroadcastTransactionInQueue  chan *BroadcastTransactionQueueItem
	BroadcastTransactionOutQueue chan *BroadcastTransactionQueueItem
	TransactionStatusQueue       chan TransactionStatusQueueItem
)

const gasFactorLimit uint = 10

func StartQueues() {
	CreateAccountQueue = make(chan *CreateAccountQueueItem, 3000)
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	BroadcastTransactionOutQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	TransactionStatusQueue = make(chan TransactionStatusQueueItem, 3000)
	go ProcessCreateAccountQueue(CreateAccountQueue)
	go ProcessInOutQueue(BroadcastTransactionInQueue, BroadcastTransactionOutQueue)
	go ProcessBroadcastTransactionQueue(BroadcastTransactionOutQueue)
	go ProcessTransactionStatusQueue(TransactionStatusQueue)
	go StartSavingState()
}
