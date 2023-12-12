package chains

import (
	"errors"
	"math/big"

	"peersyst/bridge-witness-go/internal/chains/evm"
	"peersyst/bridge-witness-go/internal/chains/xrp"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"peersyst/bridge-witness-go/internal/signer"
	aws "peersyst/bridge-witness-go/internal/signer/aws_kms"
	"strings"

	config "peersyst/bridge-witness-go/configs"

	"github.com/rs/zerolog/log"
	"golang.org/x/exp/maps"
)

const (
	PendingStatus     string = "Pending"
	NotFoundStatus    string = "NotFound"
	ErrorStatus       string = "Error"
	FailedStatus      string = "Failed"
	AcceptedStatus    string = "Accepted"
	ConfirmedStatus   string = "Confirmed"
	UnconfirmedStatus string = "Unconfirmed"
)

const (
	DecodingError   string = "decoding error:"
	NoResponseError string = "no response error:"
	InvalidNonce    string = "invalid nonce"
	IgnorableError  string = "ignorable error:"
	UnknownError    string = "unknown error:"
)

type ChainProvider interface {
	BroadcastTransaction(payload string) (string, error)
	GetAttestClaimTransaction(claimId uint64, sender, amount, destination, bridgeId string) (string, uint64)
	GetAttestAccountCreateTransaction(sender, amount, destination, signatureReward, bridgeId string) (string, uint64)
	SignTransaction(transaction string) string
	SignEncodedCreateBridgeTransaction(encodedTx string, isLocking bool, minBridgeReward, maxBridgeReward uint64, otherChainAddress, tokenAddress, currency string) (string, string, string, error)
	GetNoOpTransaction(nonce uint, gasPriceFactor uint) string
	SetTransactionGasPrice(transaction string, factor uint) string
	GetTransactionStatus(hash string) string
	GetCurrentBlockNumber() uint64
	SetCurrentBlockNumber(currentBlock uint64)
	SetNewBridgesCurrentBlockNumber(currentBlock uint64)
	GetNewCommits(toBlock uint64) interface{}
	GetNewAccountCreates(toBlock uint64) interface{}
	FetchNewBridges(toBlock uint64) error
	FetchNewBridgeRequests(toBlock uint64) (interface{}, error)
	RetryNewBridgeRequest(bridgeRequestCounter interface{}) error
	GetUnattestedClaimById(claimId uint64, bridgeId string) (interface{}, error)
	GetChainId() *big.Int
	GetNonce() *uint
	IsInSignerList() bool
	CheckWitnessHasAttestedCreateAccount(destination, bridgeId string) (bool, error)
	CheckAccountCreated(account, bridgeId string) (bool, error)
	GetCurrentCreateAccountCount() uint64
	GetTransactionCreateCount(payload string) (uint64, error)
	ConvertToDecimal(payload, bridgeId string) *big.Float
	ConvertToWhole(payload *big.Float, bridgeId string) string
	SetBridgeValidated(bridgeId string) interface{}
	GetUnpairedBridges() interface{}
	GetType() config.ChainType
	UpdateOracleData(amount, amount2 int64) error
	GetAmmInfo(asset *xrpl.AmmAsset, asset2 *xrpl.AmmAsset) (*xrpl.AmmInfoResult, error)
	GetTokenCodeFromAddress(address string) (string, error)
}

type BridgeProvider interface {
	GetId() string
	GetLockingChainDoor() string
	GetLockingChainCurrency() string
	GetLockingChainIssuer() string
	GetIssuingChainDoor() string
	GetIssuingChainCurrency() string
	GetIssuingChainIssuer() string
}

var mainChainProvider ChainProvider
var sideChainProvider ChainProvider

func StartMainChainProvider(cfg config.ChainConfig, signer signer.SignerProvider) (ChainProvider, error) {
	switch cfg.Type {
	case config.Xrp:
		provider, err := xrp.Create(signer, cfg.Node, cfg.DoorAddress, cfg.StartingBlock, cfg.SignerListSeconds, cfg.MaxGasFactor)
		mainChainProvider = provider
		return mainChainProvider, err
	case config.Evm:
		provider, err := evm.Create(signer, cfg.Node, cfg.DoorAddress, cfg.StartingBlock, cfg.SignerListSeconds, cfg.MaxGasFactor)
		mainChainProvider = provider
		return mainChainProvider, err
	}

	return nil, errors.New("invalid config type")
}

