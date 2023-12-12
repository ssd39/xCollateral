package evm

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"peersyst/bridge-witness-go/internal/common/cache"
	"peersyst/bridge-witness-go/internal/common/utils"
	"peersyst/bridge-witness-go/internal/signer"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

type EvmProvider struct {
	node                       string
	witnessAddress             common.Address
	doorAddress                common.Address
	bridgeAddress              common.Address
	currentBlock               uint64
	currentNewBridgesBlock     uint64
	currentBridgeRequestsBlock uint64
	client                     *ethclient.Client
	bridgeOpts                 *bind.CallOpts
	bridgeContract             *Bridge
	contractAbi                *abi.ABI
	nonce                      *uint64
	signerProvider             signer.SignerProvider
	inSignerList               *bool
	lastSignerCheck            time.Time
	recheckSignerDuration      time.Duration
	maxGas                     *big.Int
	maxGasFactor               *big.Int
	bridgeProviders            map[string]*EvmBridgeProvider
	unpairedBridgeProviders    map[string]*EvmBridgeProvider
	bridgeProvidersByKey       map[[32]byte]*EvmBridgeProvider
	createBridgeEventsArr      []*BridgeRequestCounter
	safeSigner                 SafeSigner
}

type EvmCommit struct {
	Block       uint64
	ClaimId     uint64
	Sender      string
	Amount      string
	Destination *string
	BridgeId    string
}

type EvmAccountCreate struct {
	Block           uint64
	Sender          string
	Amount          string
	Destination     string
	SignatureReward string
	BridgeId        string
}

type EvmClaim struct {
	ClaimId uint64
	Sender  string
	Source  string
}

type EvmBridge struct {
	IssuingChainDoor     string
	IssuingChainIssuer   string
	IssuingChainIssue    string
	LockingChainDoor     string
	LockingChainIssuer   string
	LockingChainIssue    string
	TokenContractAddress string
}

type BridgeRequestCounter struct {
	Tries uint
	Event *BridgeCreateBridgeRequest
}

var zeroAddress common.Address = common.HexToAddress("0x0000000000000000000000000000000000000000")
var maxAttestedIterations = 5

func Create(signerProvider signer.SignerProvider, node, doorAddress string, startingBlock uint64, signerListSeconds, maxGasFactor int64) (*EvmProvider, error) {
	client, err := ethclient.Dial(node)
	if err != nil {
		return nil, err
	}
	witnessAddress := common.HexToAddress(signerProvider.GetAddress())
	callOpts := bind.CallOpts{Pending: true, From: witnessAddress, Context: context.Background()}

	// If starting block is not specified, fetch the latest validated from the network
	currentBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Error().Msgf("Error getting block number: '%+v'", err)
		return nil, err
	}
	if startingBlock == 0 {
		startingBlock = currentBlock
	}

	nonce, err := client.PendingNonceAt(context.Background(), witnessAddress)
	if err != nil {
		log.Error().Msgf("Error getting pending nonce: '%+v'", err)
		return nil, err
	}

	safe, err := NewSafe(common.HexToAddress(doorAddress), client)
	if err != nil {
		log.Error().Msgf("Error instantiating safe contract: '%s'", err)
		return nil, err
	}

	moduleAddresses := []common.Address{}
	startAddress := common.HexToAddress("0x1") // 0x1 address is to indicate start
	pageSize := big.NewInt(20)
	result, err := safe.GetModulesPaginated(&callOpts, startAddress, pageSize)
	if err != nil {
		log.Error().Msgf("Error getting modules from safe contract: '%s'", err)
		return nil, err
	}
	for result.Next != startAddress {
		moduleAddresses = append(moduleAddresses, result.Array...)
		result, err = safe.GetModulesPaginated(&callOpts, startAddress, pageSize)
		if err != nil {
			log.Error().Msgf("Error getting modules from safe contract: '%s'", err)
			return nil, err
		}
	}

	// Should only be 1 bridge as there is only 1 contract now
	moduleAddresses = append(moduleAddresses, result.Array...)
	if len(moduleAddresses) != 1 {
		return nil, errors.New("Error should only be 1 module in safe")
	}

	bridgeContract, err := NewBridge(moduleAddresses[0], client)
	if err != nil {
		log.Error().Msgf("Error instantiating contract: '%s'", err)
		return nil, err
	}

	// Find create bridge events in contract
	// Temporary start/end until get bridges in smart contract
	currentPage := big.NewInt(0)
	bridgeProviders := map[string]*EvmBridgeProvider{}
	bridgeProvidersByKey := map[[32]byte]*EvmBridgeProvider{}

