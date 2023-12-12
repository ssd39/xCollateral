package evm

import (
	"encoding/hex"
	"math/big"
	config "peersyst/bridge-witness-go/configs"
	"peersyst/bridge-witness-go/internal/signer/local/evm"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func createSafeInstance(address string) (*Safe, error) {
	client, err := ethclient.Dial("https://rpc-evm-sidechain.xrpl.org")
	if err != nil {
		return nil, err
	}
	instance, err := NewSafe(common.HexToAddress(address), client)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

func TestSafeSigner_SignTransactionHash(t *testing.T) {
	privateKey := "022bfeaa81eed7d52f500990cac50e8d3561a89795e9a1121d25d1299edd0e9c"
	localSigner := evm.NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})

	safe, err := createSafeInstance("")
	if err != nil {
		t.Fatalf("Error connecting to client %v", err)
	}

	signer := SafeSigner{*safe, localSigner}

	fixtures := []struct {
		hash   string
		expect string
	}{
		{
			hash:   "7f3adbd9a09026886b4f5bdca712b92126dfaf75572ebcc14abbe410fb67fb66",
			expect: "0x1e9ceae9b4d7f08c88d39ec469ef1875571ae98bb9cfed3e979409a0ff840e925779ec35f8befa4fa35b10b61b464647e2c9d63e51c6dd43b5cc7e5f8042fdd620",
		},
		{
			hash:   "e53c9f1c06ec67d6948b6e5370e2d97482a45903f659dc6efa5828fd072bc369",
			expect: "0xd4b996e9dfd5d23d7ee77be3400115823c74d34cc5156d8515709107642bb73511f8b2643c1c0aab26c0e8f53d5f5ac2139150eb4aecc82a601421a9f388661d1f",
		},
		{
			hash:   "3852c6d77d1507e85372bc2ee0a1e7aae1dc9b482eeadaf9db85b50c1a01f88d",
			expect: "0x3c75907fa5cb907ce3f87c5f5abc5b67038f2930ddd5a7d486ea041f524898884491f54974c03dc745df66527a47b6d108b10f659e0a8ab98e5e417b2a4666bf20",
		},
		{
			hash:   "d716755ca71e498ad4bff20357a09f8d9941e2891446efd84736f988771ebf30",
			expect: "0x1233aebe418ab25fc6a334312a2fbacaf3e369c8732dae51936ec8a7f66a1ce13c0f0ffc1adbecbcf44dd37078c257fba570c8db466f186f32a9c7406d5ef59c20",
		},
		{
			hash:   "a40c37e5c233545dae4444866bb7ecaa825b73801a6119b74b574895a26102b1",
			expect: "0xef140add4ca880799483244ff7ef33f7fa9656ed976cf8269edcba7fe7aee95e06af3fbb2c3063a8d0b531ac969b7258b719371a63c7eef23e23140cba7d3f8420",
		},
		{
			hash:   "11553769da6c743281c1dadde3550474b5b8361ed7dae015162bfa37a848b3c1",
			expect: "0x64803de62ca7367faded863c57a1086fbab9cfa0a4a5d13de3d04ccdeb6915726f1f8ca8447ec039113e5890c9a3fcf488c9dd2cb07390aac4176b0bcee058f91f",
		},
	}
	for _, fixture := range fixtures {

		signature, err := signer.SignTransactionHash(fixture.hash)
		if err != nil {
			t.Fatalf("Error signing transaction %v", err)
		}
		if fixture.expect != signature.data {
			t.Fatalf("Invalid signature - expected: %+v got: %+v", fixture.expect, signature.data)
		}
	}
}

func TestSafeSigner_GetTransactionHash(t *testing.T) {
	privateKey := "022bfeaa81eed7d52f500990cac50e8d3561a89795e9a1121d25d1299edd0e9c"
	localSigner := evm.NewEvmLocalSignerProvider(config.LocalSigner{PrivateKey: privateKey})
	safe, err := createSafeInstance("0x94b6ba740973abffde3334dec43db3c20272985a")
	if err != nil {
		t.Fatalf("Error connecting to client %v", err)
	}
	signer := SafeSigner{*safe, localSigner}

	data, _ := hex.DecodeString("4a07d6730000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a00000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c000000000000000000000000000000000000000000000000000000000000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c0000000000000000000000000000000000000000000000000000000000")
	safeTransaction := SafeTransaction{
		To:             common.HexToAddress("0xA01cA6C79682358c59fCDf4a2e0Ce125bC7a41D0"),
		Value:          big.NewInt(0),
		Data:           data,
		Operation:      0,
		SafeTxGas:      big.NewInt(0x1e8480),
		BaseGas:        big.NewInt(0x1e8480),
		GasPrice:       big.NewInt(7),
		GasToken:       common.HexToAddress("0x0000000000000000000000000000000000000000"),
		RefundReceiver: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Nonce:          big.NewInt(123),
	}

	hash, _ := signer.GetTransactionHash(safeTransaction)
	expectedHash := "c9499a7d33076e1b1e5597641d08ca3f771fba06c337fa7dc7d10a224ecdfabf"
	if hex.EncodeToString(hash[:]) != expectedHash {
		t.Fatalf("Invalid hash - expected: %+v got: %+v", expectedHash, hex.EncodeToString(hash[:]))
	}
}

