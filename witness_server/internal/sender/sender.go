package sender

import (
	"peersyst/bridge-witness-go/internal/chains"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func SendTransaction(provider chains.ChainProvider, transactionData TransactionData, nonce uint, gasFactor uint, delay time.Duration) {
	AppAttestationState.AddAttestation(provider.GetChainId().Uint64(), transactionData.Block, transactionData.Id)
	// If we surpass gasFactor limit send noOp transaction
	if gasFactor > gasFactorLimit {
		log.Debug().Msgf("Gas factor surpassed the limit! Sending NoOp transaction with nonce %+v", nonce)
		transactionNoOp := provider.GetNoOpTransaction(nonce, gasFactor-1)

		if transactionNoOp == "" {
			// If transaction NoOp malformed wait and try again
			time.Sleep(time.Millisecond * delay)
			go SendTransaction(provider, transactionData, nonce, gasFactor-1, delay)
			return
		}
		transactionData.Transaction = transactionNoOp
	}
	time.Sleep(time.Millisecond * delay)
	log.Debug().Msgf("Enqueuing after sleeping Transaction %+v Nonce %+v GasFactor %+v", transactionData, nonce, gasFactor)
	BroadcastTransactionInQueue <- &BroadcastTransactionQueueItem{
		Provider:        provider,
		TransactionData: transactionData,
		GasFactor:       gasFactor,
		Nonce:           nonce}
}

func BroadcastTransaction(broadcastItem BroadcastTransactionQueueItem, hash string, expiresAt time.Time, delay time.Duration) {
	time.Sleep(time.Second * delay)
	TransactionStatusQueue <- TransactionStatusQueueItem{
		BroadcastTransactionQueueItem: broadcastItem,
		Hash:                          hash,
		ExpiresAt:                     expiresAt}
}

func ProcessBroadcastTransactionQueue(queueChannel <-chan *BroadcastTransactionQueueItem) {
	for broadcastTransactionQueueItem := range queueChannel {
		log.Info().Msgf("Processing broadcast transaction item %+v", broadcastTransactionQueueItem)
		currentNonce := broadcastTransactionQueueItem.Provider.GetNonce()
		if currentNonce == nil {
			// If can not get currentNonce send transaction again to queue
			go SendTransaction(
				broadcastTransactionQueueItem.Provider,
				broadcastTransactionQueueItem.TransactionData,
				broadcastTransactionQueueItem.Nonce,
				broadcastTransactionQueueItem.GasFactor,
				200,
			)
			continue
		}

		if *currentNonce <= broadcastTransactionQueueItem.Nonce {
			signedTx := broadcastTransactionQueueItem.Provider.SignTransaction(broadcastTransactionQueueItem.TransactionData.Transaction)
			if signedTx == "" {
				// Error signing transaction, requeue it
				go SendTransaction(
					broadcastTransactionQueueItem.Provider,
					broadcastTransactionQueueItem.TransactionData,
					broadcastTransactionQueueItem.Nonce,
					broadcastTransactionQueueItem.GasFactor,
					200,
				)
				continue
			}
			hash, err := broadcastTransactionQueueItem.Provider.BroadcastTransaction(signedTx)
			if err != nil {
				// Error broadcasting transaction
				var gasFactor uint
				var txData TransactionData

				if strings.Contains(err.Error(), chains.NoResponseError) {
					log.Warn().Msgf("Error broadcasting tx: NoResponseError %+v - Resending transaction with more gas %+v", err, broadcastTransactionQueueItem)
					txData = broadcastTransactionQueueItem.TransactionData
					gasFactor = broadcastTransactionQueueItem.GasFactor + 1
				} else if strings.Contains(err.Error(), chains.InvalidNonce) {
					log.Warn().Msgf("Error broadcasting tx: InvalidNonce %+v - Waiting for nonce %+v", err, broadcastTransactionQueueItem)
					txData = broadcastTransactionQueueItem.TransactionData
					gasFactor = broadcastTransactionQueueItem.GasFactor
				} else if strings.Contains(err.Error(), chains.IgnorableError) {
					log.Warn().Msgf("Error broadcasting tx: IgnorableError %+v - for transaction %+v", err, broadcastTransactionQueueItem)
					continue
				} else {
					log.Error().Msgf("Error broadcasting tx: %+v - for transaction %+v", err, broadcastTransactionQueueItem)
					gasFactor = broadcastTransactionQueueItem.GasFactor
					txData = TransactionData{
						Id:          broadcastTransactionQueueItem.TransactionData.Id,
						Block:       broadcastTransactionQueueItem.TransactionData.Block,
						Transaction: broadcastTransactionQueueItem.Provider.GetNoOpTransaction(broadcastTransactionQueueItem.Nonce, gasFactor),
					}
				}
				go SendTransaction(
					broadcastTransactionQueueItem.Provider,
					txData,
					broadcastTransactionQueueItem.Nonce,
					gasFactor,
					200,
				)
				continue
			}

			log.Debug().Msgf("Sending signed broadcast transaction %+v with hash %+v", broadcastTransactionQueueItem, hash)
			go BroadcastTransaction(*broadcastTransactionQueueItem, hash, time.Now().Add(time.Minute), 20)
		} else {
			log.Warn().Msgf("Ignoring transaction with past nonce %+v", broadcastTransactionQueueItem)
		}
	}
}

func ProcessTransactionStatusQueue(queueChannel <-chan TransactionStatusQueueItem) {
	for item := range queueChannel {
		log.Info().Msgf("Processing transaction status item %+v", item)
		status := item.Provider.GetTransactionStatus(item.Hash)
		log.Debug().Msgf("Status of transaction with hash %s is %s", item.Hash, status)

		if status == chains.AcceptedStatus {
			log.Info().Msgf("Transaction submitted correctly %s", item.Hash)
			AppAttestationState.SetAttested(item.Provider.GetChainId().Uint64(), item.TransactionData.Block, item.TransactionData.Id)
		} else if status == chains.PendingStatus {
			if item.ExpiresAt.Second() < time.Now().Second() {
				go SendTransaction(
					item.Provider,
					TransactionData{
						Id:          item.TransactionData.Id,
						Block:       item.TransactionData.Block,
						Transaction: item.Provider.GetNoOpTransaction(item.Nonce, item.GasFactor+1),
					},
					item.Nonce,
					item.GasFactor+1,
					0,
				)
			} else {
				item.incrementItemGasFactor()
				go SendTransaction(item.Provider, item.TransactionData, item.Nonce, item.GasFactor, 0)
			}
		} else if status == chains.NotFoundStatus || strings.Contains(status, chains.UnconfirmedStatus) {
			item.incrementItemGasFactor()
			go SendTransaction(item.Provider, item.TransactionData, item.Nonce, item.GasFactor+1, 0)
		} else if status == chains.FailedStatus {
			// Confirmed and failed: consumed nonce, report error
			log.Warn().Msgf("Transaction with hash %s has failed", item.Hash)
		} else {
			// Confirmed and unknown status due to request error
			go BroadcastTransaction(item.BroadcastTransactionQueueItem, item.Hash, item.ExpiresAt, 1)
		}
	}
}
