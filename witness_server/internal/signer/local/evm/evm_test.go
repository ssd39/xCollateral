package evm

import (
	"encoding/hex"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/signer"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestEvm_SignTransaction(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	var data []byte
	ethTransaction := types.NewTransaction(
		1,
		common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d"),
		big.NewInt(0),
		uint64(21000),
		big.NewInt(1),
		data,
	)
	rawTxBytes, _ := rlp.EncodeToBytes(ethTransaction)
	transaction := hex.EncodeToString(rawTxBytes)
	expect := "f85e0101825208944592d8f8d7b001e72cb26a73e4fa1806a51ac79d808026a0316235aabac9298a1e14ba3cec3a1e673098288bf166ec2bbaa8fa29880b94669f953481d9746dc2c33829945a24b431f8cf45e965a14993aa33dafcee92a30b"

	localSigner := NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})
	chainId := big.NewInt(1)
	opts := signer.SignEvmTransactionOpts{ChainId: chainId}
	signedTransaction := localSigner.SignTransaction(transaction, &opts)
	if expect != signedTransaction {
		t.Fatalf("Invalid transaction signed - expected: %+v got: %+v", expect, signedTransaction)
	}
}

func TestEvm_GetAddress(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	expect := "0x96216849c49358B10257cb55b28eA603c874b05E"
	localSigner := NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	if expect != localSigner.GetAddress() {
		t.Fatalf("Invalid address - expected: %+v got: %+v", expect, localSigner.GetAddress())
	}
}

func TestEvm_GetPublicKey(t *testing.T) {
	privateKey := "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	expect := "039a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fd"
	localSigner := NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	if expect != localSigner.GetPublicKey() {
		t.Fatalf("Invalid public key - expected: %+v got: %+v", expect, localSigner.GetPublicKey())
	}
}

func TestEvm_SignMessage(t *testing.T) {
	privateKey := "022bfeaa81eed7d52f500990cac50e8d3561a89795e9a1121d25d1299edd0e9c"
	fixtures := []struct {
		message string
		expect  string
	}{
		{
			message: hex.EncodeToString([]byte("hello world")),
			expect:  "d831e5984855e6fbfd8ee79eb186d99e0f73ee7d52a96ea27b7c12de060313b875825a1abd91329b0e2ff20bc20962ba7dc09eb0e438502f03960721a1b846211c",
		},
		{
			message: "7f3adbd9a09026886b4f5bdca712b92126dfaf75572ebcc14abbe410fb67fb66",
			expect:  "1e9ceae9b4d7f08c88d39ec469ef1875571ae98bb9cfed3e979409a0ff840e925779ec35f8befa4fa35b10b61b464647e2c9d63e51c6dd43b5cc7e5f8042fdd61c",
		},
	}
	for _, fixture := range fixtures {
		localSigner := NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

		signature := localSigner.SignMessage(fixture.message)
		if fixture.expect != signature {
			t.Fatalf("Invalid signature - expected: %+v got: %+v", fixture.expect, signature)
		}
	}
}
