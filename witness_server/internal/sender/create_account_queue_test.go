package sender

import (
	"math/big"
	"peersyst/bridge-witness-go/internal/chains"
	"testing"
	"time"
)

func TestSender_ProcessCreateAccountQueue(t *testing.T) {
	CreateAccountQueue = make(chan *CreateAccountQueueItem, 3000)
	go ProcessCreateAccountQueue(CreateAccountQueue)
	chains.StartXrpTestProvider(100, 49, true, big.NewInt(144), nil)
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}

	CreateAccountQueue <- &CreateAccountQueueItem{
		Provider: chains.XrpTestProvider,
		TransactionData: TransactionData{
			Id: 1, Transaction: "error", Block: 100,
		},
		GasFactor: 10,
		Nonce:     1}
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 0 {
		t.Errorf("error: tx should not be sent if error expected %+v got %+v", 0, len(BroadcastTransactionInQueue))
	}
	if len(CreateAccountQueue) != 0 {
		t.Errorf("error: create queue should be empty expected %+v got %+v", 0, len(CreateAccountQueue))
	}

	CreateAccountQueue <- &CreateAccountQueueItem{
		Provider: chains.XrpTestProvider,
		TransactionData: TransactionData{
			Id: 1, Transaction: "real-transaction", Block: 100,
		},
		GasFactor: 10,
		Nonce:     1}
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 0 {
		t.Errorf("error: tx should not be sent if maxUnconfirmed expected %+v got %+v", 1, len(BroadcastTransactionInQueue))
	}
	if len(CreateAccountQueue) != 0 {
		t.Errorf("error: create queue should be empty expected %+v got %+v", 0, len(CreateAccountQueue))
	}

	chains.XrpTestProvider.AccountCount = 50
	time.Sleep(time.Second)
	if len(BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: tx should be sent expected %+v got %+v", 1, len(BroadcastTransactionInQueue))
	}
	if len(CreateAccountQueue) != 0 {
		t.Errorf("error: create queue should be empty expected %+v got %+v", 0, len(CreateAccountQueue))
	}

	CreateAccountQueue <- &CreateAccountQueueItem{
		Provider: chains.XrpTestProvider,
		TransactionData: TransactionData{
			Id: 1, Transaction: "real-transaction", Block: 100,
		},
		GasFactor: 10,
		Nonce:     1}
	time.Sleep(time.Millisecond * 2)
	if len(BroadcastTransactionInQueue) != 2 {
		t.Errorf("error: tx should be sent expected %+v got %+v", 2, len(BroadcastTransactionInQueue))
	}
	if len(CreateAccountQueue) != 0 {
		t.Errorf("error: create queue should be empty expected %+v got %+v", 0, len(CreateAccountQueue))
	}

	close(CreateAccountQueue)
	close(BroadcastTransactionInQueue)
}
