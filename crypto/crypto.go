package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashFromPassword returns a hashed password from the given password.
func HashFromPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hashedPassword := string(hashedPasswordBytes)
	return hashedPassword, nil
}

// CompareHashAndPassword compares the given hashed password and plain password
// and returns true and no error if the are equivalent.
func CompareHashAndPassword(hashedPassword, plainPassword string) (bool, error) {
	hashedPasswordBytes := []byte(hashedPassword)
	plainPasswordBytes := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, plainPasswordBytes)
	if err != nil {
		return false, err
	}

	return true, nil
}
