package xrp

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	external "peersyst/bridge-witness-go/external/xrpl.js"
	"peersyst/bridge-witness-go/internal/common/cache"
	"peersyst/bridge-witness-go/internal/common/utils"
	"peersyst/bridge-witness-go/internal/signer"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"

	"github.com/rs/zerolog/log"
)

type XrpProvider struct {
	node                       string
	witnessAddress             string
	doorAddress                string
	currentBlock               uint64
	currentNewBridgesBlock     uint64
	currentBridgeRequestsBlock uint64
	client                     *xrpl.Client
	sequence                   uint64
	signerProvider             signer.SignerProvider
	inSignerList               *bool
	lastSignerCheck            time.Time
	recheckSignerDuration      time.Duration
	networkId                  uint64
	maxGasFactor               int64
	bridgeProviders            map[string]*XrpBridgeProvider
	unpairedBridgeProviders    map[string]*XrpBridgeProvider
}

type XrpCommit struct {
	BridgeId    string
	Block       uint64
	ClaimId     uint64
	Sender      string
	Amount      string
	Destination *string
}

type XrpAccountCreate struct {
	BridgeId        string
	Block           uint64
	Sender          string
	Amount          string
	Destination     string
	SignatureReward string
}

type XrpClaim struct {
	ClaimId uint64
	Sender  string
	Source  string
}

var XrplJs *external.XrplJs

const (
	MaxBlocksPerRequest = 5000
	MaxTxsPerRequest    = 10000
	OldestBlockDiff     = 100000
	XrpPrec             = 6
	TokenPrec           = 15
)

func Create(signerProvider signer.SignerProvider, node string, doorAddress string, startingBlock uint64, signerListSeconds, maxGasFactor int64) (*XrpProvider, error) {
	XrplJs = external.NewXrplJs()
	client, err := xrpl.Create(node)
	if err != nil {
		return nil, err
	}

	var currentSeq uint64 = 0

	// Get network id
	serverInfo, err := client.GetServerInfo()
	if err != nil {
		log.Error().Msgf("Error getting server info: '%+v'", err)
		return nil, err
	}

	// If starting block is not specified, fetch the latest validated from the network
	currentBlock, err := client.GetLedgerIndex()
	if err != nil {
		log.Error().Msgf("Error getting block number: '%+v'", err)
		return nil, err
	}
	if startingBlock == 0 {
		startingBlock = currentBlock
	}

	// Retrieve current Seq
	ledgerIndex := "current"
	lastSeq, err := client.GetAccountInfo(signerProvider.GetAddress(), &ledgerIndex)
	if err != nil {
		log.Error().Msgf("Error getting account info: '%+v'", err)
	} else {
		currentSeq = lastSeq.AccountData.Sequence
	}

	// Retrieve current account create count from bridge
	objectType := "bridge"
	bridges := map[string]*XrpBridgeProvider{}
	bridgeObjects, err := client.GetAccountObjects(doorAddress, &ledgerIndex, &objectType)
	if err != nil {
		log.Error().Msgf("Error getting bridge account objects: '%+v'", err)
	} else {
		for _, object := range bridgeObjects.Objects {
			jsonObj, _ := json.Marshal(object)
			bridgeObj := xrpl.XChainBridgeObject{}
			if err := json.Unmarshal(jsonObj, &bridgeObj); err == nil && bridgeObj.LedgerEntryType == "Bridge" {
				createCount, err := strconv.ParseUint(bridgeObj.XChainAccountClaimCount, 16, 64)
				if err != nil {
					createCount = 0
				}
				bridgeProvider := CreateXrpBridgeProvider(doorAddress, bridgeObj.XChainBridge, createCount)
				bridges[bridgeProvider.bridgeId] = bridgeProvider
			}
		}
	}

	provider := XrpProvider{
		node,
		signerProvider.GetAddress(),
		doorAddress,
		startingBlock,
		currentBlock,
		currentBlock,
		client,
		currentSeq,
		signerProvider,
		nil,
		time.Now(),
		time.Duration(signerListSeconds),
		serverInfo.Info.NetworkId,
		maxGasFactor,
		map[string]*XrpBridgeProvider{},
		bridges,
	}
	return &provider, nil
}

func (provider *XrpProvider) UpdateOracleData(amount int64, amount2 int64) error {
	return nil
}

