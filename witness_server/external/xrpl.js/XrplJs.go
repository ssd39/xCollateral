package external

import (
	"encoding/json"
	"os"
	"path/filepath"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl/transaction"
	"regexp"
	"strconv"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	v8 "rogchap.com/v8go"
)

type XrplJs struct {
	ctx *v8.Context
}

var startId uint64 = 0
var currentId *uint64 = &startId

func consumeId() uint64 {
	return atomic.AddUint64(currentId, 1)
}

func NewXrplJs() *XrplJs {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Panic().Msgf("Unable to get current dir: '%+v'", err)
	}

	// Check if current dir is this folder or project root folder (for xrp test to work)
	matches, _ := regexp.MatchString(".*internal/signer/local/xrp", currentDir)
	if matches {
		currentDir = filepath.Join(currentDir, "../../../..", "./external/xrpl.js")
	}

	// Check if current dir is this folder or project root folder (for test to work)
	matches, _ = regexp.MatchString(".*external/xrpl.js", currentDir)
	if !matches {
		currentDir = filepath.Join(currentDir, "./external/xrpl.js")
	}

	buff, err := os.ReadFile(filepath.Join(currentDir, "./xrpl.js"))
	if err != nil {
		log.Panic().Msgf("Error reading xrpl.js file: '%+v'", err)
	}

	ctx := v8.NewContext()

	_, err = ctx.RunScript(string(buff), "import-xrpl.js")
	if err != nil {
		log.Panic().Msgf("Error running script import-xrpl.js from xrpl.js: '%+v'", err)
	}

	return &XrplJs{
		ctx: ctx,
	}
}

func (xrplJs *XrplJs) EncodeForSigning(tx transaction.Transaction) string {
	id := "d" + strconv.FormatUint(consumeId(), 10)
	txStr, err := transaction.MarshalTransaction(tx)
	log.Info().Msgf("txStr: %s", txStr)
	if err != nil {
		log.Error().Msgf("Error marshaling transaction: '%+v'", err)
		return ""
	}

	_, err = xrplJs.ctx.RunScript("let "+id+" = xrpl.encodeForSigning("+txStr+")",
		"call-encodeForSigning.js")
	if err != nil {
		log.Error().Msgf("Error running xrpl.encodeForSigning script: '%+v'", err)
		return ""
	}

	val, err := xrplJs.ctx.RunScript(id, "get-value.js")
	if err != nil {
		log.Error().Msgf("Error running get-value.js of result script: '%+v'", err)
		return ""
	}

	return val.String()
}

func (xrplJs *XrplJs) EncodeForMultiSigning(tx transaction.Transaction, signerAddress string) string {
	id := "d" + strconv.FormatUint(consumeId(), 10)
	txStr, err := transaction.MarshalTransaction(tx)
	log.Info().Msgf("txStr: %s", txStr)
	if err != nil {
		log.Error().Msgf("Error marshaling transaction: '%+v'", err)
		return ""
	}

	_, err = xrplJs.ctx.RunScript("let "+id+" = xrpl.encodeForMultisigning("+txStr+", \""+signerAddress+"\")",
		"call-encodeForMultisigning.js")
	if err != nil {
		log.Error().Msgf("Error running xrpl.encodeForMultisigning script: '%+v'", err)
		return ""
	}

	val, err := xrplJs.ctx.RunScript(id, "get-value.js")
	if err != nil {
		log.Error().Msgf("Error running get-value.js of result script: '%+v'", err)
		return ""
	}

	return val.String()
}

func (xrplJs *XrplJs) Encode(jsonTx string) string {
	// Check jsonTx validity to prevent javascript injection
	var js json.RawMessage
	if json.Unmarshal([]byte(jsonTx), &js) != nil {
		log.Fatal().Msgf("")
		return ""
	}

	id := "c" + strconv.FormatUint(consumeId(), 10)
	_, err := xrplJs.ctx.RunScript("let "+id+" = xrpl.encode("+jsonTx+")",
		"call-encode.js")
	if err != nil {
		log.Error().Msgf("Error running xrpl.encode script: '%+v'", err)
		return ""
	}

	val, err := xrplJs.ctx.RunScript(id, "get-value.js")
	if err != nil {
		log.Error().Msgf("Error running get-value.js of encoded script: '%+v'", err)
		return ""
	}

	return val.String()
}

func (xrplJs *XrplJs) Decode(encodedTx string) string {
	id := "c" + strconv.FormatUint(consumeId(), 10)
	_, err := xrplJs.ctx.RunScript("let "+id+" = xrpl.decode('"+encodedTx+"')",
		"call-decode.js")
	if err != nil {
		log.Error().Msgf("Error running xrpl.decode script: '%+v'", err)
		return ""
	}

	val, err := xrplJs.ctx.RunScript(id, "get-value.js")
	if err != nil {
		log.Error().Msgf("Error running get-value.js of decoded script: '%+v'", err)
		return ""
	}

	marshalledTx, err := val.MarshalJSON()
	if err != nil {
		log.Error().Msgf("Error running MarshalJSON of encoded return value: '%+v'", err)
		return ""
	}

	return string(marshalledTx)
}