BRIDGES:
	for {
		bridges, err := bridgeContract.GetBridgesPaginated(&callOpts, currentPage)
		if err != nil {
			log.Error().Msgf("Error retrieving bridges paginated: '%s'", err)
			return nil, err
		}

		i := 0
		for i < len(bridges.Configs) && i < len(bridges.Params) {
			if bridges.Configs[i].LockingChainDoor.String() == zeroAddress.String() {
				break BRIDGES
			}

			bridgeProvider := CreateEvmBridgeProvider(client, &callOpts, bridgeContract, &bridges.Configs[i], &bridges.Params[i])
			if bridgeProvider == nil {
				return nil, errors.New("Error creating bridge provider")
			}
			bridgeProviders[bridgeProvider.bridgeId] = bridgeProvider
			bridgeProvidersByKey[bridgeProvider.bridgeKey] = bridgeProvider

			i += 1
		}

		currentPage = currentPage.Add(currentPage, big.NewInt(1))
	}

	provider := EvmProvider{
		node,
		witnessAddress,
		common.HexToAddress(doorAddress),
		moduleAddresses[0],
		startingBlock,
		currentBlock,
		currentBlock,
		client,
		&callOpts,
		bridgeContract,
		nil,
		&nonce,
		signerProvider,
		nil,
		time.Now(),
		time.Duration(signerListSeconds),
		big.NewInt(5000000),
		big.NewInt(maxGasFactor),
		map[string]*EvmBridgeProvider{},
		bridgeProviders,
		bridgeProvidersByKey,
		[]*BridgeRequestCounter{},
		SafeSigner{*safe, signerProvider},
	}

	return &provider, err
}

func (provider *EvmProvider) UpdateOracleData(amount, amount2 int64) error {
	priceOracleAbi := "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"doorAccount_\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"currency_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"currency2_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"amount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"amount2\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currency\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currency2\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount2_\",\"type\":\"uint256\"}],\"name\":\"updateData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	parsed, err := abi.JSON(strings.NewReader(priceOracleAbi))
	if err != nil {
		log.Error().Msgf("Error bridge ABI : '%s'", err)
		return err
	}
	input, err := parsed.Pack("updateData", big.NewInt(amount), big.NewInt(amount2))
	if err != nil {
		log.Error().Msgf("Error packing parameters : '%s'", err)
		return err
	}
	nonce := provider.getNextNonce()
	gasPrice := provider.getGasPrice()

	priceOracleContract := common.HexToAddress("0x133EEf561F068511bFc2740A1Fd33192E246637E")

	tx := types.NewTransaction(nonce, priceOracleContract, big.NewInt(0), 10000000, gasPrice, input)
	signedTx := provider.SignTransaction(encodeTransaction(tx))
	txHash, err := provider.BroadcastTransaction(signedTx)
	if err != nil {
		log.Error().Msgf("Error while updating oracle data  %v", err)
		return err
	}
	log.Info().Msgf("Transaction submitted to update oracle %s", txHash)
	return nil
}

func (provider *EvmProvider) BroadcastTransaction(payload string) (string, error) {
	rawTxBytes, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding transaction: '%s'", err)
		return "", fmt.Errorf("decoding error: %+v", err)
	}

	tx := new(types.Transaction)
	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		log.Error().Msgf("Error rlp decoding transaction: '%s'", err)
		return "", fmt.Errorf("decoding error: %+v", err)
	}

	err2 := provider.client.SendTransaction(context.Background(), tx)
	if err2 != nil {
		if strings.Contains(err2.Error(), "no result in JSON-RPC response") {
			return "", fmt.Errorf("no response error: %+v", err2)
		}
		if strings.Contains(err2.Error(), "invalid nonce") {
			return "", fmt.Errorf("invalid nonce")
		}
		return "", fmt.Errorf("unknown error: %+v", err2)
	}

	return tx.Hash().Hex(), nil
}

func (provider *EvmProvider) getNextNonce() uint64 {
	nonce := atomic.AddUint64(provider.nonce, 1)
	return nonce - 1
}

func (provider *EvmProvider) getAbi() *abi.ABI {
	if provider.contractAbi == nil {
		parsed, err := abi.JSON(strings.NewReader(BridgeMetaData.ABI))
		if err != nil {
			log.Error().Msgf("Error bridge ABI : '%s'", err)
			return nil
		}
		provider.contractAbi = &parsed
	}

	return provider.contractAbi
}

func encodeTransaction(tx *types.Transaction) string {
	rawTxBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		log.Error().Msgf("Error encoding transaction to bytes : '%s'", err)
		return ""
	}

	return hex.EncodeToString(rawTxBytes)
}