func (provider *XrpProvider) BroadcastTransaction(signedTx string) (string, error) {
	txResult, err := provider.client.Submit(signedTx)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return "", fmt.Errorf("no response error: %+v", err)
		}
		return "", fmt.Errorf("unknown error: %+v", err)
	}
	if txResult.EngineResult != "tesSUCCESS" {
		if txResult.EngineResult == "terPRE_SEQ" {
			return "", fmt.Errorf("invalid nonce")
		}
		if txResult.EngineResult == "tecXCHAIN_NO_CLAIM_ID" ||
			txResult.EngineResult == "tecXCHAIN_ACCOUNT_CREATE_PAST" ||
			txResult.EngineResult == "terQUEUED" ||
			txResult.EngineResult == "tefPAST_SEQ" {
			return "", fmt.Errorf("ignorable error: %+v", txResult.EngineResult)
		}
		return "", fmt.Errorf("unknown error: %+v", txResult.EngineResult)
	}

	return txResult.Tx.Hash, nil
}

func (provider *XrpProvider) GetAttestClaimTransaction(claimId uint64, sender, amount, destination, bridgeId string) (string, uint64) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return "", 0
	}

	// Prepare the part of the attestation that needs to be signed internally
	tx := &transaction.TransactionStruct{}
	tx.XChainBridge = bridgeProvider.bridge
	tx.OtherChainSource = &sender
	if bridgeProvider.isLocking && bridgeProvider.bridge.LockingChainIssue.Issuer != nil {
		tx.Amount = transaction.TokenAmount{Currency: bridgeProvider.bridge.IssuingChainIssue.Currency, Issuer: *bridgeProvider.bridge.IssuingChainIssue.Issuer, Value: amount}
	} else if !bridgeProvider.isLocking && bridgeProvider.bridge.IssuingChainIssue.Issuer != nil {
		tx.Amount = transaction.TokenAmount{Currency: bridgeProvider.bridge.LockingChainIssue.Currency, Issuer: *bridgeProvider.bridge.LockingChainIssue.Issuer, Value: amount}
	} else {
		tx.Amount = amount
	}
	tx.AttestationRewardAccount = &provider.witnessAddress
	wasLockingChainSend := uint64(0)
	if !bridgeProvider.isLocking {
		wasLockingChainSend = uint64(1)
	}
	tx.WasLockingChainSend = &wasLockingChainSend
	claimIdHex := strconv.FormatUint(claimId, 16)
	tx.XChainClaimID = &claimIdHex
	if destination != "" {
		tx.Destination = &destination
	}
	log.Debug().Msgf("Transaction to encode %+v", tx)
	// Sign the internal attestation
	jsonTx, err := transaction.MarshalTransaction(tx)
	log.Debug().Msgf("JSON transaction %+v", jsonTx)
	if err != nil {
		return "", 0
	}
	encoded := xrpl.GetXrplJs().Encode(jsonTx)
	if encoded == "" {
		return "", 0
	}
	signature := provider.signerProvider.SignMessage(encoded)
	if signature == "" {
		return "", 0
	}
	tx.Signature = &signature
	tx.AttestationSignerAccount = &provider.witnessAddress

	// Add the rest of the transaction fields
	tx.TransactionType = "XChainAddClaimAttestation"
	tx.Account = provider.witnessAddress
	publicKey := provider.signerProvider.GetPublicKey()
	if publicKey == "" {
		return "", 0
	}
	tx.PublicKey = &publicKey
	seq := provider.consumeSequence()
	tx.Sequence = &seq

	// Autofill remaining fields
	autoFilledTx := provider.client.Autofill(tx)
	if autoFilledTx == nil {
		log.Error().Msgf("Error autofilling tx: %+v", tx)
		return "", 0
	}

	marshalledTx, err := transaction.MarshalTransaction(autoFilledTx)
	if err != nil {
		log.Error().Msgf("Error marshaling tx: '%s'", err)
		return "", 0
	}

	return marshalledTx, seq
}

