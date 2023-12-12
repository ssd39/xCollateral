package oracle

import (
	"fmt"
	"math"
	"peersyst/bridge-witness-go/internal/chains"
	"peersyst/bridge-witness-go/internal/chains/xrp/xrpl"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func StartPriceOracle() {
	mainChainProvider := chains.GetMainChainProvider()
	sideChainProvider := chains.GetSideChainProvider()
	ticker := time.NewTicker(time.Second * 10)
	time.Sleep(time.Second)
	lastPrice := 0.0
	for range ticker.C {
		log.Info().Msgf("Fetching amm info.....")
		ammInfoResult, err := mainChainProvider.GetAmmInfo(&xrpl.AmmAsset{Currency: "XRP"}, &xrpl.AmmAsset{Currency: "TXT", Issuer: "rH9WvmWDk7CgcAPM9v8hAGmaVEQACfRa1Q"})
		if err != nil {
			log.Error().Msgf("Error while fetching amm info %v", err)
			continue
		}

		xrpValueDrops, err := strconv.ParseInt(ammInfoResult.Amm.Amount.Value, 10, 64)
		if err != nil {
			continue
		}
		iotValueF, err := strconv.ParseFloat(ammInfoResult.Amm.Amount2.Value, 64)
		if err != nil {
			continue
		}
		iotValueDrops := int64(iotValueF * 1000000.0)
		xrpValue := float64(xrpValueDrops) / 1000000.0
		log.Info().Msgf("Current price %v TXT/XRP", fmt.Sprintf("%.2f", iotValueF/xrpValue))
		newPrice := float64(xrpValue / iotValueF)
		if math.Abs(1.0-(newPrice/lastPrice))*100 >= 0.5 {
			log.Info().Msgf("Updating current price on evm")
			err := sideChainProvider.UpdateOracleData(xrpValueDrops, iotValueDrops)
			if err == nil {
				lastPrice = float64(xrpValue / iotValueF)
			}
		}
	}
}