func decodeTransaction(payload string) *types.Transaction {
	hexPayload, err := hex.DecodeString(payload)
	if err != nil {
		log.Error().Msgf("Error decoding transaction string to hex bytes : '%s'", err)
		return nil
	}

	var tx types.Transaction
	err = rlp.DecodeBytes(hexPayload, &tx)
	if err != nil {
		log.Error().Msgf("Error decoding transaction : '%s'", err)
		return nil
	}

	return &tx
}

func (provider *EvmProvider) packContractParams(method string, params ...interface{}) []byte {
	abi := provider.getAbi()
	if abi == nil {
		return nil
	}

	input, err := abi.Pack(method, params...)
	if err != nil {
		log.Error().Msgf("Error packing parameters : '%s'", err)
		return nil
	}

	return input
}

func (provider *EvmProvider) unpackContractParams(method string, data []byte) []interface{} {
	abi := provider.getAbi()
	if abi == nil {
		return nil
	}

	params, err := abi.Methods[method].Inputs.Unpack(data)
	if err != nil {
		log.Error().Msgf("Error packing parameters : '%s'", err)
		return nil
	}

	return params
}

func (provider *EvmProvider) getGasPrice() *big.Int {
	var v big.Int
	err := cache.GetAndSet(func() any {
		gas, err := provider.client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Error().Msgf("Error fetching gas price : %s", err)
			return big.NewInt(7)
		}
		return *gas
	}, &v, time.Now().Add(time.Minute*60))
	if err != nil {
		log.Error().Msgf("Error getting gas price from cache %s", err)
	}
	return &v
}

func (provider *EvmProvider) GetAttestClaimTransaction(claimId uint64, sender string, amount string, destination string, bridgeId string) (string, uint64) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return "", 0
	}

	senderAdd := common.HexToAddress(sender)
	destinationAdd := common.HexToAddress(destination)
	amountBI, _ := big.NewInt(0).SetString(amount, 10)
	input := provider.packContractParams("addClaimAttestation", *bridgeProvider.bridge, big.NewInt(int64(claimId)), amountBI, senderAdd, destinationAdd)
	if input == nil {
		return "", 0
	}

	nonce := provider.getNextNonce()
	gasPrice := provider.getGasPrice()

	tx := types.NewTransaction(nonce, provider.bridgeAddress, nil, 10000000, gasPrice, input)
	return encodeTransaction(tx), nonce
}

func (provider *EvmProvider) GetAttestAccountCreateTransaction(sender, amount, destination, signatureReward string, bridgeId string) (string, uint64) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return "", 0
	}

	amountBI, _ := big.NewInt(0).SetString(amount, 10)
	sigRewardBI, _ := big.NewInt(0).SetString(signatureReward, 10)
	input := provider.packContractParams("addCreateAccountAttestation", *bridgeProvider.bridge, common.HexToAddress(destination), amountBI, sigRewardBI)
	if input == nil {
		return "", 0
	}

	nonce := provider.getNextNonce()
	gasPrice := provider.getGasPrice()

	tx := types.NewTransaction(nonce, provider.bridgeAddress, nil, 10000000, gasPrice, input)
	return encodeTransaction(tx), nonce
}

func (provider *EvmProvider) SignTransaction(transaction string) string {
	chainId := provider.GetChainId()
	if chainId == nil {
		return ""
	}

	opts := &signer.SignEvmTransactionOpts{ChainId: chainId}
	return provider.signerProvider.SignTransaction(transaction, opts)
}

