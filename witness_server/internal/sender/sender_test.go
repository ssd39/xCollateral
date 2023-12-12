package sender

import (
	"math/big"
	"peersyst/bridge-witness-go/internal/chains"
	"testing"
	"time"
)

func TestSender_SendTransaction(t *testing.T) {
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	chains.StartXrpTestProvider(150, 150, true, big.NewInt(144), nil)

	go SendTransaction(chains.GetMainChainProvider(), TransactionData{
		Id: 1, Block: 200, Transaction: "TransactionData",
	}, 1, 1, 1)
	time.Sleep(time.Millisecond * 2)

	if chains.XrpTestProvider.GetNoOpTransactionCalledTimes != 0 {
		t.Errorf("error: get no op transaction should be called expected %+v got %+v", 0, chains.XrpTestProvider.GetNoOpTransactionCalledTimes)
	}
	if len(BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: in queue should have elem expected %+v got %+v", 1, len(BroadcastTransactionInQueue))
	}

	go SendTransaction(chains.GetMainChainProvider(), TransactionData{
		Id: 2, Block: 200, Transaction: "TransactionData",
	}, 2, gasFactorLimit+1, 1)
	time.Sleep(time.Millisecond * 2)

	if chains.XrpTestProvider.GetNoOpTransactionCalledTimes != 1 {
		t.Errorf("error: get no op transaction should be called expected %+v got %+v", 1, chains.XrpTestProvider.GetNoOpTransactionCalledTimes)
	}
	if len(BroadcastTransactionInQueue) != 2 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 2, len(BroadcastTransactionInQueue))
	}

	go SendTransaction(chains.GetMainChainProvider(), TransactionData{
		Id: 4, Block: 200, Transaction: "TransactionData",
	}, 0, gasFactorLimit+1, 1)
	go SendTransaction(chains.GetMainChainProvider(), TransactionData{
		Id: 3, Block: 200, Transaction: "TransactionData",
	}, 3, 1, 1)
	time.Sleep(time.Millisecond * 3)

	if chains.XrpTestProvider.GetNoOpTransactionCalledTimes != 2 {
		t.Errorf("error: get no op transaction should be called expected %+v got %+v", 2, chains.XrpTestProvider.GetNoOpTransactionCalledTimes)
	}
	if len(BroadcastTransactionInQueue) != 4 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 4, len(BroadcastTransactionInQueue))
	}

	// Id 4 should go after Id 3 in in queue (even if 4 got called first)
	go func() {
		counter := uint64(1)
		for item := range BroadcastTransactionInQueue {
			if counter != item.TransactionData.Id {
				t.Errorf("error: expected %+v got %+v", counter, item.TransactionData.Id)
			}
			counter += 1
		}
	}()

	time.Sleep(time.Millisecond)
	close(BroadcastTransactionInQueue)
}

func TestSender_BroadcastTransaction(t *testing.T) {
	TransactionStatusQueue = make(chan TransactionStatusQueueItem, 3000)

	go BroadcastTransaction(BroadcastTransactionQueueItem{}, "hash", time.Now().Add(time.Minute), 1)
	time.Sleep(time.Millisecond)
	if len(TransactionStatusQueue) != 0 {
		t.Errorf("error: queue should be empty expected %+v got %+v", 0, len(TransactionStatusQueue))
	}

	time.Sleep(time.Second)
	if len(TransactionStatusQueue) != 1 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 1, len(TransactionStatusQueue))
	}

	close(TransactionStatusQueue)
}

func createMockTxQueueItem(transaction string) *BroadcastTransactionQueueItem {
	return &BroadcastTransactionQueueItem{
		Provider:  chains.GetMainChainProvider(),
		GasFactor: 1,
		Nonce:     10,
		TransactionData: TransactionData{
			Id:          1,
			Block:       200,
			Transaction: transaction,
		},
	}
}

