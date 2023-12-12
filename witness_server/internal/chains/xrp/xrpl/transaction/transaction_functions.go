package transaction

import (
	"encoding/json"
	"strconv"

	"github.com/rs/zerolog/log"
)

// Getters
func (tx *TransactionStruct) GetAccount() string {
	return tx.Account
}

func (tx *TransactionStruct) GetTransactionType() string {
	return tx.TransactionType
}

func (tx *TransactionStruct) GetFee() *string {
	return tx.Fee
}

func (tx *TransactionStruct) GetSequence() *uint64 {
	return tx.Sequence
}

func (tx *TransactionStruct) GetFlags() uint64 {
	return *tx.Flags
}

func (tx *TransactionStruct) GetNetworkID() *uint64 {
	return tx.NetworkID
}

func (tx *TransactionStruct) GetLastLedgerSequence() *uint64 {
	return tx.LastLedgerSequence
}

func (tx *TransactionStruct) GetLedgerSequence() uint64 {
	return tx.LedgerSequence
}

func (tx *TransactionStruct) GetSigningPubKey() string {
	if tx.SigningPubKey == nil {
		return ""
	}
	return *tx.SigningPubKey
}

func (tx *TransactionStruct) GetTxnSignature() string {
	return tx.TxnSignature
}

func (tx *TransactionStruct) GetHash() string {
	return tx.Hash
}

func (tx *TransactionStruct) GetClaimId() uint64 {
	claimId, err := strconv.ParseUint(*tx.XChainClaimID, 16, 64)
	if err != nil {
		log.Error().Msgf("Error converting XChainClaimID to uint64: '%+v'", err)
		return 0
	}

	return claimId
}

func (tx *TransactionStruct) GetAccountCreateCount() uint64 {
	count, err := strconv.ParseUint(*tx.XChainAccountCreateCount, 16, 64)
	if err != nil {
		log.Error().Msgf("Error converting XChainAccountCreateCount to uint64: '%+v'", err)
		return 0
	}

	return count
}

func (tx *TransactionStruct) GetDestination() *string {
	return tx.Destination
}

func (tx *TransactionStruct) GetAmount() string {
	amountStr, isString := tx.Amount.(string)
	if isString {
		return amountStr
	}

	tokenAmount, isMap := tx.Amount.(map[string]interface{})
	if isMap {
		tknStr, isStr := tokenAmount["value"].(string)
		if isStr {
			return tknStr
		}
	}

	return ""
}

func (tx *TransactionStruct) GetSignatureReward() string {
	return tx.SignatureReward
}

func (tx *TransactionStruct) GetSourceTag() *uint32 {
	return tx.SourceTag
}

func (tx *TransactionStruct) GetDestinationTag() *uint32 {
	return tx.DestinationTag
}

func (tx *TransactionStruct) GetXChainBridge() *XChainBridge {
	return tx.XChainBridge
}

// Setters
func (tx *TransactionStruct) SetSigningPubKey(v string) {
	tx.SigningPubKey = &v
}

func (tx *TransactionStruct) SetTxnSignature(v string) {
	tx.TxnSignature = v
}

func (tx *TransactionStruct) SetSequence(seq uint64) {
	tx.Sequence = &seq
}

func (tx *TransactionStruct) SetLastLedgerSequence(llseq uint64) {
	tx.LastLedgerSequence = &llseq
}

func (tx *TransactionStruct) SetFlags(flags uint64) {
	tx.Flags = &flags
}

func (tx *TransactionStruct) SetFee(fee string) {
	tx.Fee = &fee
}

func (tx *TransactionStruct) SetNetworkID(networkID uint64) {
	tx.NetworkID = &networkID
}

func (tx *TransactionStruct) SetSourceTag(tag *uint32) {
	tx.SourceTag = tag
}

func (tx *TransactionStruct) SetDestinationTag(tag *uint32) {
	tx.DestinationTag = tag
}

func (tx *TransactionStruct) SetAccount(account string) {
	tx.Account = account
}

func (tx *TransactionStruct) SetDestination(destination *string) {
	tx.Destination = destination
}

// Global Transaction Functions
func UnmarshalTransaction(txJson string) (Transaction, error) {
	tx := TransactionStruct{}
	err := json.Unmarshal([]byte(txJson), &tx)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func MarshalTransaction(tx Transaction) (string, error) {
	encoded, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}
