package utils

import (
	"math/big"
	"testing"
)

func TestCache_Truncate(t *testing.T) {
	toTruncate := "value not pointer or nil"
	length := 10
	got := Truncate(toTruncate, length)
	if len(got) != length {
		t.Errorf("expected %+v got %+v", length, len(got))
	}
}

func TestCache_IntToFloatPrec(t *testing.T) {
	iniNumber := "10000"
	precision := 7
	expect := big.NewFloat(0.001)
	got := IntToFloatPrec(iniNumber, precision)
	if got.String() != expect.String() {
		t.Errorf("expected %+v got %+v", expect, got)
	}
}

func TestCache_FloatToIntPrec(t *testing.T) {
	iniNumber := big.NewFloat(0.001)
	precision := 7
	expect := big.NewInt(10000)
	got := FloatToIntPrec(iniNumber, precision)
	if got.String() != expect.String() {
		t.Errorf("expected %+v got %+v", expect, got)
	}
}
