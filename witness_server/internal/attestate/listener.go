package attestate

import (
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/chains/evm"
	"peersyst/bridge-witness-go/internal/chains/xrp"
)

func addToAttestateQueue(queueType QueueType, item *interface{}) {
	if queueType == mainChainQueue {
		AttestateInSideChainQueue <- item
	} else if queueType == sideChainQueue {
		AttestateInMainChainQueue <- item
	}
}

// This function could take both chain providers (listener + executer)
// Or just have the chainProvider as a singleton and be able to get it in a way (like kms)
func Fetch(queueChannel <-chan QueueType) {
	for queueType := range queueChannel {
		var chainProvider chains.ChainProvider
		if queueType == mainChainQueue {
			chainProvider = chains.GetMainChainProvider()
		} else if queueType == sideChainQueue {
			chainProvider = chains.GetSideChainProvider()
		}

		currentBlock := chainProvider.GetCurrentBlockNumber()
		if currentBlock == 0 {
			// Node is probably down, wait until next queue execution
			continue
		}

		commits := chainProvider.GetNewCommits(currentBlock)
		accCreates := chainProvider.GetNewAccountCreates(currentBlock)
		if commits == nil || accCreates == nil {
			// If either commits or accCreates are nil, continue before setting block number to search again (node down)
			continue
		}

		// Set current block number as search has been successful
		chainProvider.SetCurrentBlockNumber(currentBlock + 1)

		xrpCommits, isXrpCommit := commits.([]xrp.XrpCommit)
		xrpAccountCreates, isXrpAccountCreate := accCreates.([]xrp.XrpAccountCreate)
		evmCommits, isEvmCommit := commits.([]evm.EvmCommit)
		evmAccountCreates, isEvmAccountCreate := accCreates.([]evm.EvmAccountCreate)

		if isXrpCommit && isXrpAccountCreate && (len(xrpCommits) > 0 || len(xrpAccountCreates) > 0) {
			for _, xrpCommit := range xrpCommits {
				claim := getClaimFromXrpCommit(xrpCommit)
				addToAttestateQueue(queueType, &claim)
			}
			for _, xrpAccountCreate := range xrpAccountCreates {
				accCreate := getAccountCreateFromXrpAccCreate(xrpAccountCreate)
				addToAttestateQueue(queueType, &accCreate)
			}
		} else if isEvmCommit && isEvmAccountCreate && (len(evmCommits) > 0 || len(evmAccountCreates) > 0) {
			for _, evmCommit := range evmCommits {
				claim := getClaimFromEvmCommit(evmCommit)
				addToAttestateQueue(queueType, &claim)
			}
			for _, evmAccountCreate := range evmAccountCreates {
				accCreate := getAccountCreateFromEvmAccCreate(evmAccountCreate)
				addToAttestateQueue(queueType, &accCreate)
			}
		}
	}
}

func getClaimFromXrpCommit(commit xrp.XrpCommit) interface{} {
	var claim struct {
		Block       uint64
		ClaimId     uint64
		Sender      string
		Amount      string
		Destination string
		Nonce       int
		Fee         int
		BridgeId    string
	}
	claim.Block = commit.Block
	claim.ClaimId = commit.ClaimId
	claim.Sender = commit.Sender
	claim.Amount = commit.Amount
	if commit.Destination != nil {
		claim.Destination = *commit.Destination
	}
	claim.BridgeId = commit.BridgeId

	return &claim
}

func getAccountCreateFromXrpAccCreate(accCreate xrp.XrpAccountCreate) interface{} {
	var accountCreate struct {
		Block           uint64
		Sender          string
		Amount          string
		Destination     string
		SignatureReward string
		Nonce           int
		Fee             int
		BridgeId        string
	}
	accountCreate.Block = accCreate.Block
	accountCreate.Sender = accCreate.Sender
	accountCreate.Amount = accCreate.Amount
	accountCreate.Destination = accCreate.Destination
	accountCreate.SignatureReward = accCreate.SignatureReward
	accountCreate.BridgeId = accCreate.BridgeId

	return &accountCreate
}

func getClaimFromEvmCommit(commit evm.EvmCommit) interface{} {
	var claim struct {
		Block       uint64
		ClaimId     uint64
		Sender      string
		Amount      string
		Destination string
		Nonce       int
		Fee         int
		BridgeId    string
	}
	claim.Block = commit.Block
	claim.ClaimId = commit.ClaimId
	claim.Sender = commit.Sender
	claim.Amount = commit.Amount
	if commit.Destination != nil {
		claim.Destination = *commit.Destination
	}
	claim.BridgeId = commit.BridgeId

	return &claim
}

func getAccountCreateFromEvmAccCreate(accCreate evm.EvmAccountCreate) interface{} {
	var accountCreate struct {
		Block           uint64
		Sender          string
		Amount          string
		Destination     string
		SignatureReward string
		Nonce           int
		Fee             int
		BridgeId        string
	}
	accountCreate.Block = accCreate.Block
	accountCreate.Sender = accCreate.Sender
	accountCreate.Amount = accCreate.Amount
	accountCreate.Destination = accCreate.Destination
	accountCreate.SignatureReward = accCreate.SignatureReward
	accountCreate.BridgeId = accCreate.BridgeId

	return &accountCreate
}
