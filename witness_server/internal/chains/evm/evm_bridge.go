package evm

import (
	aws "peersyst/bridge-witness-go/internal/signer/aws_kms"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type EvmBridgeProvider struct {
	bridgeId      string
	bridgeKey     [32]byte
	bridge        *XChainTypesBridgeConfig
	bridgeParams  *XChainTypesBridgeParams
	tokenAddress  *common.Address
	assetDecimals int
	isNative      bool
}

func CreateEvmBridgeProvider(client *ethclient.Client, bridgeOpts *bind.CallOpts, contract *Bridge, bridge *XChainTypesBridgeConfig, params *XChainTypesBridgeParams) *EvmBridgeProvider {
	key, err := contract.GetBridgeKey(bridgeOpts, *bridge)
	if err != nil {
		log.Error().Msgf("Error getting bridge key: '%s'", err)
		return nil
	}

	provider := EvmBridgeProvider{
		"",
		key,
		bridge,
		params,
		&zeroAddress,
		18,
		false,
	}

	err = provider.setBridge(client, bridgeOpts, contract)
	if err != nil {
		log.Error().Msgf("Error setting up bridge: '%s'", err)
		return nil
	}

	return &provider
}

func CreateEvmBridgeProviderFromEvent(client *ethclient.Client, bridgeOpts *bind.CallOpts, contract *Bridge, event *BridgeCreateBridge) *EvmBridgeProvider {
	bridge := &XChainTypesBridgeConfig{
		LockingChainDoor:  event.LockingChainDoor,
		LockingChainIssue: XChainTypesBridgeChainIssue{event.LockingChainIssueIssuer, event.LockingChainIssueCurency},
		IssuingChainDoor:  event.IssuingChainDoor,
		IssuingChainIssue: XChainTypesBridgeChainIssue{event.IssuingChainIssueIssuer, event.IssuingChainIssueCurency},
	}
	MinCreateAmount, SignatureReward, err := contract.GetBridgeParams(bridgeOpts, *bridge)
	if err != nil {
		log.Error().Msgf("Error getting bridge params: '%s'", err)
		return nil
	}
	params := &XChainTypesBridgeParams{MinCreateAmount, SignatureReward}

	provider := EvmBridgeProvider{
		"",
		event.BridgeKey,
		bridge,
		params,
		&zeroAddress,
		18,
		false,
	}

	err = provider.setBridge(client, bridgeOpts, contract)
	if err != nil {
		log.Error().Msgf("Error setting up bridge: '%s'", err)
		return nil
	}

	return &provider
}

func (provider *EvmBridgeProvider) GetId() string {
	return provider.bridgeId
}

func (provider *EvmBridgeProvider) GetLockingChainDoor() string {
	return provider.bridge.LockingChainDoor.String()
}

func (provider *EvmBridgeProvider) GetLockingChainCurrency() string {
	return provider.bridge.LockingChainIssue.Currency
}

func (provider *EvmBridgeProvider) GetLockingChainIssuer() string {
	return provider.bridge.LockingChainIssue.Issuer.String()
}

func (provider *EvmBridgeProvider) GetIssuingChainDoor() string {
	// TODO: REMOVE ON DEPLOY
	// bridge.IssuingChainIssue + "-" + bridge.IssuingChainDoor
	// We hardcode for native to validate correctly. Ideally it will be created on genesis correctly
	if provider.isNative {
		return "0xb5f762798a53d543a014caf8b297cff8f2f937e8"
	}
	return provider.bridge.IssuingChainDoor.String()
}

func (provider *EvmBridgeProvider) GetIssuingChainCurrency() string {
	return provider.bridge.IssuingChainIssue.Currency
}

func (provider *EvmBridgeProvider) GetIssuingChainIssuer() string {
	return provider.bridge.IssuingChainIssue.Issuer.String()
}

func (provider *EvmBridgeProvider) setBridge(client *ethclient.Client, bridgeOpts *bind.CallOpts, contract *Bridge) error {
	if provider.bridge.LockingChainIssue.Issuer.String() != zeroAddress.String() {
		tokenAddress, err := contract.GetBridgeToken(bridgeOpts, *provider.bridge)
		if err != nil {
			return err
		}
		provider.tokenAddress = &tokenAddress
	}

	// Must read asset decimals
	if provider.tokenAddress.String() != zeroAddress.String() {
		instance, err := NewToken(*provider.tokenAddress, client)
		if err != nil {
			return err
		}

		Decimals, err := instance.Decimals(bridgeOpts)
		if err != nil {
			return err
		}

		provider.assetDecimals = int(Decimals)
	} else {
		provider.isNative = true
	}

	provider.SetBridgeId()

	return nil
}

func (provider *EvmBridgeProvider) SetBridgeId() {
	if provider.GetIssuingChainCurrency() == "XRP" || provider.GetLockingChainCurrency() == "XRP" {
		provider.bridgeId = strings.ToLower(provider.GetLockingChainCurrency() + "-" + aws.EvmAddressToXrplAccount(provider.GetLockingChainDoor()) + ":" +
			provider.GetIssuingChainCurrency() + "-" + provider.GetIssuingChainDoor())
	} else {
		provider.bridgeId = strings.ToLower(provider.GetLockingChainCurrency() + "-" + aws.EvmAddressToXrplAccount(provider.GetLockingChainIssuer()) + ":" +
			provider.GetIssuingChainCurrency() + "-" + provider.GetIssuingChainIssuer())
	}
}
