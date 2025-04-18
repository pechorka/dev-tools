package uuid

import (
	"encoding/hex"
	"io"
	"time"

	"github.com/pechorka/gostdlib/pkg/errs"
)

// UUID represents a Universal Unique Identifier (UUID)
type UUID [16]byte

// String returns the string representation of the UUID
func (u UUID) String() string {
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], u[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], u[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], u[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], u[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], u[10:])
	return string(buf)
}

// NewV4 generates a new random UUID v4 using the provided reader
func NewV4(r io.Reader) (UUID, error) {
	var uuid UUID
	_, err := io.ReadFull(r, uuid[:])
	if err != nil {
		return UUID{}, errs.Wrap(err, "failed to read random bytes")
	}

	// Set version 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant to RFC4122
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

// MustV4 is a helper that wraps a call to NewV4 and panics if the error is non-nil
func MustV4(r io.Reader) UUID {
	uuid, err := NewV4(r)
	if err != nil {
		panic(err)
	}
	return uuid
}

// NewV4String generates a new random UUID v4 using the provided reader and returns it as a string
func NewV4String(r io.Reader) (string, error) {
	uuid, err := NewV4(r)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func MustV4String(r io.Reader) string {
	return MustV4(r).String()
}

// NewV7 generates a new UUID v7 (time-ordered) using the provided reader
func NewV7(r io.Reader) (UUID, error) {
	uuid, err := NewV4(r)
	if err != nil {
		return UUID{}, err
	}
	return toV7(uuid), nil
}

func MustV7(r io.Reader) UUID {
	uuid, err := NewV7(r)
	if err != nil {
		panic(err)
	}
	return uuid
}

func NewV7String(r io.Reader) (string, error) {
	uuid, err := NewV7(r)
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func MustV7String(r io.Reader) string {
	return MustV7(r).String()
}

func toV7(uuid UUID) UUID {
	now := time.Now().UnixMilli()

	// Set timestamp (first 48 bits)
	uuid[0] = byte(now >> 40)
	uuid[1] = byte(now >> 32)
	uuid[2] = byte(now >> 24)
	uuid[3] = byte(now >> 16)
	uuid[4] = byte(now >> 8)
	uuid[5] = byte(now)

	uuid[6] = (uuid[6] & 0x0f) | 0x70 // Set version 7
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Set variant to RFC4122

	return uuid
}
