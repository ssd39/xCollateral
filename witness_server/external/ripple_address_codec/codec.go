package ripple_address_codec

import (
	"bytes"
	"errors"
	"github.com/mr-tron/base58"
)

type Codec struct {
	alphabet *base58.Alphabet
}

var codec = Codec{alphabet: base58.NewAlphabet("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")}

func (c *Codec) Encode(buffer []byte, opts struct {
	versions       []byte
	expectedLength *int
}) (string, error) {
	if opts.expectedLength != nil && *opts.expectedLength != len(buffer) {
		return "", errors.New("bytes.length does not match expectedLength")
	}
	return c.encodeChecked(append(opts.versions, buffer...)), nil
}

func (c *Codec) encodeChecked(buffer []byte) string {
	check := HashSha256(HashSha256(buffer))[:4]
	return c.encodeRaw(append(buffer, check...))
}

func (c *Codec) encodeRaw(buffer []byte) string {
	return base58.EncodeAlphabet(buffer, c.alphabet)
}

func (c *Codec) decodeChecked(base58string string) ([]byte, error) {
	decoded, err := base58.DecodeAlphabet(base58string, c.alphabet)
	if err != nil {
		return nil, err
	}
	if len(decoded) < 5 {
		return nil, errors.New("invalid_input_size: decoded data must have length >= 5")
	}
	if !verifyCheckSum(decoded) {
		return nil, errors.New("checksum_invalid")
	}
	return decoded[:len(decoded)-4], nil
}

func verifyCheckSum(buffer []byte) bool {
	computed := HashSha256(HashSha256(buffer[:len(buffer)-4]))[:4]
	checksum := buffer[len(buffer)-4:]
	return bytes.Equal(computed, checksum)
}
