package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/crypto"
)

func TestKeyEncodeSucceeds(t *testing.T) {

	tests := []struct {
		login        string
		passwordHash string
	}{
		{"", ""},
		{"login0", "hash0"},
		{"login1", "hash1  "},
		{"login2", "hash2 --"},
	}

	for _, tt := range tests {
		key := crypto.Key{tt.login, tt.passwordHash}
		encoded := key.Encode()
		assert.NotEmpty(t, encoded)
	}
}

func TestKeyDecodeSucceeds(t *testing.T) {

	tests := []struct {
		login        string
		passwordHash string
	}{
		{"", ""},
		{"login0", "hash0"},
		{"login1", "hash1  "},
		{"login2", "hash2 --"},
	}

	for _, tt := range tests {
		key := crypto.Key{tt.login, tt.passwordHash}

		encoded := key.Encode()
		decodedKey, err := crypto.DecodeKey(encoded)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, key, decodedKey)
	}
}
