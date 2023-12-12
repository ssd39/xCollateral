package utils

import (
	"fmt"
	"math"
	"math/big"
	"strings"
)

func Truncate(text string, width int) string {
	if width > len(text) {
		width = len(text)
	}
	r := []rune(text)
	trunc := r[:width]
	return string(trunc)
}

func IntToFloatPrec(number string, precision int) *big.Float {
	numberBI, _ := big.NewInt(0).SetString(number, 10)
	return new(big.Float).Quo(new(big.Float).SetInt(numberBI), big.NewFloat(math.Pow(10, float64(precision))))
}

func FloatToIntPrec(eth *big.Float, precision int) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(int64(math.Pow(10, float64(precision)))))
	fracStr := strings.Split(fmt.Sprintf("%.*f", precision, eth), ".")[1]
	fracStr += strings.Repeat("0", precision-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}