func GetMainChainProvider() ChainProvider {
	if mainChainProvider == nil {
		log.Error().Msgf("Error: mainChainProvider not instantiated")
	}
	return mainChainProvider
}

func StartSideChainProvider(cfg config.ChainConfig, signer signer.SignerProvider) (ChainProvider, error) {
	switch cfg.Type {
	case config.Xrp:
		provider, err := xrp.Create(signer, cfg.Node, cfg.DoorAddress, cfg.StartingBlock, cfg.SignerListSeconds, cfg.MaxGasFactor)
		sideChainProvider = provider
		return sideChainProvider, err
	case config.Evm:
		provider, err := evm.Create(signer, cfg.Node, cfg.DoorAddress, cfg.StartingBlock, cfg.SignerListSeconds, cfg.MaxGasFactor)
		sideChainProvider = provider
		return sideChainProvider, err
	}

	return nil, errors.New("invalid config type")
}

func ValidateBridges() bool {
	validatedBridgeIds := map[string]bool{}
	validatedBridges := 0

	if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Xrp {
		mainChainBridgeMap, mainChainMapIsXrp := mainChainProvider.GetUnpairedBridges().(map[string]*xrp.XrpBridgeProvider)
		sideChainBridgeMap, sideChainMapIsXrp := sideChainProvider.GetUnpairedBridges().(map[string]*xrp.XrpBridgeProvider)
		if !mainChainMapIsXrp || !sideChainMapIsXrp {
			return false
		}

		mainChainBridges := maps.Values(mainChainBridgeMap)
		for _, mainChainBridge := range mainChainBridges {
			sideChainBridge, exists := sideChainBridgeMap[mainChainBridge.GetId()]
			if exists && CheckSameChainEqualBridges(mainChainBridge, sideChainBridge) {
				validatedBridgeIds[mainChainBridge.GetId()] = true
			}
		}

		sideChainBridges := maps.Values(sideChainBridgeMap)
		for _, sideChainBridge := range sideChainBridges {
			mainChainBridge, exists := mainChainBridgeMap[sideChainBridge.GetId()]
			if exists && CheckSameChainEqualBridges(mainChainBridge, sideChainBridge) {
				if validatedBridgeIds[sideChainBridge.GetId()] {
					validatedBridges += 1
					mainChainProvider.SetBridgeValidated(sideChainBridge.GetId())
					sideChainProvider.SetBridgeValidated(sideChainBridge.GetId())
				}
			}
		}

		return true
	} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Evm {
		mainChainBridgeMap, mainChainMapIsEvm := mainChainProvider.GetUnpairedBridges().(map[string]*evm.EvmBridgeProvider)
		sideChainBridgeMap, sideChainMapIsEvm := sideChainProvider.GetUnpairedBridges().(map[string]*evm.EvmBridgeProvider)
		if !mainChainMapIsEvm || !sideChainMapIsEvm {
			return false
		}

		mainChainBridges := maps.Values(mainChainBridgeMap)
		for _, mainChainBridge := range mainChainBridges {
			sideChainBridge, exists := sideChainBridgeMap[mainChainBridge.GetId()]
			if exists && CheckSameChainEqualBridges(mainChainBridge, sideChainBridge) {
				validatedBridgeIds[mainChainBridge.GetId()] = true
			}
		}

		sideChainBridges := maps.Values(sideChainBridgeMap)
		for _, sideChainBridge := range sideChainBridges {
			mainChainBridge, exists := mainChainBridgeMap[sideChainBridge.GetId()]
			if exists && CheckSameChainEqualBridges(mainChainBridge, sideChainBridge) {
				if validatedBridgeIds[sideChainBridge.GetId()] {
					validatedBridges += 1
					mainChainProvider.SetBridgeValidated(sideChainBridge.GetId())
					sideChainProvider.SetBridgeValidated(sideChainBridge.GetId())
				}
			}
		}

		return true
	} else {
		var xrpBridgeMap map[string]*xrp.XrpBridgeProvider
		var evmBridgeMap map[string]*evm.EvmBridgeProvider

		if mainChainProvider.GetType() == config.Xrp && sideChainProvider.GetType() == config.Evm {
			mainChainBridgeMap, mainChainMapIsXrp := mainChainProvider.GetUnpairedBridges().(map[string]*xrp.XrpBridgeProvider)
			sideChainBridgeMap, sideChainMapIsEvm := sideChainProvider.GetUnpairedBridges().(map[string]*evm.EvmBridgeProvider)

			if !mainChainMapIsXrp || !sideChainMapIsEvm {
				return false
			}
			xrpBridgeMap = mainChainBridgeMap
			evmBridgeMap = sideChainBridgeMap
		} else if mainChainProvider.GetType() == config.Evm && sideChainProvider.GetType() == config.Xrp {
			mainChainBridgeMap, mainChainMapIsEvm := mainChainProvider.GetUnpairedBridges().(map[string]*evm.EvmBridgeProvider)
			sideChainBridgeMap, sideChainMapIsXrp := sideChainProvider.GetUnpairedBridges().(map[string]*xrp.XrpBridgeProvider)

			if !mainChainMapIsEvm || !sideChainMapIsXrp {
				return false
			}
			evmBridgeMap = mainChainBridgeMap
			xrpBridgeMap = sideChainBridgeMap
		} else {
			return false
		}

		xrpBridges := maps.Values(xrpBridgeMap)
		for _, xrpBridge := range xrpBridges {
			evmBridge, exists := evmBridgeMap[xrpBridge.GetId()]
			if exists && CheckEqualBridges(xrpBridge, evmBridge) {
				validatedBridgeIds[xrpBridge.GetId()] = true
			}
		}

		evmBridges := maps.Values(evmBridgeMap)
		for _, evmBridge := range evmBridges {
			xrpBridge, exists := xrpBridgeMap[evmBridge.GetId()]
			if exists && CheckEqualBridges(xrpBridge, evmBridge) {
				if validatedBridgeIds[evmBridge.GetId()] {
					validatedBridges += 1
					mainChainProvider.SetBridgeValidated(evmBridge.GetId())
					sideChainProvider.SetBridgeValidated(evmBridge.GetId())
				}
			}
		}
	}

	log.Info().Msgf("Validated %d new bridges", validatedBridges)

	return true
}

