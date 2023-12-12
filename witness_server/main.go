package main

import (
	"os"
	"os/signal"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/attestate"
	"peersyst/bridge-witness-go/internal/bridge"
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/common"
	"peersyst/bridge-witness-go/internal/oracle"
	"peersyst/bridge-witness-go/internal/sender"
	"peersyst/bridge-witness-go/internal/signer/factory"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	configFilePath := ""
	if len(os.Args) == 2 {
		configFilePath = os.Args[1:][0]
	}
	// Load config
	conf := config.LoadConfig(configFilePath)

	common.InitLogger(conf.LogFilePath, conf.LoggingLevel, conf.LogFormat)

	// Start mainChain provider
	mainChainSigner := factory.NewSignerProviderFromConfig(conf.MainChain.Type, conf.MainChain)
	mainChainProvider, err := chains.StartMainChainProvider(conf.MainChain, mainChainSigner)
	if err != nil {
		log.Fatal().Msgf("Error instantiating mainChain provider : '%s'", err)
	}
	log.Info().Msgf("MainChain provider : '%+v'", mainChainProvider)

	// Start sideChain provider
	sideChainSigner := factory.NewSignerProviderFromConfig(conf.SideChain.Type, conf.SideChain)
	sideChainProvider, err := chains.StartSideChainProvider(conf.SideChain, sideChainSigner)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error instantiating sideChain provider : '%s'", err)
	}
	log.Info().Msgf("SideChain provider : '%+v'", sideChainProvider)

	// Validate bridge
	validated := chains.ValidateBridges()
	if !validated {
		log.Fatal().Msgf("Invalid bridge error")
	}

	// Recover saved attestation state
	attestationState := sender.LoadAttestationState()
	block, found := (*attestationState).LastAttestedBlocks[mainChainProvider.GetChainId().Uint64()]
	if found && block > 0 {
		log.Info().Msgf("Recovering sidechain last attested block %v", block)
		sideChainProvider.SetCurrentBlockNumber(block)
	}
	block, found = (*attestationState).LastAttestedBlocks[sideChainProvider.GetChainId().Uint64()]
	if found && block > 0 {
		log.Info().Msgf("Recovering mainchain last attested block %v", block)
		mainChainProvider.SetCurrentBlockNumber(block)
	}

	// Start xrp and evm listener queues
	log.Info().Msgf("Starting queues...")
	attestate.StartQueues(conf.Server.QueuePeriod)
	bridge.StartListenerQueue(conf.Server.BridgeListenerQueuePeriod)
	if conf.Server.DynamicBridgeCreation {
		bridge.StartCreationQueue(
			conf.Server.BridgeCreationQueuePeriod,
			conf.Server.SequencerUrl,
			conf.MainChain.DoorAddress,
			conf.SideChain.DoorAddress,
			conf.Server.MinBridgeSignatureReward,
			conf.Server.MaxBridgeSignatureReward,
			conf.Server.MaxCreateBridgeIterations,
		)
	}
	sender.StartQueues()
	go oracle.StartPriceOracle()
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	log.Info().Msgf("Server started successfully")
	<-done // Will block here until user hits ctrl+c
}
