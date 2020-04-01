package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// Indent is the indentation used in pretty JSON responses.
const Indent = "    "

// ErrorResponse represents a JSON response for an error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// notImplemented is a handler that does nothing but return an error.
func notImplemented(c echo.Context) error {
	errResponse := ErrorResponse{
		fmt.Sprintf("'%s': path not implemented", c.Path()),
	}
	return c.JSONPretty(http.StatusInternalServerError, errResponse, Indent)
}

// Start starts and HTTP server.
func Start(address string) error {

	e := echo.New()

	// auth routes

	e.POST(
		"/login", notImplemented)

	// ride routes

	e.GET(
		"/rides", notImplemented)

	e.POST(
		"/rides", notImplemented)

	e.GET(
		"/rides/:rideID", notImplemented)

	e.PATCH(
		"/rides/:rideID", notImplemented)

	e.DELETE(
		"/rides/:rideID", notImplemented)

	e.GET(
		"/rides/:rideID/maintenance", notImplemented)

	// maintenance routes

	e.GET(
		"/maintenance", notImplemented)

	e.POST(
		"/maintenance", notImplemented)

	e.GET(
		"/maintenance/:maintenanceID", notImplemented)

	e.PATCH(
		"/maintenance/:maintenanceID", notImplemented)

	e.POST(
		"/maintenance/:maintenanceID/close", notImplemented)

	e.DELETE(
		"/maintenance/:maintenanceID", notImplemented)

	// ticket routes

	e.GET(
		"/tickets", notImplemented)

	e.POST(
		"/tickets", notImplemented)

	e.GET(
		"/tickets/:ticketID", notImplemented)

	e.PATCH(
		"/tickets/:ticketID", notImplemented)

	e.POST(
		"/tickets/:ticketID/scan/:rideID", notImplemented)

	e.DELETE(
		"/tickets/:ticketID", notImplemented)

	return e.Start(address)
}