func (provider *XrpProvider) GetAttestAccountCreateTransaction(sender, amount, destination, signatureReward, bridgeId string) (string, uint64) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return "", 0
	}

	// Prepare the part of the attestation that needs to be signed internally
	tx := &transaction.TransactionStruct{}
	tx.XChainBridge = bridgeProvider.bridge
	tx.OtherChainSource = &sender
	tx.Amount = amount
	tx.AttestationRewardAccount = &provider.witnessAddress
	wasLockingChainSend := uint64(0)
	if !bridgeProvider.isLocking {
		wasLockingChainSend = uint64(1)
	}
	tx.WasLockingChainSend = &wasLockingChainSend
	accountCreateCount := strconv.FormatUint(bridgeProvider.consumeAccountCreateCount(), 16)
	tx.XChainAccountCreateCount = &accountCreateCount
	tx.Destination = &destination
	tx.SignatureReward = signatureReward

	// Sign the internal attestation
	jsonTx, err := transaction.MarshalTransaction(tx)
	if err != nil {
		return "", 0
	}
	encoded := xrpl.GetXrplJs().Encode(jsonTx)
	if encoded == "" {
		return "", 0
	}
	signature := provider.signerProvider.SignMessage(encoded)
	if signature == "" {
		return "", 0
	}
	tx.Signature = &signature
	tx.AttestationSignerAccount = &provider.witnessAddress

	// Add the rest of the transaction fields
	tx.TransactionType = "XChainAddAccountCreateAttestation"
	tx.Account = provider.witnessAddress
	publicKey := provider.signerProvider.GetPublicKey()
	if publicKey == "" {
		return "", 0
	}
	tx.PublicKey = &publicKey
	seq := provider.consumeSequence()
	tx.Sequence = &seq

	// Autofill remaining fields
	autoFilledTx := provider.client.Autofill(tx)
	if autoFilledTx == nil {
		log.Error().Msgf("Error autofilling tx: %+v", tx)
		return "", 0
	}

	marshalledTx, err := transaction.MarshalTransaction(autoFilledTx)
	if err != nil {
		log.Error().Msgf("Error marshaling tx: '%s'", err)
		return "", 0
	}

	return marshalledTx, seq
}

func (provider *XrpProvider) consumeSequence() uint64 {
	seq := atomic.AddUint64(&provider.sequence, 1)
	return seq - 1
}

func (provider *XrpProvider) SetBridgeValidated(bridgeId string) interface{} {
	bridgeProvider, exists := provider.unpairedBridgeProviders[bridgeId]
	if !exists {
		return nil
	}

	delete(provider.unpairedBridgeProviders, bridgeId)
	provider.bridgeProviders[bridgeId] = bridgeProvider

	return nil
}

func (provider *XrpProvider) GetUnpairedBridges() interface{} {
	return provider.unpairedBridgeProviders
}

func (provider *XrpProvider) SignTransaction(payload string) string {
	return provider.signerProvider.SignTransaction(payload, struct{}{})
}