func TestSender_ProcessBroadcastTransactionQueue(t *testing.T) {
	chains.StartXrpTestProvider(150, 150, true, big.NewInt(144), nil)
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}
	BroadcastTransactionOutQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	TransactionStatusQueue = make(chan TransactionStatusQueueItem, 3000)
	go ProcessBroadcastTransactionQueue(BroadcastTransactionOutQueue)

	BroadcastTransactionOutQueue <- createMockTxQueueItem("transaction")
	time.Sleep(time.Millisecond * 210)
	if len(BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 1, len(BroadcastTransactionInQueue))
	}

	currentNonce := uint(10)
	chains.XrpTestProvider.Nonce = &currentNonce
	BroadcastTransactionOutQueue <- createMockTxQueueItem("fail")
	time.Sleep(time.Millisecond * 210)
	if len(BroadcastTransactionInQueue) != 2 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 2, len(BroadcastTransactionInQueue))
	}

	BroadcastTransactionOutQueue <- createMockTxQueueItem("timeout")
	time.Sleep(time.Millisecond * 210)
	if len(BroadcastTransactionInQueue) != 3 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 3, len(BroadcastTransactionInQueue))
	}

	BroadcastTransactionOutQueue <- createMockTxQueueItem("invalid nonce")
	time.Sleep(time.Millisecond * 210)
	if len(BroadcastTransactionInQueue) != 4 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 4, len(BroadcastTransactionInQueue))
	}

	BroadcastTransactionOutQueue <- createMockTxQueueItem("ignorable")
	time.Sleep(time.Millisecond * 210)
	// ignorable error is not sent to any queue
	if len(BroadcastTransactionInQueue) != 4 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 4, len(BroadcastTransactionInQueue))
	}
	if len(TransactionStatusQueue) != 0 {
		t.Errorf("error: queue should be empty expected %+v got %+v", 0, len(TransactionStatusQueue))
	}

	BroadcastTransactionOutQueue <- createMockTxQueueItem("unknown")
	time.Sleep(time.Millisecond * 210)
	if len(BroadcastTransactionInQueue) != 5 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 5, len(BroadcastTransactionInQueue))
	}
	if len(TransactionStatusQueue) != 0 {
		t.Errorf("error: queue should be empty expected %+v got %+v", 0, len(TransactionStatusQueue))
	}

	// Check every error sent params well
	go func() {
		counter := uint64(1)
		for item := range BroadcastTransactionInQueue {
			if counter != 3 && item.GasFactor != 1 {
				t.Errorf("error: gas factor should stay the same expected %+v got %+v", 1, item.GasFactor)
			}
			if counter == 3 && item.GasFactor != 2 {
				t.Errorf("error: gas factor augment expected %+v got %+v", 2, item.GasFactor)
			}
			if counter == 5 && item.TransactionData.Transaction != "noOpTransactionEncoded" {
				t.Errorf("error: should change to noop expected %+v got %+v", "noOpTransactionEncoded", item.TransactionData.Transaction)
			}
			counter += 1
		}
	}()

	BroadcastTransactionOutQueue <- createMockTxQueueItem("success")
	time.Sleep(time.Second * 21)
	if len(TransactionStatusQueue) != 1 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 1, len(TransactionStatusQueue))
	}

	time.Sleep(time.Millisecond)
	close(BroadcastTransactionInQueue)
}

