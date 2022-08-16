package model

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
)

type UUID uuid.UUID

// MarshalText implements encoding.TextMarshaler.
func (u UUID) MarshalText() ([]byte, error) {
	if uuid.UUID(u) == uuid.Nil {
		return []byte(""), nil
	}

	var js [36]byte

	encodeHex(js[:], u)

	return js[:], nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *UUID) UnmarshalText(data []byte) error {
	id, err := ParseBytes(data)
	if err == nil {
		*u = id
	}
	return err
}

// func byteIsZero(s []byte) bool {
// 	for _, v := range s {
// 		if v != 0 {
// 			return false
// 		}
// 	}
// 	return true
// }

func encodeHex(dst []byte, uuid UUID) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// ParseBytes is like Parse, except it parses a byte slice instead of a string.
func ParseBytes(b []byte) (UUID, error) {
	var uuid UUID

	switch len(b) {
	case 36: // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9: // urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		if !bytes.Equal(bytes.ToLower(b[:9]), []byte("urn:uuid:")) {
			return uuid, fmt.Errorf("invalid urn prefix: %q", b[:9]) //nolint:goerr113
		}
		b = b[9:]
	case 36 + 2: // {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
		b = b[1:]
	case 32: // xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		var ok bool
		for i := 0; i < 32; i += 2 {
			uuid[i/2], ok = xtob(b[i], b[i+1])
			if !ok {
				return uuid, fmt.Errorf("invalid UUID format")
			}
		}
		return uuid, nil
	default:
		return uuid, fmt.Errorf("invalid UUID length: %d", len(b))
	}
	// s is now at least 36 bytes long
	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return uuid, fmt.Errorf("invalid UUID format")
	}
	for i, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		v, ok := xtob(b[x], b[x+1])
		if !ok {
			return uuid, fmt.Errorf("invalid UUID format")
		}
		uuid[i] = v
	}
	return uuid, nil
}

// xvalues returns the value of a byte as a hexadecimal digit or 255.
var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}

// xtob converts hex characters x1 and x2 into a byte.
func xtob(x1, x2 byte) (byte, bool) {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return (b1 << 4) | b2, b1 != 255 && b2 != 255
}
