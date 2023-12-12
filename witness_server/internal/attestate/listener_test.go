package attestate

import (
	"peersyst/bridge-witness-go/internal/chains"
	"testing"
	"time"
)

func TestAttestate_Fetch(t *testing.T) {
	chains.StartXrpTestProvider(0, 2, true, nil, nil)
	chains.StartEvmTestProvider(0, 0, true, nil, nil)
	ListenerQueue = make(chan QueueType, 5)
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	go Fetch(ListenerQueue)

	ListenerQueue <- mainChainQueue
	ListenerQueue <- sideChainQueue
	time.Sleep(time.Millisecond)

	if chains.XrpTestProvider.GetNewCommitsCalledTimes != 0 {
		t.Errorf("error: Fetch should not continue if currentBlock is 0. expected %+v got %+v", 0, chains.XrpTestProvider.GetNewCommitsCalledTimes)
	}
	if chains.EvmTestProvider.GetNewCommitsCalledTimes != 0 {
		t.Errorf("error: Fetch should not continue if currentBlock is 0. expected %+v got %+v", 0, chains.XrpTestProvider.GetNewCommitsCalledTimes)
	}

	chains.XrpTestProvider.BlockNumber = 50
	chains.EvmTestProvider.BlockNumber = 50
	ListenerQueue <- mainChainQueue
	ListenerQueue <- sideChainQueue
	time.Sleep(time.Millisecond)

	if chains.XrpTestProvider.GetNewCommitsCalledTimes != 1 {
		t.Errorf("error: Fetch should continue if currentBlock is not 0. expected %+v got %+v", 1, chains.XrpTestProvider.GetNewCommitsCalledTimes)
	}
	if chains.XrpTestProvider.SetCurrentBlockCalledTimes > 0 {
		t.Errorf("error: Fetch should not continue if commits or accCreates are null. expected %+v got %+v", 0, chains.XrpTestProvider.SetCurrentBlockCalledTimes)
	}

	if chains.EvmTestProvider.GetNewCommitsCalledTimes != 1 {
		t.Errorf("error: Fetch should continue if currentBlock is not 0. expected %+v got %+v", 1, chains.XrpTestProvider.GetNewCommitsCalledTimes)
	}
	if chains.EvmTestProvider.SetCurrentBlockCalledTimes > 0 {
		t.Errorf("error: Fetch should not continue if commits or accCreates are null. expected %+v got %+v", 0, chains.XrpTestProvider.SetCurrentBlockCalledTimes)
	}

	chains.XrpTestProvider.BlockNumber = 100
	ListenerQueue <- mainChainQueue
	time.Sleep(time.Millisecond)
	if chains.XrpTestProvider.SetCurrentBlockCalledTimes != 1 {
		t.Errorf("error: Fetch should continue if commits or accCreates are not null. expected %+v got %+v", 1, chains.XrpTestProvider.SetCurrentBlockCalledTimes)
	}
	if len(AttestateInSideChainQueue) != 2 {
		t.Errorf("error: Fetch should send commits and account create to evm attestate. expected %+v got %+v", 2, len(AttestateInSideChainQueue))
	}

	chains.EvmTestProvider.BlockNumber = 1000
	ListenerQueue <- sideChainQueue
	time.Sleep(time.Millisecond)
	if chains.EvmTestProvider.SetCurrentBlockCalledTimes != 1 {
		t.Errorf("error: Fetch should continue if commits or accCreates are not null. expected %+v got %+v", 1, chains.EvmTestProvider.SetCurrentBlockCalledTimes)
	}
	if len(AttestateInMainChainQueue) != 2 {
		t.Errorf("error: Fetch should send commits and account create to xrp attestate. expected %+v got %+v", 2, len(AttestateInMainChainQueue))
	}
}
