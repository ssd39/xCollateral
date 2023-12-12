package bridge

import (
	"peersyst/bridge-witness-go/internal/bridge/sequencer"
	"time"

	"github.com/rs/zerolog/log"
)

type QueueType int32

const (
	mainChainQueue QueueType = 0
	sideChainQueue QueueType = 1
)

func StartListenerQueue(queuePeriod int) {
	go startListenerPeriodicQueue(queuePeriod)
}

func StartCreationQueue(queuePeriod int, sequencerUrl, mcDoorAddress, scDoorAddress string, minBridgeReward, maxBridgeReward, maxBridgeIterations uint64) {
	clt, err := sequencer.NewClientWithResponses(sequencerUrl)
	if err != nil {
		log.Fatal().Msgf("Error setting up sequencer client: %+v", err)
	}

	// Setup variables
	client = clt
	minBridgeCreateReward = minBridgeReward
	maxBridgeCreateReward = maxBridgeReward
	mainchainDoorAddress = mcDoorAddress
	sidechainDoorAddress = scDoorAddress
	maxBridgeCreateIterations = maxBridgeIterations

	go startCreationPeriodicQueue(queuePeriod)
}

var (
	ListenerQueue             chan interface{}
	CreationQueue             chan QueueType
	client                    *sequencer.ClientWithResponses
	minBridgeCreateReward     uint64
	maxBridgeCreateReward     uint64
	maxBridgeCreateIterations uint64
	mainchainDoorAddress      string
	sidechainDoorAddress      string
)

func startListenerPeriodicQueue(period int) {
	ListenerQueue = make(chan interface{}, 5)
	ticker := time.NewTicker(time.Second * time.Duration(period))

	time.Sleep(time.Second)
	go fetchNewBridges(ListenerQueue)

	for range ticker.C {
		ListenerQueue <- 1
	}
}

func startCreationPeriodicQueue(period int) {
	CreationQueue = make(chan QueueType, 5)
	ticker := time.NewTicker(time.Second * time.Duration(period))

	time.Sleep(time.Second)
	go fetchBridgeCreationRequests(CreationQueue)

	for range ticker.C {
		CreationQueue <- mainChainQueue
		CreationQueue <- sideChainQueue
	}
}
