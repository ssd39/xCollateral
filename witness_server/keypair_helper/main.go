package main

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	pk, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("PrivateKey:", hex.EncodeToString(crypto.FromECDSA(pk)))
	fmt.Println("PublicKey:", hex.EncodeToString(crypto.CompressPubkey(&pk.PublicKey)))
}