func (provider *EvmProvider) SignEncodedCreateBridgeTransaction(encodedTx string, isLocking bool, minBridgeReward, maxBridgeReward uint64, otherChainAddress, tokenAddress, tokenCode string) (string, string, string, error) {
	tx, err := parseSafeTransactionInJSON(encodedTx)
	if err != nil {
		return "", "", "", err
	}

	if tx.To.String() != provider.bridgeAddress.String() {
		return "", "", "", errors.New("error create bridge transaction to (" + tx.To.String() + ") should be multitoken contract address: " + provider.bridgeAddress.String())
	}
	if tx.Value.Cmp(big.NewInt(0)) != 0 {
		return "", "", "", errors.New("error create bridge transaction value (" + tx.Value.String() + ") should be 0")
	}
	if tx.Operation != 0 {
		return "", "", "", errors.New("error create bridge transaction operation (" + string(tx.Operation) + ") should be 0")
	}
	if tx.SafeTxGas.Cmp(provider.maxGas) == 1 {
		return "", "", "", errors.New("error create bridge transaction safeTxGas (" + tx.SafeTxGas.String() + ") should be less than or equal than maxGas: " + provider.maxGas.String())
	}
	if tx.BaseGas.Cmp(provider.maxGas) == 1 {
		return "", "", "", errors.New("error create bridge transaction baseGas (" + tx.BaseGas.String() + ") should be less than or equal than maxGas: " + provider.maxGas.String())
	}
	if tx.GasPrice.Cmp(provider.maxGasFactor) == 1 {
		return "", "", "", errors.New("error create bridge transaction gasPrice (" + tx.GasPrice.String() + ") should be less than or equal than maxGasFactor: " + provider.maxGasFactor.String())
	}
	if tx.GasToken.String() != zeroAddress.String() {
		return "", "", "", errors.New("error create bridge transaction gasToken (" + tx.GasToken.String() + ") should be 0 address: " + zeroAddress.String())
	}
	if tx.RefundReceiver.String() != zeroAddress.String() {
		return "", "", "", errors.New("error create bridge transaction refundReceiver (" + tx.RefundReceiver.String() + ") should be 0 address: " + zeroAddress.String())
	}

	params := provider.unpackContractParams("createBridge", tx.Data[4:])
	if params == nil {
		return "", "", "", errors.New("error unpacking create bridge transaction params")
	}

	bridgeCfgJson, err := json.Marshal(params[0])
	if err != nil {
		return "", "", "", err
	}
	var bridgeCfg XChainTypesBridgeConfig
	err = json.Unmarshal(bridgeCfgJson, &bridgeCfg)
	if err != nil {
		return "", "", "", err
	}

	if isLocking {
		if !strings.EqualFold(bridgeCfg.LockingChainDoor.String(), provider.doorAddress.String()) {
			return "", "", "", errors.New("error create bridge transaction lockingChainDoor (" + bridgeCfg.LockingChainDoor.String() + ") should be: " + provider.doorAddress.String())
		}
		if !strings.EqualFold(bridgeCfg.IssuingChainDoor.String(), otherChainAddress) {
			return "", "", "", errors.New("error create bridge transaction issuingChainDoor (" + bridgeCfg.IssuingChainDoor.String() + ") should be: " + otherChainAddress)
		}
		if !strings.EqualFold(bridgeCfg.IssuingChainIssue.Issuer.String(), otherChainAddress) {
			return "", "", "", errors.New("error create bridge transaction issuingChainIssue Issuer (" + bridgeCfg.IssuingChainIssue.Issuer.String() + ") should be: " + otherChainAddress)
		}
	} else {
		if !strings.EqualFold(bridgeCfg.LockingChainDoor.String(), otherChainAddress) {
			return "", "", "", errors.New("error create bridge transaction lockingChainDoor (" + bridgeCfg.LockingChainDoor.String() + ") should be: " + otherChainAddress)
		}
		if !strings.EqualFold(bridgeCfg.IssuingChainDoor.String(), provider.doorAddress.String()) {
			return "", "", "", errors.New("error create bridge transaction issuingChainDoor (" + bridgeCfg.IssuingChainDoor.String() + ") should be: " + provider.doorAddress.String())
		}
		if !strings.EqualFold(bridgeCfg.IssuingChainIssue.Issuer.String(), provider.doorAddress.String()) {
			return "", "", "", errors.New("error create bridge transaction issuingChainIssue Issuer (" + bridgeCfg.IssuingChainIssue.Issuer.String() + ") should be: " + provider.doorAddress.String())
		}
	}
	if !strings.EqualFold(bridgeCfg.LockingChainIssue.Issuer.String(), tokenAddress) {
		return "", "", "", errors.New("error create bridge transaction lockingChainIssue Issuer (" + bridgeCfg.LockingChainIssue.Issuer.String() + ") should be: " + tokenAddress)
	}
	if bridgeCfg.LockingChainIssue.Currency != tokenCode {
		return "", "", "", errors.New("error create bridge transaction lockingChainIssue Currency (" + bridgeCfg.LockingChainIssue.Currency + ") should be: " + tokenCode)
	}
	issuingStrings := strings.Split(bridgeCfg.IssuingChainIssue.Currency, "-")
	if len(issuingStrings) != 2 {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency (" + bridgeCfg.IssuingChainIssue.Currency + ") should be splitted to 2 by -")
	}
	if issuingStrings[0] != tokenCode {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency first part (" + issuingStrings[0] + ") should be: " + tokenCode)
	}
	_, err = hex.DecodeString(issuingStrings[1])
	if err != nil || len(issuingStrings[1]) != 8 {
		return "", "", "", errors.New("error create bridge transaction issuingChainIssue Currency second part (" + issuingStrings[1] + ") should have length 8")
	}

	bridgeParamsJson, err := json.Marshal(params[1])
	if err != nil {
		return "", "", "", err
	}
	var bridgeParams XChainTypesBridgeParams
	err = json.Unmarshal(bridgeParamsJson, &bridgeParams)
	if err != nil {
		return "", "", "", err
	}

	minBridgeRewardWhole := utils.FloatToIntPrec(big.NewFloat(float64(minBridgeReward)), 18)
	if bridgeParams.SignatureReward.Cmp(minBridgeRewardWhole) == -1 {
		return "", "", "", errors.New("error create bridge transaction signature reward (" + bridgeParams.SignatureReward.String() + ") should be more than or equal than minBridgeReward: " + string(rune(minBridgeReward)))
	}
	maxBridgeRewardWhole := utils.FloatToIntPrec(big.NewFloat(float64(maxBridgeReward)), 18)
	if bridgeParams.SignatureReward.Cmp(maxBridgeRewardWhole) == 1 {
		return "", "", "", errors.New("error create bridge transaction signature reward (" + bridgeParams.SignatureReward.String() + ") should be less than or equal than minBridgeReward: " + string(rune(maxBridgeReward)))
	}

	txHash, err := provider.safeSigner.GetTransactionHash(*tx)
	if err != nil {
		return "", "", "", err
	}

	signature, err := provider.safeSigner.SignTransactionHash(hex.EncodeToString(txHash[:]))
	if err != nil {
		return "", "", "", err
	}

	return signature.signer, signature.publicKey, signature.data, nil
}