func (provider *XrpProvider) SignEncodedCreateBridgeTransaction(encodedTx string, isLocking bool, minBridgeReward, maxBridgeReward uint64, otherChainAddress, tokenAddress, tokenCode string) (string, string, string, error) {
	tx := signer.DecodeXrpTransaction(encodedTx)
	if tx.GetAccount() != provider.doorAddress {
		return "", "", "", errors.New("error create bridge transaction account (" + tx.GetAccount() + ") should be door address: " + provider.doorAddress)
	}
	if tx.GetTransactionType() != "XChainCreateBridge" {
		return "", "", "", errors.New("error create bridge transaction type (" + tx.GetTransactionType() + ") should be 'XChainCreateBridge'")
	}
	if tx.GetSigningPubKey() != "" {
		return "", "", "", errors.New("error create bridge transaction signing pub key (" + tx.GetSigningPubKey() + ") should be ''")
	}

	sigReward, err := strconv.Atoi(tx.GetSignatureReward())
	if err != nil {
		return "", "", "", err
	}
	minBridgeRewardWhole := utils.FloatToIntPrec(big.NewFloat(float64(minBridgeReward)), XrpPrec)
	if int64(sigReward) < minBridgeRewardWhole.Int64() {
		return "", "", "", errors.New("error create bridge transaction signature reward (" + tx.GetSignatureReward() + ") should be more than or equal than minBridgeReward: " + string(rune(minBridgeReward)))
	}
	maxBridgeRewardWhole := utils.FloatToIntPrec(big.NewFloat(float64(maxBridgeReward)), XrpPrec)
	if int64(sigReward) > maxBridgeRewardWhole.Int64() {
		return "", "", "", errors.New("error create bridge transaction signature reward (" + tx.GetSignatureReward() + ") should be less than or equal than maxBridgeReward: " + string(rune(maxBridgeReward)))
	}

	fee, err := strconv.Atoi(*tx.GetFee())
	if err != nil {
		return "", "", "", err
	}
	if fee > int(provider.maxGasFactor) {
		return "", "", "", errors.New("error create bridge transaction fee (" + *tx.GetFee() + ") should be less than or equal than maxGasFactor: " + string(rune(provider.maxGasFactor)))
	}

	bridge := tx.GetXChainBridge()
	if bridge == nil {
		return "", "", "", errors.New("error create bridge transaction undefined XChainBridge")
	}
	if isLocking {
		if bridge.LockingChainDoor != provider.doorAddress {
			return "", "", "", errors.New("error create bridge transaction lockingChainDoor (" + bridge.LockingChainDoor + ") should be: " + provider.doorAddress)
		}
		if bridge.IssuingChainDoor != otherChainAddress {
			return "", "", "", errors.New("error create bridge transaction issuingChainDoor (" + bridge.IssuingChainDoor + ") should be: " + otherChainAddress)
		}
		if *bridge.IssuingChainIssue.Issuer != otherChainAddress {
			return "", "", "", errors.New("error create bridge transaction issuingChainIssue Issuer (" + *bridge.IssuingChainIssue.Issuer + ") should be: " + otherChainAddress)
		}
	} else {
		if bridge.LockingChainDoor != otherChainAddress {
			return "", "", "", errors.New("error create bridge transaction lockingChainDoor (" + bridge.LockingChainDoor + ") should be: " + otherChainAddress)
		}
		if bridge.IssuingChainDoor != provider.doorAddress {
			return "", "", "", errors.New("error create bridge transaction issuingChainDoor (" + bridge.IssuingChainDoor + ") should be: " + provider.doorAddress)
		}
		if *bridge.IssuingChainIssue.Issuer != provider.doorAddress {
			return "", "", "", errors.New("error create bridge transaction issuingChainIssue Issuer (" + *bridge.IssuingChainIssue.Issuer + ") should be: " + provider.doorAddress)
		}
	}
	if *bridge.LockingChainIssue.Issuer != tokenAddress {
		return "", "", "", errors.New("error create bridge transaction lockingChainIssue Issuer (" + *bridge.LockingChainIssue.Issuer + ") should be: " + tokenAddress)
	}

	lockingChainCurrency := bridge.LockingChainIssue.Currency
	if len(bridge.LockingChainIssue.Currency) != 3 {
		lockCurrencyByte, err := hex.DecodeString(lockingChainCurrency)
		if err != nil {
			return "", "", "", err
		}
		lockingChainCurrency = string(bytes.Trim(lockCurrencyByte, "\x00"))
	}
	if lockingChainCurrency != tokenCode {
		return "", "", "", errors.New("error create bridge transaction lockingChainIssue Currency (" + lockingChainCurrency + ") should be: " + tokenCode)
	}

	issuingChainCurrency := bridge.IssuingChainIssue.Currency
	if len(bridge.IssuingChainIssue.Currency) != 12 {
		issCurrencyByte, err := hex.DecodeString(issuingChainCurrency)
		if err != nil {
			return "", "", "", err
		}
		issuingChainCurrency = string(bytes.Trim(issCurrencyByte, "\x00"))
	}
	issuingStrings := strings.Split(issuingChainCurrency, "-")
	if len(issuingStrings) != 2 {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency (" + issuingChainCurrency + ") should be splitted to 2 by -")
	}
	if issuingStrings[0] != tokenCode {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency first part (" + issuingStrings[0] + ") should be: " + tokenCode)
	}
	_, err = hex.DecodeString(issuingStrings[1])
	if err != nil || len(issuingStrings[1]) != 8 {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency second part (" + issuingStrings[1] + ") should be have length 8")
	}

	signature := provider.signerProvider.SignMultiSigTransaction(encodedTx)

	return provider.signerProvider.GetAddress(), provider.signerProvider.GetPublicKey(), signature, nil
}

