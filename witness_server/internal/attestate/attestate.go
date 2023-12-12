package attestate

import (
	"math/rand"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/chains/evm"
	"peersyst/bridge-witness-go/internal/chains/xrp"
	"peersyst/bridge-witness-go/internal/sender"
	aws "peersyst/bridge-witness-go/internal/signer/aws_kms"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func resendToQueue(queue chan *interface{}, elem *interface{}) {
	time.Sleep(time.Second * 1)
	queue <- elem
}

func AttestateInSideChain(queueChannel <-chan *interface{}) {
	for pendingAttester := range queueChannel {
		log.Debug().Msgf("Pending attester data %+v", pendingAttester)
		claim, isClaim := (*pendingAttester).(*struct {
			Block       uint64
			ClaimId     uint64
			Sender      string
			Amount      string
			Destination string
			Nonce       int
			Fee         int
			BridgeId    string
		})
		accountCreate, isAccountCreate := (*pendingAttester).(*struct {
			Block           uint64
			Sender          string
			Amount          string
			Destination     string
			SignatureReward string
			Nonce           int
			Fee             int
			BridgeId        string
		})

		var attestTx string
		var nonce uint64
		var block uint64

		mainChainProvider := chains.GetMainChainProvider()
		sideChainProvider := chains.GetSideChainProvider()
		if !sideChainProvider.IsInSignerList() {
			continue
		}

		if isClaim {
			claimExists, err := checkSideChainClaim(claim)
			if err != nil {
				resendToQueue(AttestateInSideChainQueue, pendingAttester)
				continue
			}

			var senderSideChain string
			var destinationSideChain string
			if mainChainProvider.GetType() == sideChainProvider.GetType() {
				senderSideChain = claim.Sender
				destinationSideChain = claim.Destination
			} else if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Evm {
				senderSideChain = aws.XrplAccountToEvmAddress(claim.Sender)
				destinationSideChain = aws.XrplAccountToEvmAddress(claim.Destination)
			} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Xrp {
				senderSideChain = aws.EvmAddressToXrplAccount(claim.Sender)
				destinationSideChain = aws.EvmAddressToXrplAccount(claim.Destination)
			} else {
				continue
			}
			if !claimExists || senderSideChain == "" || destinationSideChain == "" {
				continue
			}

			amountParsed := sideChainProvider.ConvertToWhole(mainChainProvider.ConvertToDecimal(claim.Amount, claim.BridgeId), claim.BridgeId)
			attestTx, nonce = sideChainProvider.GetAttestClaimTransaction(claim.ClaimId, senderSideChain, amountParsed, destinationSideChain, claim.BridgeId)
			block = claim.Block
		} else if isAccountCreate {
			var sourceSideChain string
			var destinationSideChain string
			if mainChainProvider.GetType() == sideChainProvider.GetType() {
				sourceSideChain = accountCreate.Sender
				destinationSideChain = accountCreate.Destination
			} else if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Evm {
				sourceSideChain = aws.XrplAccountToEvmAddress(accountCreate.Sender)
				destinationSideChain = aws.XrplAccountToEvmAddress(accountCreate.Destination)
			} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Xrp {
				sourceSideChain = aws.EvmAddressToXrplAccount(accountCreate.Sender)
				destinationSideChain = aws.EvmAddressToXrplAccount(accountCreate.Destination)
			} else {
				continue
			}

			canCreateAccount, err := checkSideChainCreateAccount(destinationSideChain, accountCreate.BridgeId)
			if err != nil {
				resendToQueue(AttestateInSideChainQueue, pendingAttester)
				continue
			}
			if !canCreateAccount || destinationSideChain == "" || sourceSideChain == "" {
				continue
			}

			amountParsed := sideChainProvider.ConvertToWhole(mainChainProvider.ConvertToDecimal(accountCreate.Amount, accountCreate.BridgeId), accountCreate.BridgeId)
			sigRewardParsed := sideChainProvider.ConvertToWhole(mainChainProvider.ConvertToDecimal(accountCreate.SignatureReward, accountCreate.BridgeId), accountCreate.BridgeId)
			attestTx, nonce = sideChainProvider.GetAttestAccountCreateTransaction(sourceSideChain, amountParsed, destinationSideChain, sigRewardParsed, accountCreate.BridgeId)
			block = accountCreate.Block
		}

		if attestTx == "" {
			// If transaction fails to be constructed, requeue it
			resendToQueue(AttestateInSideChainQueue, pendingAttester)
			continue
		}

		go sender.SendTransaction(
			sideChainProvider,
			sender.TransactionData{Transaction: attestTx, Id: rand.Uint64(), Block: block},
			uint(nonce),
			1,
			0)
	}
}

func AttestateInMainChain(queueChannel <-chan *interface{}) {
	for pendingAttester := range queueChannel {
		claim, isClaim := (*pendingAttester).(*struct {
			Block       uint64
			ClaimId     uint64
			Sender      string
			Amount      string
			Destination string
			Nonce       int
			Fee         int
			BridgeId    string
		})
		accountCreate, isAccountCreate := (*pendingAttester).(*struct {
			Block           uint64
			Sender          string
			Amount          string
			Destination     string
			SignatureReward string
			Nonce           int
			Fee             int
			BridgeId        string
		})

		var attestTx string
		var nonce uint64
		var block uint64

		mainChainProvider := chains.GetMainChainProvider()
		sideChainProvider := chains.GetSideChainProvider()
		if !mainChainProvider.IsInSignerList() {
			continue
		}

		if isClaim {
			claimExists, err := checkMainChainClaim(claim)
			if err != nil {
				resendToQueue(AttestateInMainChainQueue, pendingAttester)
				continue
			}

			var senderMainChain string
			var destinationMainChain string
			if mainChainProvider.GetType() == sideChainProvider.GetType() {
				senderMainChain = claim.Sender
				destinationMainChain = claim.Destination
			} else if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Evm {
				senderMainChain = aws.EvmAddressToXrplAccount(claim.Sender)
				destinationMainChain = aws.EvmAddressToXrplAccount(claim.Destination)
			} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Xrp {
				senderMainChain = aws.XrplAccountToEvmAddress(claim.Sender)
				destinationMainChain = aws.XrplAccountToEvmAddress(claim.Destination)
			} else {
				continue
			}
			if !claimExists || senderMainChain == "" || destinationMainChain == "" {
				continue
			}

			amountParsed := mainChainProvider.ConvertToWhole(sideChainProvider.ConvertToDecimal(claim.Amount, claim.BridgeId), claim.BridgeId)
			attestTx, nonce = mainChainProvider.GetAttestClaimTransaction(claim.ClaimId, senderMainChain, amountParsed, destinationMainChain, claim.BridgeId)
			block = claim.Block
		} else if isAccountCreate {
			var sourceMainChain string
			var destinationMainChain string
			if mainChainProvider.GetType() == sideChainProvider.GetType() {
				sourceMainChain = accountCreate.Sender
				destinationMainChain = accountCreate.Destination
			} else if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Evm {
				sourceMainChain = aws.EvmAddressToXrplAccount(accountCreate.Sender)
				destinationMainChain = aws.EvmAddressToXrplAccount(accountCreate.Destination)
			} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Xrp {
				sourceMainChain = aws.XrplAccountToEvmAddress(accountCreate.Sender)
				destinationMainChain = aws.XrplAccountToEvmAddress(accountCreate.Destination)
			} else {
				continue
			}

			canCreateAccount, err := checkMainChainCreateAccount(destinationMainChain, accountCreate.BridgeId)
			if err != nil {
				resendToQueue(AttestateInMainChainQueue, pendingAttester)
				continue
			}
			if !canCreateAccount || destinationMainChain == "" || sourceMainChain == "" {
				continue
			}

			amountParsed := mainChainProvider.ConvertToWhole(sideChainProvider.ConvertToDecimal(accountCreate.Amount, accountCreate.BridgeId), accountCreate.BridgeId)
			sigRewardParsed := mainChainProvider.ConvertToWhole(sideChainProvider.ConvertToDecimal(accountCreate.SignatureReward, accountCreate.BridgeId), accountCreate.BridgeId)
			attestTx, nonce = mainChainProvider.GetAttestAccountCreateTransaction(sourceMainChain, amountParsed, destinationMainChain, sigRewardParsed, accountCreate.BridgeId)
			block = accountCreate.Block
		}
		if attestTx == "" {
			// If transaction fails to be constructed, requeue it
			resendToQueue(AttestateInMainChainQueue, pendingAttester)
			continue
		}

		if isClaim {
			go sender.SendTransaction(
				mainChainProvider,
				sender.TransactionData{Transaction: attestTx, Id: rand.Uint64(), Block: block},
				uint(nonce),
				1,
				0)
		} else if isAccountCreate {
			// If is AccountCreate send to accountCreateQueue
			sender.SendToCreateAccountQueue(
				mainChainProvider,
				sender.TransactionData{Transaction: attestTx, Id: rand.Uint64(), Block: block},
				uint(nonce),
				1,
			)
		}
	}
}

func checkMainChainClaim(commit *struct {
	Block       uint64
	ClaimId     uint64
	Sender      string
	Amount      string
	Destination string
	Nonce       int
	Fee         int
	BridgeId    string
}) (bool, error) {
	chainProvider := chains.GetMainChainProvider()
	claim, err := chainProvider.GetUnattestedClaimById(commit.ClaimId, commit.BridgeId)
	if err != nil {
		return false, err
	}

	if chainProvider.GetType() == config.Xrp {
		xrpClaim, isXrpClaim := claim.(xrp.XrpClaim)
		if !isXrpClaim {
			return false, nil
		}

		if chains.GetSideChainProvider().GetType() == config.Evm {
			sender := aws.EvmAddressToXrplAccount(commit.Sender)
			return xrpClaim.Source == sender, nil
		} else if chains.GetSideChainProvider().GetType() == config.Xrp {
			return xrpClaim.Source == commit.Sender, nil
		}

		return false, nil
	} else if chainProvider.GetType() == config.Evm {
		evmClaim, isEvmClaim := claim.(evm.EvmClaim)
		if !isEvmClaim {
			return false, nil
		}

		if chains.GetSideChainProvider().GetType() == config.Xrp {
			sender := aws.XrplAccountToEvmAddress(commit.Sender)
			return strings.EqualFold(evmClaim.Source, sender), nil
		} else if chains.GetSideChainProvider().GetType() == config.Evm {
			return strings.EqualFold(evmClaim.Source, commit.Sender), nil
		}

		return false, nil
	}

	return false, nil
}

func checkSideChainClaim(commit *struct {
	Block       uint64
	ClaimId     uint64
	Sender      string
	Amount      string
	Destination string
	Nonce       int
	Fee         int
	BridgeId    string
}) (bool, error) {
	chainProvider := chains.GetSideChainProvider()
	claim, err := chainProvider.GetUnattestedClaimById(commit.ClaimId, commit.BridgeId)
	if err != nil {
		return false, err
	}

	if chainProvider.GetType() == config.Xrp {
		xrpClaim, isXrpClaim := claim.(xrp.XrpClaim)
		if !isXrpClaim {
			return false, nil
		}

		if chains.GetMainChainProvider().GetType() == config.Evm {
			sender := aws.EvmAddressToXrplAccount(commit.Sender)
			return xrpClaim.Source == sender, nil
		} else if chains.GetMainChainProvider().GetType() == config.Xrp {
			return xrpClaim.Source == commit.Sender, nil
		}

		return false, nil
	} else if chainProvider.GetType() == config.Evm {
		evmClaim, isEvmClaim := claim.(evm.EvmClaim)
		if !isEvmClaim {
			return false, nil
		}

		if chains.GetMainChainProvider().GetType() == config.Xrp {
			sender := aws.XrplAccountToEvmAddress(commit.Sender)
			return strings.EqualFold(evmClaim.Source, sender), nil
		} else if chains.GetMainChainProvider().GetType() == config.Evm {
			return strings.EqualFold(evmClaim.Source, commit.Sender), nil
		}

		return false, nil
	}

	return false, nil
}

func checkMainChainCreateAccount(account, bridgeId string) (bool, error) {
	chainProvider := chains.GetMainChainProvider()

	accountCreated, err := chainProvider.CheckAccountCreated(account, bridgeId)
	if err != nil {
		return false, err
	}
	if accountCreated {
		return false, nil
	}

	alreadyAttested, err := chainProvider.CheckWitnessHasAttestedCreateAccount(account, bridgeId)
	if err != nil {
		return false, err
	}
	return !alreadyAttested, nil
}

func checkSideChainCreateAccount(account, bridgeId string) (bool, error) {
	chainProvider := chains.GetSideChainProvider()

	accountCreated, err := chainProvider.CheckAccountCreated(account, bridgeId)
	if err != nil {
		return false, err
	}
	if accountCreated {
		return false, nil
	}

	alreadyAttested, err := chainProvider.CheckWitnessHasAttestedCreateAccount(account, bridgeId)
	if err != nil {
		return false, err
	}
	return !alreadyAttested, nil
}
