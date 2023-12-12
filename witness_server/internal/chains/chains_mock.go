package chains

import (
	"fmt"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/chains/evm"
	"peersyst/bridge-witness-go/internal/chains/xrp"
	"peersyst/bridge-witness-go/internal/common/utils"
	"strings"
)

type TestProvider struct {
	BlockNumber                              uint64
	NewBridgesBlockNumber                    uint64
	BridgeRequestsBlockNumber                uint64
	isInSignerList                           bool
	chainId                                  *big.Int
	Nonce                                    *uint
	AccountCount                             uint64
	GetCurrentBlockCalledTimes               uint64
	SetCurrentBlockCalledTimes               uint64
	SetNewBridgesCurrentBlockCalledTimes     uint64
	SetBridgeRequestsCurrentBlockCalledTimes uint64
	GetNewCommitsCalledTimes                 uint64
	GetNewAccountCreatesCalledTimes          uint64
	CheckWitnessAttestedCalledTimes          uint64
	GetUnattestedClaimCalledTimes            uint64
	GetAttestClaimTxCalledTimes              uint64
	GetAttestCreateAccountTxCalledTimes      uint64
	CheckAccountCreatedCalledTimes           uint64
	GetNoOpTransactionCalledTimes            uint64
	SetTransactionGasPriceCalledTimes        uint64
	ChainType                                config.ChainType
}

func (provider *TestProvider) BroadcastTransaction(payload string) (string, error) {
	if strings.Contains(payload, "timeout") {
		return "", fmt.Errorf(NoResponseError + " error")
	}
	if strings.Contains(payload, "success") {
		return "hash", nil
	}
	if strings.Contains(payload, "invalid nonce") {
		return "", fmt.Errorf(InvalidNonce)
	}
	if strings.Contains(payload, "ignorable") {
		return "", fmt.Errorf(IgnorableError + " error")
	}
	if strings.Contains(payload, "decoding") {
		return "", fmt.Errorf(DecodingError + " error")
	}
	if strings.Contains(payload, "unknown") {
		return "", fmt.Errorf(UnknownError + " error")
	}
	return "", fmt.Errorf(NoResponseError + " error")
}

func (provider *TestProvider) GetAttestClaimTransaction(claimId uint64, sender, amount, destination, bridgeId string) (string, uint64) {
	provider.GetAttestClaimTxCalledTimes += 1
	if destination == "0x177adf17f5ac5df0178a24ba5b805a88a7a4be2a" || destination == "rs99jCuSAjrXzdebKm1AgpErz9M2FwHQCE" {
		return "", 0
	}
	return "attestClaimTransactionEncoded", claimId
}

func (provider *TestProvider) GetAttestAccountCreateTransaction(sender, amount, destination, signatureReward, bridgeId string) (string, uint64) {
	provider.GetAttestCreateAccountTxCalledTimes += 1
	if destination == "0x177adf17f5ac5df0178a24ba5b805a88a7a4be2a" || destination == "rs99jCuSAjrXzdebKm1AgpErz9M2FwHQCE" {
		return "", 0
	}
	return "attestCreateAccountTransactionEncoded", 1
}

func (provider *TestProvider) SignTransaction(transaction string) string {
	if transaction == "fail" {
		return ""
	}
	return transaction
}

func (provider *TestProvider) SignEncodedCreateBridgeTransaction(encodedTx string, isLocking bool, minBridgeReward, maxBridgeReward uint64, otherChainAddress, tokenAddress, tokenCode string) (string, string, string, error) {
	return "", "", "", nil
}

func (provider *TestProvider) GetNoOpTransaction(nonce uint, gasPriceFactor uint) string {
	provider.GetNoOpTransactionCalledTimes += 1
	if nonce == 0 {
		return ""
	}
	return "noOpTransactionEncoded"
}

func (provider *TestProvider) SetTransactionGasPrice(transaction string, factor uint) string {
	provider.SetTransactionGasPriceCalledTimes += 1
	if transaction == "fail" {
		return ""
	}
	return "transactionWithGasPriceFactor" + fmt.Sprint(factor)

}

func (provider *TestProvider) GetTransactionStatus(hash string) string {
	if strings.Contains(hash, PendingStatus) {
		return PendingStatus
	}
	if strings.Contains(hash, NotFoundStatus) {
		return NotFoundStatus
	}
	if strings.Contains(hash, ErrorStatus) {
		return ErrorStatus
	}
	if strings.Contains(hash, FailedStatus) {
		return FailedStatus
	}
	if strings.Contains(hash, AcceptedStatus) {
		return AcceptedStatus
	}
	if strings.Contains(hash, ConfirmedStatus) {
		return ConfirmedStatus
	}
	if strings.Contains(hash, UnconfirmedStatus) {
		return UnconfirmedStatus
	}
	return AcceptedStatus

}

func (provider *TestProvider) GetCurrentBlockNumber() uint64 {
	provider.GetCurrentBlockCalledTimes += 1
	if provider.BlockNumber == 0 {
		return 0
	}
	return provider.BlockNumber + 10
}

func (provider *TestProvider) SetCurrentBlockNumber(currentBlock uint64) {
	provider.SetCurrentBlockCalledTimes += 1
	provider.BlockNumber = currentBlock
}

func (provider *TestProvider) SetNewBridgesCurrentBlockNumber(currentBlock uint64) {
	provider.SetNewBridgesCurrentBlockCalledTimes += 1
	provider.NewBridgesBlockNumber = currentBlock
}