func (provider *XrpProvider) GetNoOpTransaction(nonce uint, gasPrice uint) string {
	tx := &transaction.TransactionStruct{}
	tx.Account = provider.witnessAddress
	tx.TransactionType = "AccountSet"
	nonceUint64 := uint64(nonce)
	tx.Sequence = &nonceUint64
	feeStr := strconv.FormatUint(uint64(gasPrice)*10, 10)
	tx.Fee = &feeStr

	// Autofill remaining fields
	autoFilledTx := provider.client.Autofill(tx)
	if autoFilledTx == nil {
		log.Error().Msgf("Error autofilling tx")
		return ""
	}

	marshalledTx, err := transaction.MarshalTransaction(autoFilledTx)
	if err != nil {
		log.Error().Msgf("Error marshaling tx: '%s'", err)
		return ""
	}

	return marshalledTx
}

func (provider *XrpProvider) SetTransactionGasPrice(payload string, factor uint) string {
	tx, err := transaction.UnmarshalTransaction(payload)
	if err != nil {
		log.Error().Msgf("Error unmarshaling tx: '%s'", err)
		return ""
	}

	currentFee, err := strconv.ParseUint(*tx.GetFee(), 10, 64)
	if err != nil {
		log.Error().Msgf("Error parsing current fee: '%s'", err)
		return ""
	}
	if factor > 1 {
		tx.SetFee(strconv.FormatUint(currentFee*uint64(factor)/uint64(factor-1), 10))
	}

	newTx, err := transaction.MarshalTransaction(tx)
	if err != nil {
		log.Error().Msgf("Error marshaling tx: '%s'", err)
		return ""
	}
	return newTx
}

func (provider *XrpProvider) GetTransactionStatus(hash string) string {
	result, err := provider.client.GetTransaction(hash)
	if err != nil && strings.Contains(err.Error(), "txnNotFound") {
		return "NotFound"
	}
	if err != nil {
		return "Unconfirmed" + err.Error()
	}
	if !result.Validated {
		return "Pending"
	}

	// Transaction has been confirmed
	if result.MetaData.TransactionResult == "tesSUCCESS" {
		return "Accepted"
	}
	return "Failed"
}

func (provider *XrpProvider) GetCurrentBlockNumber() uint64 {
	lIndex, err := provider.client.GetLedgerIndex()
	if err != nil {
		log.Error().Msgf("Error getting ledger index: '%+v'", err)
		return 0
	}

	return uint64(lIndex)
}

func (provider *XrpProvider) SetCurrentBlockNumber(currentBlock uint64) {
	(*provider).currentBlock = currentBlock
}

func (provider *XrpProvider) SetNewBridgesCurrentBlockNumber(currentBlock uint64) {
	(*provider).currentNewBridgesBlock = currentBlock
}

func (provider *XrpProvider) setBridgeRequestsCurrentBlockNumber(currentBlock uint64) {
	(*provider).currentBridgeRequestsBlock = currentBlock
}

func (provider *XrpProvider) GetNewCommits(toBlock uint64) interface{} {
	log.Info().Msgf("Fetching commits from block %d to block %d", (*provider).currentBlock, toBlock)

	transactions, err := provider.GetTransactions(provider.doorAddress, int64(provider.currentBlock), int64(toBlock))
	if err != nil {
		log.Error().Msgf("Error retrieving new commits: '%s'", err)
		return nil
	}

	commits := []XrpCommit{}
	for _, tx := range transactions {
		if tx.Transaction.GetTransactionType() == "XChainCommit" {
			commits = append(commits, XrpCommit{
				BridgeId:    GetIdFromBridge(tx.Transaction.GetXChainBridge()),
				Block:       provider.currentBlock,
				ClaimId:     tx.Transaction.GetClaimId(),
				Sender:      tx.Transaction.GetAccount(),
				Amount:      tx.Transaction.GetAmount(),
				Destination: tx.Transaction.OtherChainDestination,
			})
		}
	}

	log.Debug().Msgf("Fetched %d commits in XRP", len(commits))

	return commits
}

