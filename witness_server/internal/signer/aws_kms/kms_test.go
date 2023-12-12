package aws_kms

import (
	"encoding/hex"
	"testing"
)

func TestKms_AddressDerivation(t *testing.T) {
	evmAddress := "0x96329a50d10a3f69311e4f4e108672926c51c474"
	evmAddressInXrpl := "rN6wy5CiaM8Ng5JcM7YFT4vDKapYyk3FLk"
	xrpAddress := "rpSspP5yYyomcSrgsohyKMCnu5oJsTMkYP"
	xrpAddressInEvm := "0x0fb436e1514eb41310c50ac60d675c2542851715"

	evmAddressDerived := EvmAddressToXrplAccount(evmAddress)
	if evmAddressDerived != evmAddressInXrpl {
		t.Errorf("Invalid EvmAddressToXrplAccount %+v expected %+v\n", evmAddressDerived, evmAddressInXrpl)
	}

	evmAddressInEvmAgain := XrplAccountToEvmAddress(evmAddressDerived)
	if evmAddressInEvmAgain != evmAddress {
		t.Errorf("Invalid XrplAccountToEvmAddress %+v expected %+v\n", evmAddressInEvmAgain, evmAddress)
	}

	xrpAddressDerived := XrplAccountToEvmAddress(xrpAddress)
	if xrpAddressDerived != xrpAddressInEvm {
		t.Errorf("Invalid EvmAddressToXrplAccount %+v expected %+v\n", xrpAddressDerived, xrpAddressInEvm)
	}

	xrpAddressDerivedAgain := EvmAddressToXrplAccount(xrpAddressDerived)
	if xrpAddressDerivedAgain != xrpAddress {
		t.Errorf("Invalid EvmAddressToXrplAccount %+v expected %+v\n", xrpAddressDerivedAgain, xrpAddress)
	}
}

func TestKms_HashSha256(t *testing.T) {
	textToHash := "text to sha256 hash"
	expectedHash := "58c9296d6f1ac5b0643b8da4b56ca932cde7767ec2c426482594b6790ab49385"

	hashed := HashSha256([]byte(textToHash))
	hashedStr := hex.EncodeToString(hashed)
	if expectedHash != hashedStr {
		t.Errorf("Invalid HashSha256 %+v expected %+v\n", hashedStr, expectedHash)
	}
}

func TestKms_HashSha512(t *testing.T) {
	textToHash := "text to sha512 hash"
	expectedHash := "f56f3cb4c5da37dbaed1ce245ddeb4dc65e7ba5d7d2499323d55e4e80fd0fe88e729dda8373cd6abcac5037eb0154813a15cbd6e9fc4eedaab3dd733700afa5b"

	hashed := HashSha512([]byte(textToHash))
	hashedStr := hex.EncodeToString(hashed)
	if expectedHash != hashedStr {
		t.Errorf("Invalid HashSha512 %+v expected %+v\n", hashedStr, expectedHash)
	}
}

func TestKms_HashRipemd160(t *testing.T) {
	textToHash := "text to ripemd160 hash"
	expectedHash := "ffe64e297d5c812da40c85c93a62b09b29453abd"

	hashed := HashRipemd160([]byte(textToHash))
	hashedStr := hex.EncodeToString(hashed)
	if expectedHash != hashedStr {
		t.Errorf("Invalid HashRipemd160 %+v expected %+v\n", hashedStr, expectedHash)
	}
}

func TestKms_GetEthereumSignature(t *testing.T) {
	pubKeyHex := "041df461959a0973ca929c85db78d2370eb612d2671a9b7af4922b0519ebddf58c959b8e96175d8d8a675265b3852434b76887727e33e216211f308f07e6f937d6"
	pubKeyBytes, _ := hex.DecodeString(pubKeyHex)
	txHashHex := "98ef5b4f285e3bae9d09e1f0f9869ec7d7f7509a028b53b67606f6fda8d6c75b"
	txHashBytes, _ := hex.DecodeString(txHashHex)
	rHex := "ea6dd0ba4f206360c8c6e81e24996d315c9a7592024afea3509fe6fd909a1f7b"
	rBytes, _ := hex.DecodeString(rHex)
	sHex := "7f55d465f5f8531a3e615edc5b525ba0cc71acb7999748f770eeb6c6e8385b9a"
	sBytes, _ := hex.DecodeString(sHex)
	expectedSignatureHex := "ea6dd0ba4f206360c8c6e81e24996d315c9a7592024afea3509fe6fd909a1f7b7f55d465f5f8531a3e615edc5b525ba0cc71acb7999748f770eeb6c6e8385b9a00"

	// Should return nil if could not recover public key from signature
	invalid := GetEthereumSignature(pubKeyBytes[1:], txHashBytes[1:], rBytes, sBytes)
	if invalid != nil {
		t.Errorf("Invalid GetEthereumSignature %+v expected %+v\n", invalid, nil)
	}

	// Should return nil if could not reconstruct public key from signature
	invalid = GetEthereumSignature(pubKeyBytes[1:], txHashBytes, rBytes, sBytes)
	if invalid != nil {
		t.Errorf("Invalid GetEthereumSignature %+v expected %+v\n", invalid, nil)
	}

	// Should return return signature signature
	signatureBytes := GetEthereumSignature(pubKeyBytes, txHashBytes, rBytes, sBytes)
	signatureHex := hex.EncodeToString(signatureBytes)
	if expectedSignatureHex != signatureHex {
		t.Errorf("Invalid GetEthereumSignature %+v expected %+v\n", signatureHex, expectedSignatureHex)
	}
}

func TestKms_BigIntBytesToSignatureHex(t *testing.T) {
	rHex := "ea6dd0ba4f206360c8c6e81e24996d315c9a7592024afea3509fe6fd909a1f7b"
	rBytes, _ := hex.DecodeString(rHex)
	sHex := "7f55d465f5f8531a3e615edc5b525ba0cc71acb7999748f770eeb6c6e8385b9a"
	sBytes, _ := hex.DecodeString(sHex)
	expectedSignatureHex := "3045022100ea6dd0ba4f206360c8c6e81e24996d315c9a7592024afea3509fe6fd909a1f7b02207f55d465f5f8531a3e615edc5b525ba0cc71acb7999748f770eeb6c6e8385b9a"

	signatureHex := BigIntBytesToSignatureHex(rBytes, sBytes)
	if expectedSignatureHex != signatureHex {
		t.Errorf("Invalid GetEthereumSignature %+v expected %+v\n", signatureHex, expectedSignatureHex)
	}
}
