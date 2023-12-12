package xrp

import (
	"bytes"
	"encoding/hex"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"strings"
	"sync/atomic"

	aws "peersyst/bridge-witness-go/internal/signer/aws_kms"
)

type XrpBridgeProvider struct {
	bridgeId           string
	accountCreateCount uint64
	bridge             *transaction.XChainBridge
	isLocking          bool
	isToken            bool
}

func CreateXrpBridgeProvider(doorAddress string, bridge *transaction.XChainBridge, createCount uint64) *XrpBridgeProvider {
	isLocking := true
	isToken := false
	if bridge.IssuingChainDoor == doorAddress {
		isLocking = false
	}
	if bridge.IssuingChainIssue.Issuer != nil {
		isToken = true
	}

	bridgeProvider := XrpBridgeProvider{
		bridgeId:           GetIdFromBridge(bridge),
		accountCreateCount: createCount,
		bridge:             bridge,
		isLocking:          isLocking,
		isToken:            isToken,
	}

	return &bridgeProvider
}

func (provider *XrpBridgeProvider) GetId() string {
	return provider.bridgeId
}

func (provider *XrpBridgeProvider) GetLockingChainDoor() string {
	return provider.bridge.LockingChainDoor
}

func (provider *XrpBridgeProvider) GetLockingChainCurrency() string {
	lckCurrency := provider.bridge.LockingChainIssue.Currency
	if len(lckCurrency) != 3 {
		lockCurrencyByte, err := hex.DecodeString(lckCurrency)
		if err != nil {
			return ""
		}
		lckCurrency = string(bytes.Trim(lockCurrencyByte, "\x00"))
	}
	return lckCurrency
}

func (provider *XrpBridgeProvider) GetLockingChainIssuer() string {
	return *provider.bridge.LockingChainIssue.Issuer
}

func (provider *XrpBridgeProvider) GetIssuingChainDoor() string {
	return provider.bridge.IssuingChainDoor
}

func (provider *XrpBridgeProvider) GetIssuingChainCurrency() string {
	issCurrency := provider.bridge.IssuingChainIssue.Currency
	if len(issCurrency) != 12 && len(issCurrency) != 3 {
		lockCurrencyByte, err := hex.DecodeString(issCurrency)
		if err != nil {
			return ""
		}
		issCurrency = string(bytes.Trim(lockCurrencyByte, "\x00"))
	}
	return issCurrency
}

func (provider *XrpBridgeProvider) GetIssuingChainIssuer() string {
	return *provider.bridge.IssuingChainIssue.Issuer
}

func (provider *XrpBridgeProvider) consumeAccountCreateCount() uint64 {
	return atomic.AddUint64(&provider.accountCreateCount, 1)
}

func bridgesEqual(bridge1, bridge2 *transaction.XChainBridge) bool {
	return bridge1.IssuingChainDoor == bridge2.IssuingChainDoor &&
		bridge1.LockingChainDoor == bridge2.LockingChainDoor &&
		bridge1.IssuingChainIssue.Currency == bridge2.IssuingChainIssue.Currency &&
		bridge1.LockingChainIssue.Currency == bridge2.LockingChainIssue.Currency &&
		// Both bridges have Issuer nil pointer
		((bridge1.IssuingChainIssue.Issuer == nil && bridge2.IssuingChainIssue.Issuer == nil) ||
			// Both bridges have Issuer not nil pointer and same values
			(bridge1.IssuingChainIssue.Issuer != nil &&
				bridge2.IssuingChainIssue.Issuer != nil &&
				*bridge1.IssuingChainIssue.Issuer == *bridge2.IssuingChainIssue.Issuer)) &&
		// Both bridges have Issuer not nil pointer and same values
		((bridge1.LockingChainIssue.Issuer == nil && bridge2.LockingChainIssue.Issuer == nil) ||
			// Both bridges have Issuer not nil pointer and same values
			(bridge1.LockingChainIssue.Issuer != nil &&
				bridge2.LockingChainIssue.Issuer != nil &&
				*bridge1.LockingChainIssue.Issuer == *bridge2.LockingChainIssue.Issuer))
}

func GetIdFromBridge(bridge *transaction.XChainBridge) string {
	if bridge.LockingChainIssue.Issuer == nil || bridge.IssuingChainIssue.Issuer == nil || bridge.LockingChainIssue.Currency == "XRP" || bridge.IssuingChainIssue.Currency == "XRP" {
		return strings.ToLower(bridge.LockingChainIssue.Currency + "-" + bridge.LockingChainDoor + ":" +
			bridge.IssuingChainIssue.Currency + "-" + aws.XrplAccountToEvmAddress(bridge.IssuingChainDoor))
	}

	lckCurrency := bridge.LockingChainIssue.Currency
	issCurrency := bridge.IssuingChainIssue.Currency
	if len(lckCurrency) != 3 {
		lockCurrencyByte, err := hex.DecodeString(lckCurrency)
		if err != nil {
			return ""
		}
		lckCurrency = string(bytes.Trim(lockCurrencyByte, "\x00"))
	}
	if len(issCurrency) != 12 && len(issCurrency) != 3 {
		lockCurrencyByte, err := hex.DecodeString(issCurrency)
		if err != nil {
			return ""
		}
		issCurrency = string(bytes.Trim(lockCurrencyByte, "\x00"))
	}
	return strings.ToLower(lckCurrency + "-" + *bridge.LockingChainIssue.Issuer + ":" +
		issCurrency + "-" + aws.XrplAccountToEvmAddress(*bridge.IssuingChainIssue.Issuer))
}
