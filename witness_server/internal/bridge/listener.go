package bridge

import (
	"errors"
	"peersyst/bridge-witness-go/internal/chains"

	"github.com/rs/zerolog/log"
)

// This function could take both chain providers (listener + executer)
// Or just have the chainProvider as a singleton and be able to get it in a way (like kms)
func fetchNewBridges(queueChannel <-chan interface{}) {
	for range queueChannel {
		mainChainProvider := chains.GetMainChainProvider()
		sideChainProvider := chains.GetSideChainProvider()

		err := fetchProviderNewBridges(mainChainProvider)
		if err != nil {
			log.Warn().Msgf("Error fetching new bridges: '%+v'", err)
		}

		err = fetchProviderNewBridges(sideChainProvider)
		if err != nil {
			log.Warn().Msgf("Error fetching new bridges: '%+v'", err)
		}

		chains.ValidateBridges()
	}
}

func fetchProviderNewBridges(provider chains.ChainProvider) error {
	currentBlock := provider.GetCurrentBlockNumber()
	if currentBlock == 0 {
		// Node is probably down
		return errors.New("invalid current block number")
	}

	err := provider.FetchNewBridges(currentBlock)
	if err != nil {
		// If fetch fails return to not set current number
		return err
	}

	// Set current block number as search has been successful
	provider.SetNewBridgesCurrentBlockNumber(currentBlock + 1)

	return nil
}
