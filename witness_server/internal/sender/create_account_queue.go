package sender

import (
	"peersyst/bridge-witness-go/internal/chains"
	"time"

	"github.com/rs/zerolog/log"
)

const MaxUnconfirmedAccountCreates = 100

func SendToCreateAccountQueue(provider chains.ChainProvider, transactionData TransactionData, nonce uint, gasFactor uint) {
	AppAttestationState.AddAttestation(provider.GetChainId().Uint64(), transactionData.Block, transactionData.Id)

	CreateAccountQueue <- &CreateAccountQueueItem{
		Provider:        provider,
		TransactionData: transactionData,
		GasFactor:       gasFactor,
		Nonce:           nonce}
}

func ProcessCreateAccountQueue(queueChannel <-chan *CreateAccountQueueItem) {
	// This queue will process createAccountCounts in order, we should only pass to next one when first is finished
QUEUE:
	for item := range queueChannel {
		for {
			currentCAC := item.Provider.GetCurrentCreateAccountCount()
			txCAC, err := item.Provider.GetTransactionCreateCount(item.TransactionData.Transaction)
			if err != nil {
				log.Error().Msgf("Error getting create account count from attestate create account tx: %+v", err)
				continue QUEUE
			}

			if txCAC <= currentCAC+MaxUnconfirmedAccountCreates {
				go SendTransaction(item.Provider, item.TransactionData, item.Nonce, item.GasFactor, 0)
				continue QUEUE
			}

			time.Sleep(1 * time.Second)
		}
	}
}