func (provider *EvmProvider) GetNoOpTransaction(nonce uint, feeFactor uint) string {
	gasPrice := big.NewInt(int64(feeFactor))
	gasPrice.Mul(gasPrice, provider.getGasPrice())
	tx := types.NewTransaction(uint64(nonce), zeroAddress, nil, 21000, gasPrice, nil)
	return encodeTransaction(tx)
}

func (provider *EvmProvider) SetTransactionGasPrice(payload string, feeFactor uint) string {
	tx := decodeTransaction(payload)
	if tx == nil {
		return ""
	}
	newTx := types.NewTransaction(tx.Nonce(), *tx.To(), tx.Value(), tx.Gas(), big.NewInt(tx.GasPrice().Int64()*int64(feeFactor)), tx.Data())
	return encodeTransaction(newTx)
}

func (provider *EvmProvider) GetTransactionStatus(hash string) string {
	_, pending, err := provider.client.TransactionByHash(context.Background(), common.HexToHash(hash))
	if err == ethereum.NotFound {
		return "NotFound"
	}
	if pending {
		return "Pending"
	}
	if err != nil {
		return "Unconfirmed" + err.Error()
	}

	// Confirmed transaction
	receipt, err := provider.client.TransactionReceipt(context.Background(), common.HexToHash(hash))
	if err != nil {
		return "Confirmed" + err.Error()
	}
	if types.ReceiptStatusFailed == receipt.Status {
		return "Failed"
	}
	if types.ReceiptStatusSuccessful == receipt.Status {
		return "Accepted"
	}

	// Should never get here as receipt.Status is 1 or 0
	return "Error"
}

func (provider *EvmProvider) GetCurrentBlockNumber() uint64 {
	block, err := (*provider).client.BlockNumber(context.Background())
	if err != nil {
		log.Error().Msgf("Error getting block number : '%s'", err)
		return 0
	}

	return block
}

func (provider *EvmProvider) SetCurrentBlockNumber(currentBlock uint64) {
	endBlock := getEndBlock((*provider).currentBlock, currentBlock)
	(*provider).currentBlock = endBlock
}

func (provider *EvmProvider) SetNewBridgesCurrentBlockNumber(currentBlock uint64) {
	endBlock := getEndBlock((*provider).currentNewBridgesBlock, currentBlock)
	(*provider).currentNewBridgesBlock = endBlock
}

func (provider *EvmProvider) setBridgeRequestsCurrentBlockNumber(currentBlock uint64) {
	endBlock := getEndBlock((*provider).currentBridgeRequestsBlock, currentBlock)
	(*provider).currentBridgeRequestsBlock = endBlock
}

func (provider *EvmProvider) SetBridgeValidated(bridgeId string) interface{} {
	bridgeProvider, exists := provider.unpairedBridgeProviders[bridgeId]
	if !exists {
		return nil
	}

	delete(provider.unpairedBridgeProviders, bridgeId)
	provider.bridgeProviders[bridgeId] = bridgeProvider

	return nil
}

func (provider *EvmProvider) GetUnpairedBridges() interface{} {
	return provider.unpairedBridgeProviders
}

func getEndBlock(fromBlock, toBlock uint64) uint64 {
	end := toBlock
	if fromBlock+10000 < toBlock {
		end = fromBlock + 10000
	}
	return end
}

