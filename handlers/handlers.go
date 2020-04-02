package handlers

import (
	"github.com/labstack/echo/v4"
)

// Indent is the constant used in JSONPretty responses.
const Indent = "    "

// Handler is an interface that lets you bind to an echo.Echo instance.
type Handler interface {
	Bind(e *echo.Echo) error
}

// ResponseError represents an error response with a (sometimes useful) message.
type ResponseError struct {
	Error string `json:"error"`
}
