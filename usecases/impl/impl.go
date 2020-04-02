package impl

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID (V4) string.
func GenerateUUID() (string, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}
