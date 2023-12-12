package evm

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type evmEncodingFixtures struct {
	tx      *types.Transaction
	encoded string
}

func getFixtures() []evmEncodingFixtures {
	fixtures := []evmEncodingFixtures{}

	address := common.HexToAddress("0x4DBeE27B94c970B6A7916628236ad6D9369a4518")
	fixtures = append(fixtures, evmEncodingFixtures{
		tx:      types.NewTransaction(1, address, nil, 10000000, nil, nil),
		encoded: "e0018083989680944dbee27b94c970b6a7916628236ad6d9369a45188080808080",
	})

	amount := big.NewInt(50000000000)
	fixtures = append(fixtures, evmEncodingFixtures{
		tx:      types.NewTransaction(2, address, amount, 10000000, nil, nil),
		encoded: "e5028083989680944dbee27b94c970b6a7916628236ad6d9369a4518850ba43b740080808080",
	})

	address2 := common.HexToAddress("0x4C5033DB823538d398e84Bf65fAdEbA0b4d71599")
	gasPrice := big.NewInt(100)
	fixtures = append(fixtures, evmEncodingFixtures{
		tx:      types.NewTransaction(3, address2, amount, 10000000, gasPrice, nil),
		encoded: "e5036483989680944c5033db823538d398e84bf65fadeba0b4d71599850ba43b740080808080",
	})

	fixtures = append(fixtures, evmEncodingFixtures{
		tx:      types.NewTransaction(4, address2, nil, 10000000, gasPrice, nil),
		encoded: "e0046483989680944c5033db823538d398e84bf65fadeba0b4d715998080808080",
	})

	parsedAbi, _ := abi.JSON(strings.NewReader(BridgeMetaData.ABI))
	address3 := common.HexToAddress("0xc2cD370bAdC28A01682394E8072824c1D7300D96")
	bridge := XChainTypesBridgeConfig{
		LockingChainDoor: address2,
		LockingChainIssue: XChainTypesBridgeChainIssue{
			Currency: "XRP",
			Issuer:   zeroAddress,
		},
		IssuingChainDoor: address3,
		IssuingChainIssue: XChainTypesBridgeChainIssue{
			Currency: "XRP",
			Issuer:   zeroAddress,
		},
	}
	input, _ := parsedAbi.Pack("addClaimAttestation", bridge, big.NewInt(int64(1)), amount, address, address3)
	fixtures = append(fixtures, evmEncodingFixtures{
		tx:      types.NewTransaction(4, address2, nil, 10000000, gasPrice, input),
		encoded: "f90246046483989680944c5033db823538d398e84bf65fadeba0b4d7159980b90224192dd3cc00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000ba43b74000000000000000000000000004dbee27b94c970b6a7916628236ad6d9369a4518000000000000000000000000c2cd370badc28a01682394e8072824c1d7300d960000000000000000000000004c5033db823538d398e84bf65fadeba0b4d715990000000000000000000000000000000000000000000000000000000000000080000000000000000000000000c2cd370badc28a01682394e8072824c1d7300d96000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000358525000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000035852500000000000000000000000000000000000000000000000000000000000808080",
	})

	return fixtures
}

func TestEvm_encodeTransaction(t *testing.T) {
	for _, fixture := range getFixtures() {
		encoded := encodeTransaction(fixture.tx)
		if encoded != fixture.encoded {
			t.Errorf("expected %+v got %+v", fixture.encoded, encoded)
		}
	}
}

func TestEvm_decodeTransaction(t *testing.T) {
	tx := "z1-.1"
	decoded := decodeTransaction(tx)
	if decoded != nil {
		t.Errorf("expected %+v got %+v", nil, decoded)
	}

	tx = "e0012180839896809400000000000000000000000000000000000000008080808080"
	decoded2 := decodeTransaction(tx)
	if decoded2 != nil {
		t.Errorf("expected %+v got %+v", nil, decoded2)
	}

	for _, fixture := range getFixtures() {
		decoded3 := decodeTransaction(fixture.encoded)
		if decoded3.Nonce() != fixture.tx.Nonce() {
			t.Errorf("nonce expected %+v got %+v", fixture.tx.Nonce(), decoded3.Nonce())
		}
		if decoded3.To().String() != fixture.tx.To().String() {
			t.Errorf("to expected %+v got %+v", fixture.tx.To().String(), decoded3.To().String())
		}
		if decoded3.Gas() != fixture.tx.Gas() {
			t.Errorf("to expected %+v got %+v", fixture.tx.Gas(), decoded3.Gas())
		}
		if decoded3.Value().String() != fixture.tx.Value().String() {
			t.Errorf("to expected %+v got %+v", fixture.tx.Value().String(), decoded3.Value().String())
		}
		if decoded3.GasPrice().String() != fixture.tx.GasPrice().String() {
			t.Errorf("to expected %+v got %+v", fixture.tx.GasPrice().String(), decoded3.GasPrice().String())
		}
		if string(decoded3.Data()) != string(fixture.tx.Data()) {
			t.Errorf("to expected %+v got %+v", string(fixture.tx.Data()), string(decoded3.Data()))
		}
	}
}