func (provider *XrpProvider) GetNewAccountCreates(toBlock uint64) interface{} {
	log.Info().Msgf("Fetching account creates from block %d to block %d", (*provider).currentBlock, toBlock)
	transactions, err := provider.GetTransactions(provider.doorAddress, int64(provider.currentBlock), int64(toBlock))
	if err != nil {
		log.Error().Msgf("Error retrieving new account creates: '%s'", err)
		return nil
	}

	accountCreates := []XrpAccountCreate{}
	for _, tx := range transactions {
		if tx.Transaction.GetTransactionType() == "XChainAccountCreateCommit" {
			accountCreates = append(accountCreates, XrpAccountCreate{
				BridgeId:        GetIdFromBridge(tx.Transaction.GetXChainBridge()),
				Block:           provider.currentBlock,
				Sender:          tx.Transaction.GetAccount(),
				Amount:          tx.Transaction.GetAmount(),
				Destination:     *tx.Transaction.GetDestination(),
				SignatureReward: tx.Transaction.GetSignatureReward(),
			})
		}
	}

	return accountCreates
}

func (provider *XrpProvider) FetchNewBridges(toBlock uint64) error {
	log.Info().Msgf("Fetching new bridges from block %d to block %d", (*provider).currentNewBridgesBlock, toBlock)
	transactions, err := provider.GetTransactions(provider.doorAddress, int64(provider.currentNewBridgesBlock), int64(toBlock))
	if err != nil {
		log.Error().Msgf("Error fetching new bridges: '%s'", err)
		return err
	}

	for _, tx := range transactions {
		if tx.Transaction.GetTransactionType() == "XChainCreateBridge" {
			bridge := tx.Transaction.GetXChainBridge()
			_, exists := provider.bridgeProviders[GetIdFromBridge(bridge)]
			if !exists {
				bridgeProvider := CreateXrpBridgeProvider(provider.doorAddress, bridge, 0)
				provider.unpairedBridgeProviders[bridgeProvider.bridgeId] = bridgeProvider
			}
		}
	}

	return nil
}

func (provider *XrpProvider) FetchNewBridgeRequests(toBlock uint64) (interface{}, error) {
	provider.setBridgeRequestsCurrentBlockNumber(toBlock + 1)
	return nil, nil
}

func (provider *XrpProvider) RetryNewBridgeRequest(bridgeRequestCounter interface{}) error {
	return nil
}

type TransactionsCache struct {
	transactions []xrpl.TransactionAndMetadata
	fromBlock    int64
	toBlock      int64
}

func (provider *XrpProvider) GetAmmInfo(asset *xrpl.AmmAsset, asset2 *xrpl.AmmAsset) (*xrpl.AmmInfoResult, error) {
	return provider.client.GetAmmInfo(asset, asset2)
}

func (provider *XrpProvider) GetTransactions(accountId string, fromBlock, toBlock int64) ([]xrpl.TransactionAndMetadata, error) {
	key := "xrp-bridge-transactions-" + accountId
	var startIdx *int
	var endIdx *int

	value, err := cache.Get(key)
	cached, isTxs := value.(TransactionsCache)
	if err != nil || !isTxs {
		log.Debug().Msgf("Querying all transactions")
		result, err := provider.client.GetAccountTransactions(accountId, fromBlock, toBlock, MaxTxsPerRequest, nil)
		if err != nil {
			return nil, err
		}

		cached = TransactionsCache{transactions: result.Transactions, fromBlock: fromBlock, toBlock: toBlock}
		start := 0
		end := len(cached.transactions)
		startIdx = &start
		endIdx = &end
	}
	if cached.fromBlock > fromBlock {
		log.Debug().Msgf("Querying old transactions")
		result, err := provider.client.GetAccountTransactions(accountId, fromBlock, cached.fromBlock-1, MaxTxsPerRequest, nil)
		if err != nil {
			return nil, err
		}

		cached.fromBlock = fromBlock
		cached.transactions = append(cached.transactions, result.Transactions...)
		end := len(cached.transactions)
		endIdx = &end
	}
	if cached.toBlock < toBlock {
		log.Debug().Msgf("Querying new transactions")
		result, err := provider.client.GetAccountTransactions(accountId, cached.toBlock+1, toBlock, MaxTxsPerRequest, nil)
		if err != nil {
			return nil, err
		}

		cached.toBlock = toBlock
		cached.transactions = append(result.Transactions, cached.transactions...)
		start := 0
		startIdx = &start
	}

	if toBlock-OldestBlockDiff > 0 {
		// Remove transaction from older blocks
		idx := findEarliestBlockInTxs(cached.transactions, uint64(toBlock-OldestBlockDiff))
		if *idx > 0 {
			cached.transactions = cached.transactions[0:*idx]
			endIdx = nil
		}
	}

	cache.Set(key, cached, nil)

	// return only relevant transactions
	if startIdx == nil {
		startIdx = findEarliestBlockInTxs(cached.transactions, uint64(toBlock))
	}
	if endIdx == nil {
		endIdx = findEarliestBlockInTxs(cached.transactions, uint64(fromBlock)-1)
	}

	return cached.transactions[*startIdx:*endIdx], nil
}

