package xrpl

import (
	"math"
	rippleAddressCodec "peersyst/bridge-witness-go/external/ripple_address_codec"
	external "peersyst/bridge-witness-go/external/xrpl.js"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transport"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Client struct {
	Transport *transport.Transport
}

var (
	xrplJs             *external.XrplJs
	restrictedNetworks uint64 = 1024
)

func GetXrplJs() *external.XrplJs {
	if xrplJs == nil {
		xrplJs = external.NewXrplJs()
	}

	return xrplJs
}

func Create(nodeUrl string) (*Client, error) {
	xrplJs = external.NewXrplJs()
	t, err := transport.NewTransport(nodeUrl)
	if err != nil {
		return nil, err
	}
	c := &Client{&t}
	return c, nil
}

func (c *Client) Close() {
	(*c.Transport).Close()
}

func (c *Client) GetLedgerHeader() (*LedgerHeaderResult, error) {
	out := &LedgerHeaderResult{}
	err := (*c.Transport).Call("ledger_header", out, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetLedgerIndex() (uint64, error) {
	out := &LedgerIndexResult{}
	params := &LedgerIndexCommand{"validated"}
	err := (*c.Transport).Call("ledger", out, params)
	if err != nil {
		return 0, err
	}

	return out.LedgerIndex, nil
}

func (c *Client) GetAmmInfo(asset *AmmAsset, asset2 *AmmAsset) (*AmmInfoResult, error) {
	if asset.Currency == "XRP" || asset2.Currency == "XRP" {
		out := &AmmInfoResultXRP{}

		if asset.Currency == "XRP" {
			err := (*c.Transport).Call("amm_info", out, &AmmInfoCommand{Asset: *asset, Asset2: *asset2})
			if err != nil {
				return nil, err
			}
		} else {
			err := (*c.Transport).Call("amm_info", out, &AmmInfoCommand{Asset: *asset2, Asset2: *asset})
			if err != nil {
				return nil, err
			}
		}

		return &AmmInfoResult{Amm: AmmObj{
			Amount: AmmAmount{
				Value: out.Amm.Amount,
			},
			Amount2: out.Amm.Amount2,
		}}, nil
	} else {
		out := &AmmInfoResult{}
		err := (*c.Transport).Call("amm_info", out, &AmmInfoCommand{Asset: *asset, Asset2: *asset2})
		if err != nil {
			return nil, err
		}
		return out, nil
	}
}

func (c *Client) GetAccountInfo(account string, ledgerIndex *string) (*AccountInfoResult, error) {
	out := &AccountInfoResult{}
	params := &AccountInfoCommand{account, ledgerIndex}
	err := (*c.Transport).Call("account_info", out, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetAccountObjects(account string, ledgerIndex, objectType *string) (*AccountObjectsResult, error) {
	out := &AccountObjectsResult{}
	params := &AccountObjectsCommand{account, ledgerIndex, objectType}
	err := (*c.Transport).Call("account_objects", out, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetTransaction(txHash string) (*TxResult, error) {
	out := &TxResult{}
	params := &TxCommand{txHash}
	err := (*c.Transport).Call("tx", out, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetAccountTransactions(account string, minLedger, maxLedger, limit int64, marker *Marker) (*AccountTxResult, error) {
	out := &AccountTxResult{}
	params := &AccountTxCommand{account, minLedger, maxLedger, limit, marker}
	err := (*c.Transport).Call("account_tx", out, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) GetServerInfo() (*ServerInfoResult, error) {
	out := &ServerInfoResult{}
	err := (*c.Transport).Call("server_info", out, nil)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Submit(encodedTx string) (*SubmitTxResult, error) {
	out := &SubmitTxResult{}
	params := &SubmitTxCommand{encodedTx}
	err := (*c.Transport).Call("submit", out, params)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *Client) Autofill(tx transaction.Transaction) transaction.Transaction {
	// Set classic address and tags for XAddresses
	if rippleAddressCodec.IsValidXAddress(tx.GetAccount()) {
		classic, tag := rippleAddressCodec.XAddressToClassicAddress(tx.GetAccount())
		tx.SetAccount(classic)
		if tag != nil {
			tx.SetSourceTag(tag)
		}
	}
	if tx.GetDestination() != nil && rippleAddressCodec.IsValidXAddress(*tx.GetDestination()) {
		classic, tag := rippleAddressCodec.XAddressToClassicAddress(*tx.GetDestination())
		tx.SetDestination(&classic)
		if tag != nil {
			tx.SetDestinationTag(tag)
		}
	}

	// For our transaction types we have no flags
	tx.SetFlags(0)

	if tx.GetSequence() == nil {
		tx = c.SetNextValidSequenceNumber(tx)
		if tx == nil {
			return nil
		}
	}
	if tx.GetFee() == nil {
		tx = c.CalculateFeePerTransactionType(tx)
	}
	if tx.GetLastLedgerSequence() == nil {
		tx = c.SetLastValidatedLedgerSequence(tx)
		if tx == nil {
			return nil
		}
	}
	if tx.GetNetworkID() == nil {
		tx = c.AutofillNetworkID(tx)
	}

	return tx
}

var offsetLedger uint64 = 20

func (c *Client) SetLastValidatedLedgerSequence(tx transaction.Transaction) transaction.Transaction {
	lastSeq, err := c.GetLedgerIndex()
	if err != nil {
		log.Error().Msgf("Error getting ledger index: '%+v'", err)
		return nil
	}

	tx.SetLastLedgerSequence(lastSeq + offsetLedger)
	return tx
}

func (c *Client) SetNextValidSequenceNumber(tx transaction.Transaction) transaction.Transaction {
	currentLI := "current"
	lastSeq, err := c.GetAccountInfo(tx.GetAccount(), &currentLI)
	if err != nil {
		log.Error().Msgf("Error getting account info: '%+v'", err)
		return nil
	}

	tx.SetSequence(lastSeq.AccountData.Sequence)
	return tx
}

func (c *Client) CalculateFeePerTransactionType(tx transaction.Transaction) transaction.Transaction {
	serverInfo, err := c.GetServerInfo()
	if err != nil {
		tx.SetFee("10")
	} else {
		feeInDrops := math.Ceil(serverInfo.Info.ValidatedLedger.BaseFeeXrp * 1000000)
		tx.SetFee(strconv.FormatInt(int64(feeInDrops), 10))
	}

	return tx
}

func (c *Client) AutofillNetworkID(tx transaction.Transaction) transaction.Transaction {
	serverInfo, err := c.GetServerInfo()
	if err != nil {
		log.Error().Msgf("Error getting server info: '%+v'", err)
		return tx
	}

	if serverInfo.Info.BuildVersion != "" &&
		serverInfo.Info.NetworkId > restrictedNetworks {
		tx.SetNetworkID(serverInfo.Info.NetworkId)
	} else {
		return tx
	}

	return tx
}
