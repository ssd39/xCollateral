package ripple_address_codec

import (
	"bytes"
	"errors"
	"github.com/rs/zerolog/log"
)

func IsValidXAddress(xAddress string) bool {
	_, _, err := decodeXAddress(xAddress)
	return err == nil
}

func XAddressToClassicAddress(xAddress string) (string, *uint32) {
	accountId, tag, err := decodeXAddress(xAddress)
	if err != nil {
		log.Error().Msgf("Error running decodeXAddress '%+v'", err)
		return "", nil
	}
	classicAddress, err := encodeAccountId(accountId)
	if err != nil {
		log.Error().Msgf("Error running encodeAccountId '%+v'", err)
		return "", nil
	}
	return classicAddress, tag
}

func decodeXAddress(xAddress string) ([]byte, *uint32, error) {
	decoded, err := codec.decodeChecked(xAddress)
	if err != nil {
		return nil, nil, err
	}
	tag, err := tagFromBuffer(decoded)
	if err != nil {
		return nil, nil, err
	}
	return decoded[2:22], tag, nil

}

func tagFromBuffer(buffer []byte) (*uint32, error) {
	if len(buffer) < 22 {
		return nil, errors.New("invalid_input_size: tag buffer data must have length >= 22")
	}
	flag := buffer[22]
	if flag >= 2 {
		return nil, errors.New("unsupported x-address")
	}
	if flag == 1 {
		result := uint32(buffer[23]) + uint32(buffer[24])*0x100 + uint32(buffer[25])*0x10000 + uint32(buffer[26])*0x1000000
		return &result, nil
	}
	if flag != 0 {
		return nil, errors.New("flag must be zero to indicate no tag")
	}
	empty := make([]byte, 8)
	if !bytes.Equal(empty, buffer[23:23+8]) {
		return nil, errors.New("remaining bytes must be zero")
	}

	return nil, nil
}

func encodeAccountId(buffer []byte) (string, error) {
	length := 20
	return codec.Encode(buffer, struct {
		versions       []byte
		expectedLength *int
	}{versions: []byte{0}, expectedLength: &length})
}
