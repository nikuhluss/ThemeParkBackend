package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// TicketHandler handles HTTP requests for tickets.
type TicketHandler struct {
	ticketUsecase usecases.TicketUsecase
}

// NewTicketHandler returns a new TicketHandler instance.
func NewTicketHandler(ticketUsecase usecases.TicketUsecase) *TicketHandler {
	return &TicketHandler{
		ticketUsecase,
	}
}

// Bind sets up the routes for the handler.
func (th *TicketHandler) Bind(e *echo.Echo) error {
	e.GET("/tickets", th.Fetch)
	e.POST("/tickets", th.Store)
	e.GET("/tickets/:rideID", th.GetByID)
	e.PUT("/tickets/:rideID", th.Update)
	e.DELETE("/tickets/:rideID", th.Delete)

	e.GET("/scans", th.FetchScans)
	e.POST("/scans/:ticketID/on/:rideID", th.StoreScan)

	e.GET("/users/:userID/tickets", th.FetchForUser)
	e.GET("/rides/:rideID/scans", th.FetchScansForRide)
	e.GET("/users/:userID/scans", th.FetchScansForUser)
	return nil
}

// Fetch fetches all tickets.
func (th *TicketHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	tickets, err := th.ticketUsecase.Fetch(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, tickets, Indent)
}

// FetchForUser fetches all tickets for the given user.
func (th *TicketHandler) FetchForUser(c echo.Context) error {
	ctx := c.Request().Context()

	userID := c.Param("userID")
	tickets, err := th.ticketUsecase.FetchForUser(ctx, userID)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, tickets, Indent)
}

// Store creates a new ticket.
func (th *TicketHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	ticket := &models.Ticket{}

	err := c.Bind(ticket)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = th.ticketUsecase.Store(ctx, ticket)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, ticket, Indent)
}

// GetByID gets a specific ticket.
func (th *TicketHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	ticketID := c.Param("ticketID")

	ticket, err := th.ticketUsecase.GetByID(ctx, ticketID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusFound, ticket, Indent)
}

// Update updates a specific ticket.
func (th *TicketHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	ticketID := c.Param("ticketID")

	ticket := &models.Ticket{}
	ticket.ID = ticketID

	err := c.Bind(ticket)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = th.ticketUsecase.Update(ctx, ticket)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, ticket, Indent)
}

// Delete deletes a specific ticket.
func (th *TicketHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	ticketID := c.Param("ticketID")

	err := th.ticketUsecase.Delete(ctx, ticketID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, "", Indent)
}

// FetchScans fetches all ticket scans.
func (th *TicketHandler) FetchScans(c echo.Context) error {
	ctx := c.Request().Context()

	scans, err := th.ticketUsecase.FetchScans(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, scans, Indent)
}

// StoreScan creates a new ticket scan.
func (th *TicketHandler) StoreScan(c echo.Context) error {
	ctx := c.Request().Context()
	ticketID := c.Param("ticketID")
	rideID := c.Param("rideID")

	scan, err := th.ticketUsecase.ScanTicket(ctx, ticketID, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, scan, Indent)
}

// FetchScansForRide fetches all scans for the given ride.
func (th *TicketHandler) FetchScansForRide(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("rideID")

	scans, err := th.ticketUsecase.FetchScansForRide(ctx, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, scans, Indent)
}

// FetchScansForUser fetches all scans for the given user.
func (th *TicketHandler) FetchScansForUser(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	scans, err := th.ticketUsecase.FetchScansForUser(ctx, userID)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, scans, Indent)
}
