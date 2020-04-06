package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// EventHandler handles HTTP requests for events.
type EventHandler struct {
	eventUsecase usecases.EventUsecase
}

// NewEventHandler returns a new event handler instance.
func NewEventHandler(eventUsecase usecases.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase,
	}
}

// Bind sets up the routes for the handler.
func (eh *EventHandler) Bind(e *echo.Echo) error {
	e.GET("/events", eh.Fetch)
	e.POST("/events", eh.Store)
	e.GET("/events/:eventID", eh.GetByID)
	e.PUT("/events/:eventID", eh.Update)
	e.DELETE("/event/:eventID", eh.Delete)
	return nil
}

// GetByID gets a specific event.
func (eh *EventHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	eventID := c.Param("eventID")

	event, err := eh.eventUsecase.GetByID(ctx, eventID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, event, Indent)
}

// Fetch fetches all events.
func (eh *EventHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	day, _ := time.Parse(time.RFC3339, c.QueryParam("date"))
	var err error
	var event []*models.Event
	if len(c.QueryParam("date")) == 0 {
		event, err = eh.eventUsecase.Fetch(ctx)
	} else {
		event, err = eh.eventUsecase.FetchSince(ctx, day)
	}

	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, event, Indent)
}

// Store creates a new event.
func (eh *EventHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	event := &models.Event{}

	err := c.Bind(event)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = eh.eventUsecase.Store(ctx, event)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, event, Indent)
}

// Update updates a specific event.
func (eh *EventHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	eventID := c.Param("eventID")

	event := &models.Event{}
	event.ID = eventID

	err := c.Bind(event)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = eh.eventUsecase.Update(ctx, event)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, event, Indent)
}

// Delete deletes a specific event.
func (eh *EventHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	eventID := c.Param("eventID")

	err := eh.eventUsecase.Delete(ctx, eventID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, "", Indent)
}
