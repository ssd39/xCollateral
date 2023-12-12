package ripple_address_codec

import (
	"testing"
)

func Test_IsValidXAddress(t *testing.T) {
	fixtures := []struct {
		Address string
		Valid   bool
	}{{
		"",
		false,
	}, {
		"rPLi4pWSjq98fJSPvu3EZmZSWxeq122wzp",
		false,
	}, {
		"rPLi4pWSjq98fJSPvu3EZmZSWxeq122wzz",
		false,
	}, {
		"rOPLi4pWSjq98fJSPvu3EZmZSWxeq122wz",
		false,
	}, {
		"rOPLi4pWSjq98fJSPvu3EZmZSWxeq1",
		false,
	}, {
		"X7YDPC4TJvjVxLc4QNDgCfaAocYVbWBE8jzpaKPZBy8mKDf",
		true,
	}, {
		"X7YDPC4TJvjVxLc4QNDgCfaAocYVbWBE8jzpaKPZBy8mKDz",
		false,
	}}

	for _, fixture := range fixtures {
		valid := IsValidXAddress(fixture.Address)
		if valid != fixture.Valid {
			t.Errorf("Invalid result for address %s expected: %+v got: %+v\n", fixture.Address, fixture.Valid, valid)
		}
	}
}

func newUint32(i uint32) *uint32 {
	return &i
}

func Test_XAddressToClassicAddress(t *testing.T) {
	fixtures := []struct {
		Address                string
		ExpectedClassicAddress string
		ExpectedTag            *uint32
	}{{
		"X7AcgcsBL6XDcUb289X4mJ8djcdyKaB5hJDWMArnXr61cqZ",
		"r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
		nil,
	}, {
		"X7AcgcsBL6XDcUb289X4mJ8djcdyKaGZMhc9YTE92ehJ2Fu",
		"r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
		newUint32(1),
	}, {
		"X7AcgcsBL6XDcUb289X4mJ8djcdyKaGo2K5VpXpmCqbV2gS",
		"r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
		newUint32(14),
	}, {
		"X7AcgcsBL6XDcUb289X4mJ8djcdyKaLFuhLRuNXPrDeJd9A",
		"r9cZA1mLK5R5Am25ArfXFmqgNwjZgnfk59",
		newUint32(11747),
	}, {
		"XVLhHMPHU98es4dbozjVtdWzVrDjtV8AqEL4xcZj5whKbmc",
		"rGWrZyQqhTp9Xu7G5Pkayo7bXjH4k4QYpf",
		newUint32(0),
	}}

	for _, fixture := range fixtures {
		classicAddress, tag := XAddressToClassicAddress(fixture.Address)
		if classicAddress == "" {
			t.Errorf("XAddressToClassicAddress erroed")
		}
		if classicAddress != fixture.ExpectedClassicAddress {
			t.Errorf("Invalid result for address %s expected: %+v got: %+v\n", fixture.Address, fixture.ExpectedClassicAddress, classicAddress)
		}
		if tag != nil && *tag != *fixture.ExpectedTag {
			t.Errorf("Invalid result for tag %s expected: %+v got: %+v\n", fixture.Address, *fixture.ExpectedTag, *tag)
		}
		if fixture.ExpectedTag == nil && tag != nil {
			t.Errorf("Invalid result for tag %s expected: %+v got: %+v\n", fixture.Address, fixture.ExpectedTag, tag)
		}
	}
}
