package uuid

import (
	"crypto/rand"
)

// NewV4Crypto generates a new random UUID v4 using crypto/rand
func NewV4Crypto() (UUID, error) {
	return NewV4(rand.Reader)
}

// MustV4Crypto is a helper that wraps a call to NewV4Crypto and panics if the error is non-nil
func MustV4Crypto() UUID {
	return MustV4(rand.Reader)
}

// NewV4CryptoString generates a new random UUID v4 using crypto/rand and returns it as a string
func NewV4CryptoString() (string, error) {
	return NewV4String(rand.Reader)
}

func MustV4CryptoString() string {
	return MustV4String(rand.Reader)
}

// NewV7Crypto generates a new UUID v7 (time-ordered) using crypto/rand
func NewV7Crypto() (UUID, error) {
	return NewV7(rand.Reader)
}

func MustV7Crypto() UUID {
	return MustV7(rand.Reader)
}

func NewV7CryptoString() (string, error) {
	return NewV7String(rand.Reader)
}

func MustV7CryptoString() string {
	return MustV7String(rand.Reader)
}