func CheckEqualBridges(xrpBridge, evmBridge BridgeProvider) bool {
	validated := xrpBridge.GetId() == evmBridge.GetId()
	validated = validated && xrpBridge.GetLockingChainDoor() == aws.EvmAddressToXrplAccount(evmBridge.GetLockingChainDoor())
	validated = validated && strings.EqualFold(aws.XrplAccountToEvmAddress(xrpBridge.GetIssuingChainDoor()), evmBridge.GetIssuingChainDoor())
	validated = validated && xrpBridge.GetLockingChainCurrency() == evmBridge.GetLockingChainCurrency()
	validated = validated && xrpBridge.GetIssuingChainCurrency() == evmBridge.GetIssuingChainCurrency()
	if xrpBridge.GetLockingChainCurrency() != "XRP" {
		// Bridge asset is token
		validated = validated && xrpBridge.GetLockingChainIssuer() == aws.EvmAddressToXrplAccount(evmBridge.GetLockingChainIssuer())
		// This should be the same as evmBridge tokenContractAddress but can not create bridge with different issuer and door
		validated = validated && strings.EqualFold(aws.XrplAccountToEvmAddress(xrpBridge.GetIssuingChainIssuer()), evmBridge.GetIssuingChainIssuer())
	}
	return validated
}

func CheckSameChainEqualBridges(mainChainBridge, sideChainBridge BridgeProvider) bool {
	validated := mainChainBridge.GetId() == sideChainBridge.GetId()
	validated = validated && mainChainBridge.GetLockingChainDoor() == sideChainBridge.GetLockingChainDoor()
	validated = validated && mainChainBridge.GetIssuingChainDoor() == sideChainBridge.GetIssuingChainDoor()
	validated = validated && mainChainBridge.GetLockingChainCurrency() == sideChainBridge.GetLockingChainCurrency()
	validated = validated && mainChainBridge.GetIssuingChainCurrency() == sideChainBridge.GetIssuingChainCurrency()
	if mainChainBridge.GetLockingChainCurrency() != "XRP" {
		// Bridge asset is token
		validated = validated && mainChainBridge.GetLockingChainIssuer() == sideChainBridge.GetLockingChainIssuer()
		// This should be the same as sideChainBridge tokenContractAddress but can not create bridge with different issuer and door
		validated = validated && mainChainBridge.GetIssuingChainIssuer() == sideChainBridge.GetIssuingChainIssuer()
	}
	return validated
}

func GetSideChainProvider() ChainProvider {
	if sideChainProvider == nil {
		log.Error().Msgf("Error: sideChainProvider not instantiated")
	}
	return sideChainProvider
}
