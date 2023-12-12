package attestate

import "time"

type QueueType int32

const (
	mainChainQueue QueueType = 0
	sideChainQueue QueueType = 1
)

func StartQueues(queuePeriod int) {
	go startPeriodicQueue(queuePeriod)
}

var (
	ListenerQueue             chan QueueType
	AttestateInMainChainQueue chan *interface{}
	AttestateInSideChainQueue chan *interface{}
)

func startPeriodicQueue(period int) {
	ListenerQueue = make(chan QueueType, 5)
	AttestateInMainChainQueue = make(chan *interface{}, 1000)
	AttestateInSideChainQueue = make(chan *interface{}, 1000)
	ticker := time.NewTicker(time.Second * time.Duration(period))

	go AttestateInMainChain(AttestateInMainChainQueue)
	go AttestateInSideChain(AttestateInSideChainQueue)

	// Call before loop for immediate fetch
	// Calling go 2 times implies 2 workers for same channel which we send both types
	time.Sleep(time.Second)
	go Fetch(ListenerQueue)
	go Fetch(ListenerQueue)
	ListenerQueue <- mainChainQueue
	ListenerQueue <- sideChainQueue

	for range ticker.C {
		ListenerQueue <- mainChainQueue
		ListenerQueue <- sideChainQueue
	}
}