func (provider *XrpProvider) GetUnattestedClaimById(claimId uint64, bridgeId string) (interface{}, error) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return nil, errors.New("bridge provider not found")
	}

	// TODO: decide maximum block to look for and maybe do it by parts
	var claimCreator *string = nil
	fromBlock := int64(provider.currentBlock) - MaxBlocksPerRequest
	toBlock := int64(provider.currentBlock)

	for i := 0; i < 10 && claimCreator == nil; i++ {
		transactions, err := provider.GetTransactions(provider.doorAddress, fromBlock, toBlock)
		if err != nil {
			log.Error().Msgf("Error retrieving new commits: '%s'", err)
			return nil, err
		}

	TXS:
		for _, tx := range transactions {
			if tx.Transaction.GetTransactionType() == "XChainCreateClaimID" && bridgesEqual(bridgeProvider.bridge, tx.Transaction.GetXChainBridge()) {
				for _, affectedNode := range tx.MetaData.AffectedNodes {
					if affectedNode.CreatedNode != nil && affectedNode.CreatedNode.LedgerEntryType == "XChainOwnedClaimID" {
						newFields := affectedNode.CreatedNode.NewFields
						if newFields.XChainClaimID != nil {
							claimIdParsed, err := strconv.ParseUint(*newFields.XChainClaimID, 16, 64)
							if err != nil {
								log.Error().Msgf("Error converting XChainClaimID to uint64: '%+v'", err)
							} else if claimIdParsed == claimId {
								account := tx.Transaction.GetAccount()
								claimCreator = &account
								break TXS
							}
						}
					}
				}
			}
		}

		toBlock = fromBlock
		fromBlock -= MaxBlocksPerRequest
	}

	if claimCreator == nil {
		return nil, nil
	}

	objectType := "xchain_owned_claim_id"
	ledgerIndex := "current"
	accObjects, err := provider.client.GetAccountObjects(*claimCreator, &ledgerIndex, &objectType)
	if err != nil {
		log.Error().Msgf("Error getting account objects: '%s'", err)
		return nil, err
	} else {
		for _, object := range accObjects.Objects {
			jsonObj, _ := json.Marshal(object)
			claim := xrpl.XChainClaimObject{}
			if err := json.Unmarshal(jsonObj, &claim); err == nil {
				claimIdUint, err := strconv.ParseUint(claim.XChainClaimID, 16, 64)
				if err == nil && claimIdUint == claimId && bridgesEqual(bridgeProvider.bridge, claim.XChainBridge) {
					xrpClaim := XrpClaim{claimIdUint, *claimCreator, claim.OtherChainSource}

					// Check claim has been attested
					for _, attestation := range claim.XChainClaimAttestations {
						if attestation.XChainClaimProofSig.AttestationSignerAccount == provider.witnessAddress {
							return nil, nil
						}
					}
					return xrpClaim, nil
				}
			}
		}
	}

	return nil, nil
}

func (provider *XrpProvider) CheckWitnessHasAttestedCreateAccount(account string, bridgeId string) (bool, error) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return false, errors.New("bridge provider not found")
	}

	ledgerIndex := "current"
	objectType := "xchain_owned_create_account_claim_id"
	accObjects, err := provider.client.GetAccountObjects(provider.doorAddress, &ledgerIndex, &objectType)
	if err != nil {
		log.Error().Msgf("Error getting account objects: '%s'", err)
		return false, err
	}

	for _, object := range accObjects.Objects {
		jsonObj, _ := json.Marshal(object)
		createAccount := xrpl.XChainCreateAccountObject{}
		if err := json.Unmarshal(jsonObj, &createAccount); err == nil {
			if bridgesEqual(bridgeProvider.bridge, createAccount.XChainBridge) {
				for _, attestation := range createAccount.XChainCreateAccountAttestations {
					if attestation.XChainCreateAccountProofSig.AttestationSignerAccount == provider.witnessAddress &&
						attestation.XChainCreateAccountProofSig.Destination == account {
						return true, nil
					}
				}
			}
		}
	}

	return false, nil
}

