package handlers

import (
	"net/http"
	"strings"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"

	"github.com/labstack/echo/v4"
	middlew "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/middleware"
)

// LoginHandler handles HTTP requests for log in.
type LoginHandler struct {
	keyAuth     *middlew.KeyAuth
	userUsecase usecases.UserUsecase
}

// NewLoginHandler returns a new LoginHandler instance.
func NewLoginHandler(keyAuth *middlew.KeyAuth, userUsecase usecases.UserUsecase) *LoginHandler {
	return &LoginHandler{keyAuth, userUsecase}
}

// Bind creates the required HTTP routes.
func (lh *LoginHandler) Bind(e *echo.Echo) error {
	e.POST("/login", lh.Login)
	return nil
}

// Login handles user login.
func (lh *LoginHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	type Credentials struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	credentials := Credentials{}
	err := c.Bind(&credentials)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	credentials.Login = strings.TrimSpace(credentials.Login)
	credentials.Password = strings.TrimSpace(credentials.Password)
	if len(credentials.Login) <= 0 || len(credentials.Password) <= 0 {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{"please provided both login and password"}, Indent)
	}

	user, err := lh.userUsecase.GetByEmail(ctx, credentials.Login)
	if err != nil {
		return c.JSONPretty(http.StatusUnauthorized, ResponseError{err.Error()}, Indent)
	}

	key, err := lh.keyAuth.Login(ctx, user, credentials.Password)
	if err != nil {
		return c.JSONPretty(http.StatusUnauthorized, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"userId": user.ID, "key": key.Encode()}, Indent)
}