func TestSafeSigner_UnmarshalTransaction(t *testing.T) {
	jsonTx := `{"to":"0xa67f7EE8179FA6B87349981C10ac74d0eB93DA3E","value":"0","data":"0x4a07d6730000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a00000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c000000000000000000000000000000000000000000000000000000000000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c0000000000000000000000000000000000000000000000000000000000","operation":0,"baseGas":"0x1e8480","gasPrice":"7","gasToken":"0x0000000000000000000000000000000000000000","refundReceiver":"0x0000000000000000000000000000000000000000","nonce":125,"safeTxGas":"0x1e8480"}`

	tx, err := parseSafeTransactionInJSON(jsonTx)
	if err != nil {
		t.Fatalf("Error parsing transaction %s", err)
	}

	dataStr := "4a07d6730000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000008000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a00000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c000000000000000000000000000000000000000000000000000000000000000000000000000000000094b6ba740973abffde3334dec43db3c20272985a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000034d554c0000000000000000000000000000000000000000000000000000000000"
	data, err := hex.DecodeString(dataStr)
	if err != nil {
		t.Fatalf("Error decoding data string %s", err)
	}

	safeTransaction := SafeTransaction{
		To:             common.HexToAddress("0xa67f7EE8179FA6B87349981C10ac74d0eB93DA3E"),
		Value:          big.NewInt(0),
		Data:           data,
		Operation:      0,
		SafeTxGas:      big.NewInt(0x1e8480),
		BaseGas:        big.NewInt(0x1e8480),
		GasPrice:       big.NewInt(7),
		GasToken:       common.HexToAddress("0x0000000000000000000000000000000000000000"),
		RefundReceiver: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		Nonce:          big.NewInt(125),
	}

	if tx.To.String() != safeTransaction.To.String() {
		t.Fatalf("Invalid to - expected: %+v got: %+v", tx.To.String(), safeTransaction.To.String())
	}
	if tx.Value.Cmp(safeTransaction.Value) != 0 {
		t.Fatalf("Invalid value - expected: %+v got: %+v", tx.Value, safeTransaction.Value)
	}
	if hex.EncodeToString(tx.Data) != dataStr {
		t.Fatalf("Invalid data - expected: %+v got: %+v", hex.EncodeToString(tx.Data), dataStr)
	}
	if tx.Operation != safeTransaction.Operation {
		t.Fatalf("Invalid operation - expected: %+v got: %+v", tx.Operation, safeTransaction.Operation)
	}
	if tx.SafeTxGas.Cmp(safeTransaction.SafeTxGas) != 0 {
		t.Fatalf("Invalid safeTxGas - expected: %+v got: %+v", tx.SafeTxGas, safeTransaction.SafeTxGas)
	}
	if tx.BaseGas.Cmp(safeTransaction.BaseGas) != 0 {
		t.Fatalf("Invalid baseGas - expected: %+v got: %+v", tx.BaseGas, safeTransaction.BaseGas)
	}
	if tx.GasPrice.Cmp(safeTransaction.GasPrice) != 0 {
		t.Fatalf("Invalid gasPrice - expected: %+v got: %+v", tx.GasPrice, safeTransaction.GasPrice)
	}
	if tx.GasToken.String() != safeTransaction.GasToken.String() {
		t.Fatalf("Invalid gasToken - expected: %+v got: %+v", tx.GasToken.String(), safeTransaction.GasToken.String())
	}
	if tx.RefundReceiver.String() != safeTransaction.RefundReceiver.String() {
		t.Fatalf("Invalid refundReceiver - expected: %+v got: %+v", tx.RefundReceiver.String(), safeTransaction.RefundReceiver.String())
	}
	if tx.Nonce.Cmp(safeTransaction.Nonce) != 0 {
		t.Fatalf("Invalid nonce - expected: %+v got: %+v", tx.Nonce, safeTransaction.Nonce)
	}
}