func (provider *TestProvider) GetNewCommits(toBlock uint64) interface{} {
	provider.GetNewCommitsCalledTimes += 1
	if toBlock < 100 {
		return nil
	}

	destination := "mockDestination"
	if toBlock < 1000 {
		commits := []xrp.XrpCommit{}
		commits = append(commits, xrp.XrpCommit{
			Block:       provider.BlockNumber,
			ClaimId:     1,
			Sender:      "mockAccount",
			Amount:      "100",
			Destination: &destination,
		})
		return commits
	}

	commits := []evm.EvmCommit{}
	commits = append(commits, evm.EvmCommit{
		Block:       provider.BlockNumber,
		ClaimId:     2,
		Sender:      "mockSender",
		Amount:      "150",
		Destination: &destination,
	})
	return commits
}

func (provider *TestProvider) GetNewAccountCreates(toBlock uint64) interface{} {
	provider.GetNewAccountCreatesCalledTimes += 1
	if toBlock < 100 {
		return nil
	}

	if toBlock < 1000 {
		commits := []xrp.XrpAccountCreate{}
		commits = append(commits, xrp.XrpAccountCreate{
			Block:           provider.BlockNumber,
			SignatureReward: "1000",
			Sender:          "mockAccount",
			Amount:          "100",
			Destination:     "mockDestination",
		})
		return commits
	}

	commits := []evm.EvmAccountCreate{}
	commits = append(commits, evm.EvmAccountCreate{
		Block:           provider.BlockNumber,
		Sender:          "mockSender",
		Amount:          "150",
		Destination:     "mockDestination",
		SignatureReward: "1000",
	})
	return commits
}

func (provider *TestProvider) GetUnattestedClaimById(claimId uint64, bridgeId string) (interface{}, error) {
	provider.GetUnattestedClaimCalledTimes += 1
	if claimId == 0 {
		return nil, fmt.Errorf("error getting unattested claim by id")
	}
	if claimId == 1 {
		xrpClaim := xrp.XrpClaim{ClaimId: 1, Sender: "claimCreator", Source: "r3an6Cz2MgHQT9q3Kj3QzwQv9ARkgkkxqo"}
		return xrpClaim, nil
	}
	if claimId == 2 {
		evmClaim := evm.EvmClaim{ClaimId: 2, Sender: "claimSender", Source: "0x4DBeE27B94c970B6A7916628236ad6D9369a4518"}
		return evmClaim, nil
	}
	return nil, nil
}

func (provider *TestProvider) GetChainId() *big.Int {
	return provider.chainId
}

func (provider *TestProvider) GetNonce() *uint {
	return provider.Nonce
}

func (provider *TestProvider) IsInSignerList() bool {
	return provider.isInSignerList
}

func (provider *TestProvider) CheckWitnessHasAttestedCreateAccount(destination, bridgeId string) (bool, error) {
	provider.CheckWitnessAttestedCalledTimes += 1
	if destination == "error-attested" {
		return false, fmt.Errorf("error checking witness has attested create account")
	}
	if destination == "found-attested" {
		return true, nil
	}
	return false, nil
}

func (provider *TestProvider) CheckAccountCreated(account, bridgeId string) (bool, error) {
	provider.CheckAccountCreatedCalledTimes += 1
	if account == "error-account" {
		return false, fmt.Errorf("error checking account created")
	}
	if account == "found-account" {
		return true, nil
	}
	return false, nil
}

func (provider *TestProvider) GetCurrentCreateAccountCount() uint64 {
	return provider.AccountCount
}

func (provider *TestProvider) GetTransactionCreateCount(payload string) (uint64, error) {
	if payload == "error" {
		return 0, fmt.Errorf("error checking field")
	}
	return 150, nil
}

func (provider *TestProvider) ConvertToDecimal(payload, bridgeId string) *big.Float {
	return utils.IntToFloatPrec(payload, 10)
}

func (provider *TestProvider) ConvertToWhole(payload *big.Float, bridgeId string) string {
	return utils.FloatToIntPrec(payload, 10).String()
}

func (provider *TestProvider) SetBridgeValidated(bridgeId string) interface{} {
	return nil
}

func (provider *TestProvider) GetUnpairedBridges() interface{} {
	return nil
}

func (provider *TestProvider) FetchNewBridges(toBlock uint64) error {
	return nil
}

func (provider *TestProvider) FetchNewBridgeRequests(toBlock uint64) (interface{}, error) {
	return nil, nil
}

func (provider *TestProvider) RetryNewBridgeRequest(bridgeRequestCounter interface{}) error {
	return nil
}

func (provider *TestProvider) GetType() config.ChainType {
	return provider.ChainType
}

func (provider *TestProvider) GetTokenCodeFromAddress(address string) (string, error) {
	return address, nil
}

var XrpTestProvider *TestProvider

func StartXrpTestProvider(blockNumber, accountCount uint64, inSignerList bool, chainId *big.Int, nonce *uint) {
	XrpTestProvider = &TestProvider{BlockNumber: blockNumber, isInSignerList: inSignerList, chainId: chainId, Nonce: nonce, AccountCount: accountCount, ChainType: config.Xrp}
	//mainChainProvider = XrpTestProvider
}

var EvmTestProvider *TestProvider

func StartEvmTestProvider(blockNumber, accountCount uint64, inSignerList bool, chainId *big.Int, nonce *uint) {
	EvmTestProvider = &TestProvider{BlockNumber: blockNumber, isInSignerList: inSignerList, chainId: chainId, Nonce: nonce, AccountCount: accountCount, ChainType: config.Evm}
	//sideChainProvider = EvmTestProvider
}
