package xrpl

import "peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"

type AccountInfoCommand struct {
	Account     string  `json:"account"`
	LedgerIndex *string `json:"ledger_index,omitempty"`
}

type AmmAsset struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer,omitempty"`
}

type AmmAmount struct {
	AmmAsset
	Value string `json:"value"`
}

type AmmInfoCommand struct {
	Asset  AmmAsset `json:"asset"`
	Asset2 AmmAsset `json:"asset2"`
}

type AmmObjXRP struct {
	Amount  string    `json:"amount"`
	Amount2 AmmAmount `json:"amount2"`
}

type AmmObj struct {
	Amount  AmmAmount `json:"amount"`
	Amount2 AmmAmount `json:"amount2"`
}

type AmmInfoResultXRP struct {
	Amm AmmObjXRP `json:"amm"`
}

type AmmInfoResult struct {
	Amm AmmObj `json:"amm"`
}

type AccountObjectsCommand struct {
	Account     string  `json:"account"`
	LedgerIndex *string `json:"ledger_index,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type BaseData struct {
	LedgerEntryType   string `json:",omitempty"`
	LedgerIndex       string `json:"index,omitempty"`
	PreviousTxnID     string `json:",omitempty"`
	PreviousTxnLgrSeq int    `json:",omitempty"`
}

type AccountData struct {
	BaseData
	Account    string      `json:",omitempty"`
	Balance    interface{} `json:",omitempty"`
	Flags      int         `json:",omitempty"`
	OwnerCount int         `json:",omitempty"`
	Sequence   uint64      `json:",omitempty"`
}

type XChainBridgeObject struct {
	Account                  string                    `json:"Account,omitempty"`
	Flags                    uint64                    `json:"Flags,omitempty"`
	LedgerEntryType          string                    `json:"LedgerEntryType,omitempty"`
	MinAccountCreateAmount   string                    `json:"MinAccountCreateAmount,omitempty"`
	OwnerNode                string                    `json:"OwnerNode,omitempty"`
	PreviousTxnID            string                    `json:"PreviousTxnID,omitempty"`
	PreviousTxnLgrSeq        uint64                    `json:"PreviousTxnLgrSeq,omitempty"`
	SignatureReward          string                    `json:"SignatureReward,omitempty"`
	XChainAccountClaimCount  string                    `json:"XChainAccountClaimCount,omitempty"`
	XChainAccountCreateCount string                    `json:"XChainAccountCreateCount,omitempty"`
	XChainBridge             *transaction.XChainBridge `json:"XChainBridge,omitempty"`
	XChainClaimID            string                    `json:"XChainClaimID,omitempty"`
	Index                    string                    `json:"index,omitempty"`
}

type XChainClaimProofSigElem struct {
	Amount                   string `json:"Amount"`
	AttestationRewardAccount string `json:"AttestationRewardAccount"`
	AttestationSignerAccount string `json:"AttestationSignerAccount"`
	Destination              string `json:"Destination"`
	WasLockingChainSend      uint64 `json:"WasLockingChainSend"`
}

type XChainClaimProofSig struct {
	XChainClaimProofSig *XChainClaimProofSigElem `json:"XChainClaimProofSig,omitempty"`
}

type XChainClaimObject struct {
	Account                 string                    `json:"Account,omitempty"`
	Flags                   uint64                    `json:"Flags,omitempty"`
	LedgerEntryType         string                    `json:"LedgerEntryType,omitempty"`
	OtherChainSource        string                    `json:"OtherChainSource,omitempty"`
	OwnerNode               string                    `json:"OwnerNode,omitempty"`
	PreviousTxnID           string                    `json:"PreviousTxnID,omitempty"`
	PreviousTxnLgrSeq       uint64                    `json:"PreviousTxnLgrSeq,omitempty"`
	SignatureReward         string                    `json:"SignatureReward,omitempty"`
	XChainBridge            *transaction.XChainBridge `json:"XChainBridge,omitempty"`
	XChainClaimID           string                    `json:"XChainClaimID,omitempty"`
	XChainClaimAttestations []XChainClaimProofSig     `json:"XChainClaimAttestations,omitempty"`
	Index                   string                    `json:"index,omitempty"`
}

type XChainCreateAccountProofSigElem struct {
	Amount                   string `json:"Amount"`
	AttestationRewardAccount string `json:"AttestationRewardAccount"`
	AttestationSignerAccount string `json:"AttestationSignerAccount"`
	Destination              string `json:"Destination"`
	SignatureReward          string `json:"SignatureReward"`
	WasLockingChainSend      uint64 `json:"WasLockingChainSend"`
}

type XChainCreateAccountProofSig struct {
	XChainCreateAccountProofSig *XChainCreateAccountProofSigElem `json:"XChainCreateAccountProofSig,omitempty"`
}

type XChainCreateAccountObject struct {
	Account                         string                        `json:"Account,omitempty"`
	Flags                           uint64                        `json:"Flags,omitempty"`
	LedgerEntryType                 string                        `json:"LedgerEntryType,omitempty"`
	OwnerNode                       string                        `json:"OwnerNode,omitempty"`
	PreviousTxnID                   string                        `json:"PreviousTxnID,omitempty"`
	PreviousTxnLgrSeq               uint64                        `json:"PreviousTxnLgrSeq,omitempty"`
	XChainAccountCreateCount        string                        `json:"XChainAccountCreateCount,omitempty"`
	XChainBridge                    *transaction.XChainBridge     `json:"XChainBridge,omitempty"`
	XChainCreateAccountAttestations []XChainCreateAccountProofSig `json:"XChainCreateAccountAttestations,omitempty"`
	Index                           string                        `json:"index,omitempty"`
}

type SignerEntryElem struct {
	Account      string `json:"Account,omitempty"`
	SignerWeight uint64 `json:"SignerWeight,omitempty"`
}

type SignerEntry struct {
	SignerEntry SignerEntryElem `json:"SignerEntry,omitempty"`
}

type SignerListObject struct {
	Flags             uint64         `json:"Flags,omitempty"`
	LedgerEntryType   string         `json:"LedgerEntryType,omitempty"`
	OwnerNode         string         `json:"OwnerNode,omitempty"`
	PreviousTxnID     string         `json:"PreviousTxnID,omitempty"`
	PreviousTxnLgrSeq uint64         `json:"PreviousTxnLgrSeq,omitempty"`
	SignerEntries     []*SignerEntry `json:"SignerEntries,omitempty"`
	SignerListID      uint64         `json:"SignerListID,omitempty"`
	SignerQuorum      uint64         `json:"SignerQuorum,omitempty"`
	Index             string         `json:"index,omitempty"`
}

type AccountObjectsResult struct {
	Account            string        `json:"account"`
	Objects            []interface{} `json:"account_objects,omitempty"`
	LedgerCurrentIndex uint32        `json:"ledger_current_index,omitempty"`
	LedgerIndex        uint32        `json:"ledger_index,omitempty"`
	LedgerHash         string        `json:"ledger_hash,omitempty"`
	Validated          bool          `json:"validated"`
}

type AccountInfoResult struct {
	LedgerCurrentIndex uint32      `json:"ledger_current_index,omitempty"`
	LedgerIndex        uint32      `json:"ledger_index,omitempty"`
	LedgerHash         string      `json:"ledger_hash,omitempty"`
	Validated          bool        `json:"validated"`
	AccountData        AccountData `json:"account_data"`
}

type Ledger struct {
	Accepted            bool   `json:"accepted,omitempty"`
	AccountHash         string `json:"account_hash,omitempty"`
	CloseFlags          uint64 `json:"close_flags,omitempty"`
	CloseTime           uint64 `json:"close_time,omitempty"`
	CloseTimeHuman      string `json:"close_time_human,omitempty"`
	CloseTimeResolution uint64 `json:"close_time_resolution,omitempty"`
	Closed              bool   `json:"closed,omitempty"`
	Hash                string `json:"hash,omitempty"`
	LedgerHash          string `json:"ledger_hash,omitempty"`
	LedgerIndex         string `json:"ledger_index,omitempty"`
	ParentCloseTime     uint64 `json:"parent_close_time,omitempty"`
	ParentHash          string `json:"parent_hash,omitempty"`
	SeqNum              string `json:",omitempty"`
	TotalCoins          string `json:",omitempty"`
	TransactionHash     string `json:"transaction_hash,omitempty"`
}

type LedgerHeaderResult struct {
	Ledger             Ledger `json:"ledger,omitempty"`
	LedgerCurrentIndex uint32 `json:"ledger_current_index,omitempty"`
	LedgerData         string `json:"ledger_data,omitempty"`
	Validated          bool   `json:"validated,omitempty"`
}

type LedgerIndexCommand struct {
	LedgerIndex string `json:"ledger_index,omitempty"`
}

type LedgerIndexResult struct {
	Ledger      Ledger `json:"ledger,omitempty"`
	LedgerHash  string `json:"ledger_hash,omitempty"`
	LedgerIndex uint64 `json:"ledger_index,omitempty"`
	Validated   bool   `json:"validated,omitempty"`
}

type Marker struct {
	Ledger   uint64 `json:"ledger"`
	Sequence uint64 `json:"seq"`
}

// Incompleted, only parsed some for get claim id usage
type Fields struct {
	Account                  *string                   `json:"Account,omitempty"`
	Balance                  *interface{}              `json:"Balance,omitempty"`
	Flags                    *uint64                   `json:"Flags,omitempty"`
	MinAccountCreateAmount   *string                   `json:"MinAccountCreateAmount,omitempty"`
	OtherChainSource         *string                   `json:"OtherChainSource,omitempty"`
	OwnerNode                *string                   `json:"OwnerNode,omitempty"`
	XChainAccountClaimCount  *string                   `json:"XChainAccountClaimCount,omitempty"`
	XChainAccountCreateCount *string                   `json:"XChainAccountCreateCount,omitempty"`
	SignatureReward          *string                   `json:"SignatureReward,omitempty"`
	XChainBridge             *transaction.XChainBridge `json:"XChainBridge,omitempty"`
	XChainClaimID            *string                   `json:"XChainClaimID,omitempty"`
	Owner                    *string                   `json:"Owner,omitempty"`
	OwnerCount               *uint64                   `json:"OwnerCount,omitempty"`
	RootIndex                *string                   `json:"RootIndex,omitempty"`
	Sequence                 *uint64                   `json:"Sequence,omitempty"`
}

type AffectedNode struct {
	FinalFields       *Fields `json:",omitempty"`
	LedgerEntryType   string  `json:",omitempty"`
	LedgerIndex       string  `json:",omitempty"`
	PreviousFields    *Fields `json:",omitempty"`
	NewFields         *Fields `json:",omitempty"`
	PreviousTxnID     string  `json:",omitempty"`
	PreviousTxnLgrSeq *uint32 `json:",omitempty"`
}

type NodeEffect struct {
	ModifiedNode *AffectedNode `json:",omitempty"`
	CreatedNode  *AffectedNode `json:",omitempty"`
	DeletedNode  *AffectedNode `json:",omitempty"`
}

type Metadata struct {
	AffectedNodes     []NodeEffect
	TransactionIndex  uint32
	TransactionResult string
	DeliveredAmount   string `json:"delivered_amount,omitempty"`
}

type TransactionAndMetadata struct {
	MetaData    Metadata                      `json:"meta,omitempty"`
	Transaction transaction.TransactionStruct `json:"tx"`
}

type AccountTxCommand struct {
	Account   string  `json:"account"`
	MinLedger int64   `json:"ledger_index_min"`
	MaxLedger int64   `json:"ledger_index_max"`
	Limit     int64   `json:"limit,omitempty"`
	Marker    *Marker `json:"marker,omitempty"`
}

type AccountTxResult struct {
	AccountTxCommand
	Transactions []TransactionAndMetadata `json:"transactions,omitempty"`
	Validated    bool                     `json:"validated"`
}

type TxCommand struct {
	Transaction string `json:"transaction"`
}

type TxResult struct {
	transaction.TransactionStruct
	MetaData  Metadata `json:"meta,omitempty"`
	Validated bool     `json:"validated"`
}

type SubmitTxCommand struct {
	TxBlob string `json:"tx_blob"`
}

type SubmitTxResult struct {
	EngineResult        string                        `json:"engine_result"`
	EngineResultCode    int                           `json:"engine_result_code"`
	EngineResultMessage string                        `json:"engine_result_message"`
	TxBlob              string                        `json:"tx_blob"`
	Tx                  transaction.TransactionStruct `json:"tx_json"`
}

type ValidatedLedger struct {
	Age            uint64  `json:"age,omitempty"`
	BaseFeeXrp     float64 `json:"base_fee_xrp,omitempty"`
	Hash           string  `json:"hash,omitempty"`
	ReserveBaseXrp uint64  `json:"reserve_base_xrp,omitempty"`
	ReserveIncXrp  uint64  `json:"reserve_inc_xrp,omitempty"`
	Sequence       uint64  `json:"sequence,omitempty"`
}

type StateAccountingFields struct {
	DurationUs  string `json:"duration_us,omitempty"`
	Transitions string `json:"transitions,omitempty"`
}

type StateAccounting struct {
	Connected    StateAccountingFields `json:"connected,omitempty"`
	Disconnected StateAccountingFields `json:"disconnected,omitempty"`
	Full         StateAccountingFields `json:"full,omitempty"`
	Syncing      StateAccountingFields `json:"syncing,omitempty"`
	Tracking     StateAccountingFields `json:"tracking,omitempty"`
}

type LastClose struct {
	ConvergeTimeS float32 `json:"converge_time_s,omitempty"`
	Proposers     uint64  `json:"proposers,omitempty"`
}

type ServerInfo struct {
	BuildVersion             string          `json:"build_version,omitempty"`
	CompleteLedgers          string          `json:"complete_ledgers,omitempty"`
	HostId                   string          `json:"hostid,omitempty"`
	InitialSyncDuration      string          `json:"initial_sync_duration_us,omitempty"`
	IOLatencyMs              uint64          `json:"io_latency_ms,omitempty"`
	TransOverflow            string          `json:"jq_trans_overflow,omitempty"`
	LastClose                LastClose       `json:"last_close,omitempty"`
	LoadFactor               uint64          `json:"load_factor,omitempty"`
	NetworkId                uint64          `json:"network_id,omitempty"`
	PeerDisconnects          string          `json:"peer_disconnects,omitempty"`
	PeerDisconnectsResources string          `json:"peer_disconnects_resources,omitempty"`
	Peers                    uint64          `json:"peers,omitempty"`
	PubkeyNode               string          `json:"pubkey_node,omitempty"`
	ServerState              string          `json:"server_state,omitempty"`
	ServerStateDuration      string          `json:"server_state_duration_us,omitempty"`
	StateAccounting          StateAccounting `json:"state_accounting,omitempty"`
	Time                     string          `json:"time,omitempty"`
	Uptime                   uint64          `json:"uptime,omitempty"`
	ValidatedLedger          ValidatedLedger `json:"validated_ledger,omitempty"`
	ValidationQuorum         uint64          `json:"validation_quorum,omitempty"`
}

type ServerInfoResult struct {
	Info ServerInfo `json:"info"`
}
