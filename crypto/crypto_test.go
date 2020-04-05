package crypto_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/crypto"
)

func TestHashFromPasswordSucceeds(t *testing.T) {

	tests := []string{
		"",
		"one",
		"the quick brown fox jumps over the lazy dog",
	}

	for _, tt := range tests {
		hash, err := crypto.HashFromPassword(tt)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.NotEmpty(t, hash)
	}
}

func TestCompareHashAndPasswordSucceeds(t *testing.T) {

	tests := []string{
		"",
		"one",
		"the quick brown fox jumps over the lazy dog",
	}

	for _, tt := range tests {
		hash, err := crypto.HashFromPassword(tt)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.True(t, crypto.CompareHashAndPassword(hash, tt))
		assert.False(t, crypto.CompareHashAndPassword(hash, "some-other-not-equal-password"))
	}
}
