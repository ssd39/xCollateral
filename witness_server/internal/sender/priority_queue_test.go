package sender

import (
	"peersyst/bridge-witness-go/internal/chains"
	"testing"
	"time"
)

func addToBroadCastTxInQueue(nonce uint) {
	BroadcastTransactionInQueue <- &BroadcastTransactionQueueItem{
		Provider: chains.XrpTestProvider,
		TransactionData: TransactionData{
			Id: uint64(nonce), Block: 1000, Transaction: "transactionData",
		},
		GasFactor: 1,
		Nonce:     nonce,
	}
}

func TestSender_ProcessInOutQueue(t *testing.T) {
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	BroadcastTransactionOutQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	go ProcessInOutQueue(BroadcastTransactionInQueue, BroadcastTransactionOutQueue)

	time.Sleep(time.Millisecond)
	addToBroadCastTxInQueue(3)
	addToBroadCastTxInQueue(8)
	addToBroadCastTxInQueue(4)
	addToBroadCastTxInQueue(1)
	addToBroadCastTxInQueue(2)
	addToBroadCastTxInQueue(12)
	addToBroadCastTxInQueue(5)
	time.Sleep(time.Millisecond)

	if len(BroadcastTransactionInQueue) != 0 {
		t.Errorf("error: in queue should be empty expected %+v got %+v", 0, len(BroadcastTransactionInQueue))
	}
	if len(BroadcastTransactionOutQueue) != 7 {
		t.Errorf("error: out queue should have expected %+v got %+v", 7, len(BroadcastTransactionOutQueue))
	}
	go func() {
		failCounter := 0
		lastNonce := uint(0)
		for item := range BroadcastTransactionOutQueue {
			if lastNonce > item.Nonce {
				failCounter += 1
			}
			lastNonce = item.Nonce
		}
		if failCounter > 3 {
			t.Errorf("error: too many non sorted nonces %+v currentNonce %+v", 3, failCounter)
		}
	}()

	time.Sleep(time.Millisecond * 50)
	close(BroadcastTransactionOutQueue)
}