func (provider *EvmProvider) GetNewCommits(toBlock uint64) interface{} {
	commits := []EvmCommit{}
	fromBlock := provider.currentBlock
	endBlock := getEndBlock(fromBlock, toBlock)
	filterOpts := bind.FilterOpts{Start: fromBlock, End: &endBlock, Context: context.Background()}
	log.Info().Msgf("Fetching commits from block %d to block %d", fromBlock, endBlock)

	commitIterator, err := provider.bridgeContract.BridgeFilterer.FilterCommit(&filterOpts, [][32]byte{}, []*big.Int{}, []common.Address{})
	if err != nil {
		log.Error().Msgf("Error filtering commits for bridge: '%s'", err)
		return commits
	}

	for commitIterator.Next() {
		bridgeProvider, exists := provider.bridgeProvidersByKey[commitIterator.Event.BridgeKey]
		if !exists {
			continue
		}

		destination := commitIterator.Event.Receiver.String()
		commits = append(commits, EvmCommit{
			Block:       commitIterator.Event.Raw.BlockNumber,
			ClaimId:     commitIterator.Event.ClaimId.Uint64(),
			Sender:      commitIterator.Event.Sender.String(),
			Amount:      commitIterator.Event.Value.Text(10),
			Destination: &destination,
			BridgeId:    bridgeProvider.bridgeId,
		})
	}

	commitWOAddressIterator, err := provider.bridgeContract.BridgeFilterer.FilterCommitWithoutAddress(&filterOpts, [][32]byte{}, []*big.Int{}, []common.Address{})
	if err != nil {
		log.Error().Msgf("Error filtering commits without address for bridge: '%s'", err)
		return commits
	}

	for commitWOAddressIterator.Next() {
		bridgeProvider, exists := provider.bridgeProvidersByKey[commitWOAddressIterator.Event.BridgeKey]
		if !exists {
			continue
		}

		commits = append(commits, EvmCommit{
			Block:       commitWOAddressIterator.Event.Raw.BlockNumber,
			ClaimId:     commitWOAddressIterator.Event.ClaimId.Uint64(),
			Sender:      commitWOAddressIterator.Event.Sender.String(),
			Amount:      commitWOAddressIterator.Event.Value.Text(10),
			Destination: nil,
			BridgeId:    bridgeProvider.bridgeId,
		})
	}

	log.Debug().Msgf("Fetched %d commits in EVM", len(commits))

	return commits
}

func (provider *EvmProvider) GetNewAccountCreates(toBlock uint64) interface{} {
	accountCreates := []EvmAccountCreate{}
	fromBlock := provider.currentBlock
	endBlock := getEndBlock(fromBlock, toBlock)
	filterOpts := bind.FilterOpts{Start: fromBlock, End: &endBlock, Context: context.Background()}
	log.Info().Msgf("Fetching account creates from block %d to block %d", fromBlock, endBlock)

	accountCreateIterator, err := provider.bridgeContract.BridgeFilterer.FilterCreateAccountCommit(&filterOpts, [][32]byte{}, []common.Address{}, []common.Address{})
	if err != nil {
		log.Error().Msgf("Error filtering account creates for bridge: '%s'", err)
		return accountCreates
	}

	for accountCreateIterator.Next() {
		bridgeProvider, exists := provider.bridgeProvidersByKey[accountCreateIterator.Event.BridgeKey]
		if !exists {
			continue
		}

		accountCreates = append(accountCreates, EvmAccountCreate{
			Block:           accountCreateIterator.Event.Raw.BlockNumber,
			Sender:          accountCreateIterator.Event.Creator.String(),
			Amount:          accountCreateIterator.Event.Value.Text(10),
			Destination:     accountCreateIterator.Event.Destination.String(),
			SignatureReward: accountCreateIterator.Event.SignatureReward.Text(10),
			BridgeId:        bridgeProvider.bridgeId,
		})
	}

	return accountCreates
}

func (provider *EvmProvider) FetchNewBridges(toBlock uint64) error {
	fromBlock := provider.currentNewBridgesBlock
	endBlock := getEndBlock(fromBlock, toBlock)
	filterOpts := bind.FilterOpts{Start: fromBlock, End: &endBlock, Context: context.Background()}
	log.Info().Msgf("Fetching new bridges from block %d to block %d", fromBlock, endBlock)

	createBridgeIterator, err := provider.bridgeContract.BridgeFilterer.FilterCreateBridge(&filterOpts, [][32]byte{})
	if err != nil {
		log.Error().Msgf("Error retrieving create bridge events: '%s'", err)
		return err
	}

	for createBridgeIterator.Next() {
		_, exists := provider.bridgeProvidersByKey[createBridgeIterator.Event.BridgeKey]
		if exists {
			continue
		}

		bridgeProvider := CreateEvmBridgeProviderFromEvent(provider.client, provider.bridgeOpts, provider.bridgeContract, createBridgeIterator.Event)
		if bridgeProvider == nil {
			return errors.New("Error creating bridge provider")
		}
		provider.unpairedBridgeProviders[bridgeProvider.bridgeId] = bridgeProvider
		provider.bridgeProvidersByKey[bridgeProvider.bridgeKey] = bridgeProvider
	}

	return nil
}

