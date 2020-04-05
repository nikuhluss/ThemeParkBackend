package crypto

import (
	"encoding/base64"
	"encoding/json"
)

// Key holds an authorization key information.
type Key struct {
	Login        string `json:"login"`
	PasswordHash string `json:"passwordHash"`
}

// EmptyKey struct holds an empty key.
var EmptyKey = Key{}

// DecodeKey creates a new Key from the given encoded key.
func DecodeKey(encoded string) (Key, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return EmptyKey, err
	}

	key := Key{}
	err = json.Unmarshal(keyBytes, &key)
	if err != nil {
		return EmptyKey, err
	}

	return key, nil
}

// Encode encodes the key into a base64 string.
func (k *Key) Encode() string {
	keyBytes, _ := json.Marshal(k)
	return base64.StdEncoding.EncodeToString(keyBytes)
}
