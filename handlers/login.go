package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	middlew "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/middleware"
)

// LoginHandler handles HTTP requests for log in.
type LoginHandler struct {
	keyAuth *middlew.KeyAuth
}

// Bind creates the required HTTP routes.
func (lh *LoginHandler) Bind(e *echo.Echo) {
	e.POST("/login", lh.Login)
}

// Login handles user login.
func (lh *LoginHandler) Login(c echo.Context) error {

	type Credentials struct {
		Login    string `json:"email"`
		Password string `json:"password"`
	}

	credentials := Credentials{}
	err := c.Bind(&credentials)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return nil
}