func (provider *EvmProvider) FetchNewBridgeRequests(toBlock uint64) (interface{}, error) {
	fromBlock := provider.currentBridgeRequestsBlock
	endBlock := getEndBlock(fromBlock, toBlock)
	filterOpts := bind.FilterOpts{Start: fromBlock, End: &endBlock, Context: context.Background()}
	log.Info().Msgf("Fetching new bridge requests from block %d to block %d", fromBlock, endBlock)

	createBridgeRequestIterator, err := provider.bridgeContract.BridgeFilterer.FilterCreateBridgeRequest(&filterOpts)
	if err != nil {
		log.Error().Msgf("Error retrieving create bridge request events: '%s'", err)
		return nil, err
	}

	events := provider.createBridgeEventsArr
	provider.createBridgeEventsArr = []*BridgeRequestCounter{}
	for createBridgeRequestIterator.Next() {
		events = append(events, &BridgeRequestCounter{
			Tries: 1,
			Event: createBridgeRequestIterator.Event,
		})
	}

	provider.setBridgeRequestsCurrentBlockNumber(toBlock + 1)

	return events, nil
}

func (provider *EvmProvider) RetryNewBridgeRequest(bridgeRequestCounter interface{}) error {
	brc, ok := bridgeRequestCounter.(*BridgeRequestCounter)
	if !ok {
		return errors.New("Error invalid type of bridge request counter")
	}

	provider.createBridgeEventsArr = append(provider.createBridgeEventsArr, &BridgeRequestCounter{
		Tries: brc.Tries + 1,
		Event: brc.Event,
	})

	return nil
}

func (provider *EvmProvider) GetUnattestedClaimById(claimId uint64, bridgeId string) (interface{}, error) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return nil, errors.New("Error finding bridge provider")
	}

	hasAttested, err := provider.checkWitnessHasAttestedClaim(claimId, bridgeProvider.bridgeKey)
	if err != nil {
		return nil, err
	}
	if hasAttested {
		return nil, nil
	}

	creator, sender, exists, err := provider.bridgeContract.GetBridgeClaim(provider.bridgeOpts, *bridgeProvider.bridge, big.NewInt(int64(claimId)))
	if err != nil {
		log.Error().Msgf("Error fetching claim by id: '%+v'", err)
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	return EvmClaim{
		ClaimId: claimId,
		Sender:  creator.String(),
		Source:  sender.String(),
	}, nil
}

func (provider *EvmProvider) checkWitnessHasAttestedClaim(claimId uint64, bridgeKey [32]byte) (bool, error) {
	i := 0
	endBlock := provider.GetCurrentBlockNumber()
	start := endBlock - 10000
	if endBlock < 10000 {
		start = 0
	}

	for i < maxAttestedIterations {
		filterOpts := bind.FilterOpts{Start: start, End: &endBlock, Context: context.Background()}
		claimIdBI := big.NewInt(int64(claimId))

		claimIterator, err := provider.bridgeContract.BridgeFilterer.FilterAddClaimAttestation(&filterOpts, [][32]byte{bridgeKey}, []*big.Int{claimIdBI}, []common.Address{provider.witnessAddress})
		if err != nil {
			log.Error().Msgf("Error finding claim attestation event: '%s'", err)
			return false, err
		}

		for claimIterator.Next() {
			if (*claimIterator.Event.ClaimId).String() == (*claimIdBI).String() && claimIterator.Event.Witness == provider.witnessAddress {
				return true, nil
			}
		}

		if start > 10000 {
			start = start - 10000
			endBlock = endBlock - 10000
		} else if start > 0 {
			endBlock = start
			start = 0
		} else {
			// Start is 0, not found, break
			break
		}
		i += 1
	}

	return false, nil
}

