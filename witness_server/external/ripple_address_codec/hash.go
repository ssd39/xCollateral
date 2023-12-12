package ripple_address_codec

import "crypto/sha256"

func HashSha256(bytes []byte) []byte {
	hashSha256 := sha256.New()
	hashSha256.Write(bytes)
	return hashSha256.Sum(nil)
}