func (provider *XrpProvider) CheckAccountCreated(account string, bridgeId string) (bool, error) {
	_, err := provider.client.GetAccountInfo(account, nil)
	if err != nil && strings.Contains(err.Error(), "actNotFound") {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (provider *XrpProvider) GetChainId() *big.Int {
	return big.NewInt(int64(provider.networkId))
}

func (provider *XrpProvider) GetNonce() *uint {
	currentLI := "current"
	var seq uint
	lastSeq, err := provider.client.GetAccountInfo(provider.signerProvider.GetAddress(), &currentLI)
	if err != nil {
		log.Error().Msgf("Error getting account info: '%+v'", err)
		return nil
	}

	seq = uint(lastSeq.AccountData.Sequence)
	return &seq
}

func (provider *XrpProvider) IsInSignerList() bool {
	timeToCheck := time.Now().Add(-1 * provider.recheckSignerDuration * time.Second)
	if provider.inSignerList == nil || timeToCheck.After(provider.lastSignerCheck) {
		provider.lastSignerCheck = time.Now()

		inSignerList := false
		objectType := "signer_list"
		ledgerIndex := "current"
		accObjects, err := provider.client.GetAccountObjects(provider.doorAddress, &ledgerIndex, &objectType)
		if err != nil {
			log.Error().Msgf("Error getting account objects: '%s'", err)
		} else {
			for _, object := range accObjects.Objects {
				jsonObj, _ := json.Marshal(object)
				signerList := xrpl.SignerListObject{}
				if err := json.Unmarshal(jsonObj, &signerList); err == nil {
					for _, signerEntry := range signerList.SignerEntries {
						if signerEntry.SignerEntry.Account == provider.witnessAddress {
							inSignerList = true
						}
					}
				}
			}
		}

		provider.inSignerList = &inSignerList
	}

	return *provider.inSignerList
}

func (provider *XrpProvider) GetCurrentCreateAccountCount() uint64 {
	var v uint64

	err := cache.GetAndSet(func() any {
		// Retrieve current account create count from bridge
		ledgerIndex := "current"
		objectType := "bridge"
		bridgeObjects, err := provider.client.GetAccountObjects(provider.doorAddress, &ledgerIndex, &objectType)
		if err != nil {
			log.Error().Msgf("Error getting bridge account objects: '%+v'", err)
			return 0
		}

		for _, object := range bridgeObjects.Objects {
			jsonObj, _ := json.Marshal(object)
			bridge := xrpl.XChainBridgeObject{}
			if err = json.Unmarshal(jsonObj, &bridge); err == nil && bridge.LedgerEntryType == "Bridge" && bridge.XChainBridge.IssuingChainIssue.Currency == "XRP" {
				createCount, err := strconv.ParseUint(bridge.XChainAccountCreateCount, 16, 64)
				if err == nil {
					return createCount
				}
			}
		}

		return 0
	}, &v, time.Now().Add(time.Second*2))
	if err != nil {
		log.Error().Msgf("Error getting current create account count from cache %s", err)
	}
	return v
}

func (provider *XrpProvider) GetTransactionCreateCount(payload string) (uint64, error) {
	tx, err := transaction.UnmarshalTransaction(payload)
	if err != nil {
		return 0, err
	}

	count := tx.GetAccountCreateCount()
	if count == 0 {
		return 0, errors.New("transaction has no field XChainAccountCreateCount")
	}

	return count, nil
}

func (provider *XrpProvider) ConvertToDecimal(payload string, bridgeId string) *big.Float {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists || bridgeProvider.isToken {
		wholeN, _ := big.NewFloat(0).SetString(payload)
		return wholeN
	}

	return utils.IntToFloatPrec(payload, XrpPrec)
}

func (provider *XrpProvider) ConvertToWhole(payload *big.Float, bridgeId string) string {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists || bridgeProvider.isToken {
		return payload.String()
	}
	return utils.FloatToIntPrec(payload, XrpPrec).String()
}

func (provider *XrpProvider) GetType() config.ChainType {
	return config.Xrp
}

func (provider *XrpProvider) GetTokenCodeFromAddress(address string) (string, error) {
	return "", errors.New("error can not get token code from address")
}
