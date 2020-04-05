package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/crypto"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// KeyAuth struct helps with log in and log in validation.
type KeyAuth struct {
	userUsecase usecases.UserUsecase
	expirations map[crypto.Key]time.Time
}

// NewKeyAuth returns a new KeyAuth instance.
func NewKeyAuth(userUsecase usecases.UserUsecase) *KeyAuth {
	return &KeyAuth{
		userUsecase,
		make(map[crypto.Key]time.Time),
	}
}

// Login takes an user and its password and returns a crypto.Key if login was successful.
func (ka *KeyAuth) Login(ctx context.Context, userID, password string) (crypto.Key, error) {
	user, err := ka.userUsecase.GetByID(ctx, userID)
	if err != nil {
		return crypto.EmptyKey, fmt.Errorf("loginError: %s", err)
	}

	if !crypto.CompareHashAndPassword(user.PasswordHash, password) {
		return crypto.EmptyKey, fmt.Errorf("loginError: invalid username or password")
	}

	return crypto.Key{Login: user.ID, PasswordHash: user.PasswordHash}, nil
}

// Validator is a validator function for middleware.KeyAuth from echo.
func (ka *KeyAuth) Validator(encodedKey string, c echo.Context) (bool, error) {

	key, err := crypto.DecodeKey(encodedKey)
	if err != nil {
		return false, err
	}

	// check if key was used before and has not expired

	if exp, ok := ka.expirations[key]; ok && time.Now().Before(exp) {
		return true, nil
	}

	// key has not been used or has expired

	ctx := c.Request().Context()
	user, err := ka.userUsecase.GetByID(ctx, key.Login)
	if err != nil {
		return false, err
	}

	if key.PasswordHash != user.PasswordHash {
		return false, fmt.Errorf("keyAuthValidator: given key is not valid")
	}

	// key is valid, set up new expiration time

	ka.expirations[key] = time.Now().Add(time.Hour * 24)
	return true, nil
}
