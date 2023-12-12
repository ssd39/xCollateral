package bridge

import (
	"context"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/bridge/sequencer"
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/chains/evm"
	aws "peersyst/bridge-witness-go/internal/signer/aws_kms"
	"strconv"

	"github.com/rs/zerolog/log"
)

func fetchBridgeCreationRequests(queueChannel <-chan QueueType) {
	for queueType := range queueChannel {
		var lockingChainProvider, issuingChainProvider chains.ChainProvider
		var lockingDoor, issuingDoor string

		if queueType == mainChainQueue {
			lockingChainProvider = chains.GetMainChainProvider()
			issuingChainProvider = chains.GetSideChainProvider()
			lockingDoor = mainchainDoorAddress
			issuingDoor = sidechainDoorAddress
		} else if queueType == sideChainQueue {
			lockingChainProvider = chains.GetSideChainProvider()
			issuingChainProvider = chains.GetMainChainProvider()
			lockingDoor = sidechainDoorAddress
			issuingDoor = mainchainDoorAddress
		}

		if !lockingChainProvider.IsInSignerList() || !issuingChainProvider.IsInSignerList() {
			continue
		}

		currentBlock := lockingChainProvider.GetCurrentBlockNumber()
		if currentBlock == 0 {
			// Node is probably down, wait until next queue execution
			continue
		}

		newBridgeRequests, err := lockingChainProvider.FetchNewBridgeRequests(currentBlock)
		if err != nil {
			log.Warn().Msgf("Error fetching new bridge requests: %+v", err)
			continue
		}

		evmBridgeRequests, isEvmRequest := newBridgeRequests.([]*evm.BridgeRequestCounter)
		if isEvmRequest {
			for _, evmBridgeRequest := range evmBridgeRequests {
				log.Debug().Msgf("Bridge request event: %+v", evmBridgeRequest.Event)

				// Call sequencer to get mounted transaction
				txId := evmBridgeRequest.Event.Raw.TxHash.String() + "-" + strconv.FormatUint(uint64(evmBridgeRequest.Event.Raw.Index), 10)
				resp, err := client.GetCreateBridgeWithResponse(context.Background(), txId)
				if err != nil {
					log.Error().Msgf("Error fetching create bridge transaction from sequencer: %+v", err)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
				if resp.StatusCode() >= 300 {
					log.Error().Msgf("Invalid sequencer response status: %+v", resp.StatusCode())
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}

				bodyResp := resp.JSON200

				// Validate chains types
				lockingEvm := bodyResp.LockingChainType == sequencer.CreateBridgeEntryDtoLockingChainTypeEvm
				lockingType := lockingChainProvider.GetType()
				issuingEvm := bodyResp.IssuingChainType == sequencer.CreateBridgeEntryDtoIssuingChainTypeEvm
				issuingType := issuingChainProvider.GetType()
				if lockingEvm && lockingType != config.Evm || !lockingEvm && lockingType == config.Evm {
					log.Error().Msgf("Invalid locking chain type, expected: %s got %s", bodyResp.LockingChainType, lockingType)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
				if issuingEvm && issuingType != config.Evm || !issuingEvm && issuingType == config.Evm {
					log.Error().Msgf("Invalid issuing chain type, expected: %s got %s", bodyResp.IssuingChainType, issuingType)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}

				// Convert addresses if necessary
				tokenAddress := evmBridgeRequest.Event.TokenAddress.String()
				tokenAddressInIssChain := tokenAddress
				lockingDoorInIssuing := lockingDoor
				issuingDoorInLocking := issuingDoor
				if lockingType == config.Evm && issuingType == config.Xrp {
					lockingDoorInIssuing = aws.EvmAddressToXrplAccount(lockingDoor)
					issuingDoorInLocking = aws.XrplAccountToEvmAddress(issuingDoor)
					tokenAddressInIssChain = aws.EvmAddressToXrplAccount(tokenAddress)
				} else if lockingType == config.Xrp && issuingType == config.Evm { // NOT POSSIBLE AS WE ARE LOOPING EVM EVENTS
					// lockingDoor = aws.XrplAccountToEvmAddress(lockingDoor)
					// issuingDoor = aws.EvmAddressToXrplAccount(issuingDoor)
					// tokenAddressInIssChain = aws.XrplAccountToEvmAddress(tokenAddress)
					log.Warn().Msgf("Locking is XRP and issuing EVM: should not happen, ignoring")
					continue
				} else if lockingType != issuingType {
					// NOT SUPPORTED SITUATION (and currently impossible)
					log.Warn().Msgf("Did we implement a 3rd chain and not implement here?")
					continue
				}

				// Get token codes
				tokenCode, err := lockingChainProvider.GetTokenCodeFromAddress(tokenAddress)
				if err != nil {
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}

				// Sign and validate locking chain transaction
				account, publicKey, signature, err := lockingChainProvider.SignEncodedCreateBridgeTransaction(
					bodyResp.LockingChainDoorCreateBridgeTransaction,
					true,
					minBridgeCreateReward,
					maxBridgeCreateReward,
					issuingDoorInLocking,
					tokenAddress,
					tokenCode,
				)
				if err != nil {
					log.Error().Msgf("Error signing locking chain create bridge transaction: '%s'", err)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
				lockingSignedTx := sequencer.SignedTransactionRequest{Account: account, PublicKey: publicKey, Signature: signature}

				// Sign and validate issuing chain transaction
				account, publicKey, signature, err = issuingChainProvider.SignEncodedCreateBridgeTransaction(
					bodyResp.IssuingChainDoorCreateBridgeTransaction,
					false,
					minBridgeCreateReward,
					maxBridgeCreateReward,
					lockingDoorInIssuing,
					tokenAddressInIssChain,
					tokenCode,
				)
				if err != nil {
					log.Error().Msgf("Error signing issuing chain create bridge transaction: '%s'", err)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
				issuingSignedTx := sequencer.SignedTransactionRequest{Account: account, PublicKey: publicKey, Signature: signature}

				// Send signed transaction to sequencer
				bodyReq := sequencer.SignCreateBridgeJSONRequestBody{
					SignedLockingChainDoorCreateBridgeTransaction: lockingSignedTx,
					SignedIssuingChainDoorCreateBridgeTransaction: issuingSignedTx,
				}
				respSign, err := client.SignCreateBridgeWithResponse(context.Background(), txId, bodyReq)
				if err != nil {
					log.Error().Msgf("Error sending signed transactions to sequencer: %s", err)
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
				if respSign.StatusCode() == 409 {
					log.Warn().Msgf("Sequencer status already signed, passing to next")
					continue
				}
				if respSign.StatusCode() >= 300 {
					log.Error().Msgf("Invalid sequencer response status: %d - %+v", respSign.StatusCode(), string(respSign.Body))
					retryIfNotBigger(lockingChainProvider, evmBridgeRequest.Tries, evmBridgeRequest)
					continue
				}
			}
		}
	}
}

func retryIfNotBigger(provider chains.ChainProvider, tries uint, request interface{}) {
	if tries < uint(maxBridgeCreateIterations) {
		err := provider.RetryNewBridgeRequest(request)
		if err != nil {
			log.Error().Msgf("Error retrying new bridge request: %+v", err)
		}
	}
}
