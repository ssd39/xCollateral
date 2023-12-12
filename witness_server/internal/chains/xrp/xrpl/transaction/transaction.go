package transaction

type Transaction interface {
	// Getters
	GetAccount() string
	GetTransactionType() string
	GetFee() *string
	GetSequence() *uint64
	GetNetworkID() *uint64
	GetFlags() uint64
	GetLastLedgerSequence() *uint64
	GetSigningPubKey() string
	GetTxnSignature() string
	GetHash() string
	GetClaimId() uint64
	GetAccountCreateCount() uint64
	GetDestination() *string
	GetAmount() string
	GetSignatureReward() string
	GetSourceTag() *uint32
	GetDestinationTag() *uint32
	GetXChainBridge() *XChainBridge

	// Setters
	SetAccount(string)
	SetDestination(*string)
	SetSourceTag(*uint32)
	SetDestinationTag(*uint32)
	SetSequence(uint64)
	SetLastLedgerSequence(uint64)
	SetFlags(uint64)
	SetFee(string)
	SetNetworkID(uint64)
	SetSigningPubKey(string)
	SetTxnSignature(string)
}

type ChainIssue struct {
	Currency string  `json:"currency,omitempty"`
	Issuer   *string `json:"issuer,omitempty"`
}

type TokenAmount struct {
	Currency string `json:"currency,omitempty"`
	Issuer   string `json:"issuer,omitempty"`
	Value    string `json:"value,omitempty"`
}

type XChainBridge struct {
	IssuingChainDoor  string     `json:"IssuingChainDoor,omitempty"`
	IssuingChainIssue ChainIssue `json:"IssuingChainIssue,omitempty"`
	LockingChainDoor  string     `json:"LockingChainDoor,omitempty"`
	LockingChainIssue ChainIssue `json:"LockingChainIssue,omitempty"`
}

type XChainCreateAccountAttestationBatchElement struct {
	Account                  string `json:"Account,omitempty"`
	Amount                   string `json:"Amount,omitempty"`
	AttestationRewardAccount string `json:"AttestationRewardAccount,omitempty"`
	Destination              string `json:"Destination,omitempty"`
	PublicKey                string `json:"PublicKey,omitempty"`
	Signature                string `json:"Signature,omitempty"`
	SignatureReward          string `json:"SignatureReward,omitempty"`
	WasLockingChainSend      int    `json:"WasLockingChainSend,omitempty"`
	XChainAccountCreateCount string `json:"XChainAccountCreateCount,omitempty"`
}

type XChainClaimAttestationBatchElement struct {
	Account                  string `json:"Account,omitempty"`
	Amount                   string `json:"Amount,omitempty"`
	AttestationRewardAccount string `json:"AttestationRewardAccount,omitempty"`
	Destination              string `json:"Destination,omitempty"`
	PublicKey                string `json:"PublicKey,omitempty"`
	Signature                string `json:"Signature,omitempty"`
	WasLockingChainSend      int    `json:"WasLockingChainSend,omitempty"`
	XChainClaimID            string `json:"XChainClaimID,omitempty"`
}

type XChainCreateAccountAttestationBatch struct {
	XChainCreateAccountAttestationBatchElement XChainCreateAccountAttestationBatchElement `json:"XChainCreateAccountAttestationBatchElement,omitempty"`
}

type XChainClaimAttestationBatch struct {
	XChainClaimAttestationBatchElement XChainClaimAttestationBatchElement `json:"XChainClaimAttestationBatchElement,omitempty"`
}

type XChainAttestationBatch struct {
	XChainBridge                        XChainBridge                          `json:"XChainBridge"`
	XChainClaimAttestationBatch         []XChainClaimAttestationBatch         `json:"XChainClaimAttestationBatch"`
	XChainCreateAccountAttestationBatch []XChainCreateAccountAttestationBatch `json:"XChainCreateAccountAttestationBatch"`
}

type SignerEntryElem struct {
	Account      string `json:"Account,omitempty"`
	SignerWeight int    `json:"SignerWeight,omitempty"`
}

type SignerEntry struct {
	SignerEntry SignerEntryElem `json:"SignerEntry,omitempty"`
}

type SignerElem struct {
	Account       string `json:"Account,omitempty"`
	TxnSignature  string `json:"TxnSignature,omitempty"`
	SigningPubKey string `json:"SigningPubKey"`
}

type Signer struct {
	Signer SignerElem `json:"Signer,omitempty"`
}

// This transaction struct includes all possible transactions fields
// with the value omitempty to marshal and unmarshal
type TransactionStruct struct {
	Account                  string                  `json:"Account,omitempty"`
	Amount                   interface{}             `json:"Amount,omitempty"`
	AttestationRewardAccount *string                 `json:"AttestationRewardAccount,omitempty"`
	AttestationSignerAccount *string                 `json:"AttestationSignerAccount,omitempty"`
	Destination              *string                 `json:"Destination,omitempty"`
	Fee                      *string                 `json:"Fee,omitempty"`
	Flags                    *uint64                 `json:"Flags,omitempty"`
	NetworkID                *uint64                 `json:"NetworkID,omitempty"`
	LastLedgerSequence       *uint64                 `json:"LastLedgerSequence,omitempty"`
	OtherChainDestination    *string                 `json:"OtherChainDestination,omitempty"`
	OtherChainSource         *string                 `json:"OtherChainSource,omitempty"`
	Memos                    interface{}             `json:"Memos,omitempty"`
	MinAccountCreateAmount   string                  `json:"MinAccountCreateAmount,omitempty"`
	PublicKey                *string                 `json:"PublicKey,omitempty"`
	Sequence                 *uint64                 `json:"Sequence,omitempty"`
	Signature                *string                 `json:"Signature,omitempty"`
	SignatureReward          string                  `json:"SignatureReward,omitempty"`
	SignerEntries            []SignerEntry           `json:"SignerEntries,omitempty"`
	Signers                  []Signer                `json:"Signers,omitempty"`
	SignerQuorum             int                     `json:"SignerQuorum,omitempty"`
	SigningPubKey            *string                 `json:"SigningPubKey,omitempty"`
	SourceTag                *uint32                 `json:"SourceTag,omitempty"`
	DestinationTag           *uint32                 `json:"DestinationTag,omitempty"`
	TransactionType          string                  `json:"TransactionType,omitempty"`
	TxnSignature             string                  `json:"TxnSignature,omitempty"`
	XChainAttestationBatch   *XChainAttestationBatch `json:"XChainAttestationBatch,omitempty"`
	XChainBridge             *XChainBridge           `json:"XChainBridge,omitempty"`
	XChainAccountCreateCount *string                 `json:"XChainAccountCreateCount,omitempty"`
	XChainClaimID            *string                 `json:"XChainClaimID,omitempty"`
	WasLockingChainSend      *uint64                 `json:"WasLockingChainSend,omitempty"`
	Date                     uint64                  `json:"date,omitempty"`
	Hash                     string                  `json:"hash,omitempty"`
	InLedger                 uint64                  `json:"inLedger,omitempty"`
	LedgerSequence           uint64                  `json:"ledger_index,omitempty"`
}