func TestSender_ProcessTransactionStatusQueue(t *testing.T) {
	chainId := big.NewInt(144).Uint64()
	chains.StartXrpTestProvider(150, 150, true, big.NewInt(144), nil)
	AppAttestationState = AttestationState{
		LastAttestedBlocks: make(LastAttestedBlocksState),
		BlockAttestations:  make(BlockAttestationsState),
	}

	TransactionStatusQueue = make(chan TransactionStatusQueueItem, 3000)
	BroadcastTransactionInQueue = make(chan *BroadcastTransactionQueueItem, 3000)
	go ProcessTransactionStatusQueue(TransactionStatusQueue)

	AppAttestationState.AddAttestation(chainId, 200, 1)
	broadcastItem := BroadcastTransactionQueueItem{
		Provider:  chains.GetMainChainProvider(),
		GasFactor: 1,
		Nonce:     1,
		TransactionData: TransactionData{
			Id:          1,
			Block:       200,
			Transaction: "transaction",
		},
	}
	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Minute),
		Hash:                          chains.AcceptedStatus,
	}
	time.Sleep(time.Millisecond)
	if AppAttestationState.BlockAttestations[chainId][200][1] != true {
		t.Errorf("error: expected set attested %+v got %+v", true, AppAttestationState.BlockAttestations[chainId][200][1])
	}

	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Minute),
		Hash:                          chains.PendingStatus,
	}
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 1 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 1, len(BroadcastTransactionInQueue))
	}
	if chains.XrpTestProvider.SetTransactionGasPriceCalledTimes != 1 {
		t.Errorf("error: should have called set gas expected %+v got %+v", 1, chains.XrpTestProvider.SetTransactionGasPriceCalledTimes)
	}

	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Second * -2),
		Hash:                          chains.PendingStatus,
	}
	chains.XrpTestProvider.SetTransactionGasPriceCalledTimes = 0
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 2 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 2, len(BroadcastTransactionInQueue))
	}
	if chains.XrpTestProvider.SetTransactionGasPriceCalledTimes != 0 {
		t.Errorf("error: should not have called set gas expected %+v got %+v", 0, chains.XrpTestProvider.SetTransactionGasPriceCalledTimes)
	}

	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Minute),
		Hash:                          chains.NotFoundStatus,
	}
	chains.XrpTestProvider.SetTransactionGasPriceCalledTimes = 0
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 3 {
		t.Errorf("error: queue should have elem expected %+v got %+v", 3, len(BroadcastTransactionInQueue))
	}
	if chains.XrpTestProvider.SetTransactionGasPriceCalledTimes != 1 {
		t.Errorf("error: should have called set gas expected %+v got %+v", 1, chains.XrpTestProvider.SetTransactionGasPriceCalledTimes)
	}

	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Minute),
		Hash:                          chains.FailedStatus,
	}
	chains.XrpTestProvider.SetTransactionGasPriceCalledTimes = 0
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 3 {
		t.Errorf("error: queue should have same elem expected %+v got %+v", 3, len(BroadcastTransactionInQueue))
	}
	if chains.XrpTestProvider.SetTransactionGasPriceCalledTimes != 0 {
		t.Errorf("error: should not have called set gas expected %+v got %+v", 0, chains.XrpTestProvider.SetTransactionGasPriceCalledTimes)
	}
	if len(TransactionStatusQueue) != 0 {
		t.Errorf("error: queue should be empty expected %+v got %+v", 0, len(TransactionStatusQueue))
	}

	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		ExpiresAt:                     time.Now().Add(time.Minute),
		Hash:                          chains.ConfirmedStatus,
	}
	chains.XrpTestProvider.SetTransactionGasPriceCalledTimes = 0
	time.Sleep(time.Millisecond)
	if len(BroadcastTransactionInQueue) != 3 {
		t.Errorf("error: queue should have same elem expected %+v got %+v", 3, len(BroadcastTransactionInQueue))
	}
	if chains.XrpTestProvider.SetTransactionGasPriceCalledTimes != 0 {
		t.Errorf("error: should not have called set gas expected %+v got %+v", 0, chains.XrpTestProvider.SetTransactionGasPriceCalledTimes)
	}

	TransactionStatusQueue = make(chan TransactionStatusQueueItem, 3000)
	time.Sleep(time.Second)
	if len(TransactionStatusQueue) != 1 {
		t.Errorf("error: queue should be empty expected %+v got %+v", 1, len(TransactionStatusQueue))
	}
}