func (provider *EvmProvider) CheckWitnessHasAttestedCreateAccount(destination string, bridgeId string) (bool, error) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return false, errors.New("Error finding bridge provider")
	}

	i := 0
	endBlock := provider.GetCurrentBlockNumber()
	start := endBlock - 10000
	if endBlock < 10000 {
		start = 0
	}

	for i < maxAttestedIterations {
		filterOpts := bind.FilterOpts{Start: start, End: &endBlock, Context: context.Background()}
		receiver := common.HexToAddress(destination)

		createAccountIterator, err := provider.bridgeContract.BridgeFilterer.FilterAddCreateAccountAttestation(&filterOpts, [][32]byte{bridgeProvider.bridgeKey}, []common.Address{provider.witnessAddress}, []common.Address{receiver})
		if err != nil {
			log.Error().Msgf("Error finding claim attestation event: '%s'", err)
			return false, err
		}

		for createAccountIterator.Next() {
			if createAccountIterator.Event.Receiver == receiver && createAccountIterator.Event.Witness == provider.witnessAddress {
				return true, nil
			}
		}

		if start > 10000 {
			start = start - 10000
			endBlock = endBlock - 10000
		} else if start > 0 {
			endBlock = start
			start = 0
		} else {
			// Start is 0, not found, break
			break
		}
		i += 1
	}

	return false, nil
}

func (provider *EvmProvider) CheckAccountCreated(destination string, bridgeId string) (bool, error) {
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if !exists {
		return false, errors.New("Error finding bridge provider")
	}

	_, isCreated, _, err := provider.bridgeContract.GetBridgeCreateAccount(provider.bridgeOpts, *bridgeProvider.bridge, common.HexToAddress(destination))
	if err != nil {
		return false, err
	}
	if !isCreated {
		return false, nil
	}

	return true, nil
}

func (provider *EvmProvider) GetChainId() *big.Int {
	var v big.Int
	err := cache.GetAndSet(func() any {
		chainId, err := provider.client.ChainID(context.Background())
		if err != nil {
			log.Error().Msgf("Error fetching chainId: '%s'", err)
			return nil
		}
		log.Info().Msgf("Fetched chain id %v", &chainId)
		return *chainId
	}, &v, time.Now().Add(time.Hour*24))
	if err != nil {
		log.Error().Msgf("Error getting chain id from cache %s", err)
	}
	return &v
}

func (provider *EvmProvider) GetNonce() *uint {
	var v uint
	err := cache.GetAndSet(func() any {
		log.Debug().Msgf("Fetching nonce")
		nonce, err := provider.client.PendingNonceAt(context.Background(), provider.witnessAddress)
		if err != nil {
			log.Error().Msgf("Error getting pending nonce : '%s'", err)
			return nil
		}
		return uint(nonce)
	}, &v, time.Now().Add(time.Second*5))
	if err != nil {
		log.Error().Msgf("Error getting nonce from cache %s", err)
	}
	return &v
}

func (provider *EvmProvider) GetAmmInfo(asset *xrpl.AmmAsset, asset2 *xrpl.AmmAsset) (*xrpl.AmmInfoResult, error) {
	return &xrpl.AmmInfoResult{}, nil
}

func (provider *EvmProvider) IsInSignerList() bool {
	timeToCheck := time.Now().Add(-1 * provider.recheckSignerDuration * time.Second)
	if provider.inSignerList == nil || timeToCheck.After(provider.lastSignerCheck) {
		provider.lastSignerCheck = time.Now()
		inSignerList := false

		witnesses, err := provider.bridgeContract.GetWitnesses(provider.bridgeOpts)
		if err == nil {
			for _, witness := range witnesses {
				if witness == provider.witnessAddress {
					inSignerList = true
				}
			}
		}

		provider.inSignerList = &inSignerList
	}

	return *provider.inSignerList
}

func (provider *EvmProvider) GetCurrentCreateAccountCount() uint64 {
	return 0
}

func (provider *EvmProvider) GetTransactionCreateCount(payload string) (uint64, error) {
	return 0, nil
}

func (provider *EvmProvider) ConvertToDecimal(payload string, bridgeId string) *big.Float {
	decimals := 18
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if exists {
		decimals = bridgeProvider.assetDecimals
	}
	return utils.IntToFloatPrec(payload, decimals)
}

func (provider *EvmProvider) ConvertToWhole(payload *big.Float, bridgeId string) string {
	decimals := 18
	bridgeProvider, exists := provider.bridgeProviders[bridgeId]
	if exists {
		decimals = bridgeProvider.assetDecimals
	}
	return utils.FloatToIntPrec(payload, decimals).String()
}

func (provider *EvmProvider) GetType() config.ChainType {
	return config.Evm
}

func (provider *EvmProvider) GetTokenCodeFromAddress(address string) (string, error) {
	instance, err := NewToken(common.HexToAddress(address), provider.client)
	if err != nil {
		return "", err
	}

	code, err := instance.Symbol(provider.bridgeOpts)
	if err != nil {
		return "", err
	}

	return code, nil
}
