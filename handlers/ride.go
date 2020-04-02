package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// RideHandler handles HTTP requests for rides.
type RideHandler struct {
	rideUsecase        usecases.RideUsecase
	maintenanceUsecase usecases.MaintenanceUsecase
}

// NewRideHandler returns a new RideHandler instance.
func NewRideHandler(rideUsecase usecases.RideUsecase, maintenanceUsecase usecases.MaintenanceUsecase) *RideHandler {
	return &RideHandler{
		rideUsecase,
		maintenanceUsecase,
	}
}

// Bind sets up the routes for the handler.
func (rh *RideHandler) Bind(e *echo.Echo) error {
	e.GET("/rides", rh.Fetch)
	e.POST("/rides", rh.Store)
	e.GET("/rides/:rideID", rh.GetByID)
	e.PUT("/rides/:rideID", rh.Update)
	e.DELETE("/rides/:rideID", rh.Delete)
	e.GET("/rides/:rideID/maintenance", rh.Maintenance)
	return nil
}

// Fetch fetches all rides.
func (rh *RideHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	rides, err := rh.rideUsecase.Fetch(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, rides, Indent)
}

// Store creates a new ride.
func (rh *RideHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	ride := &models.Ride{}

	err := c.Bind(ride)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = rh.rideUsecase.Store(ctx, ride)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, ride, Indent)
}

// GetByID gets a specific ride.
func (rh *RideHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("rideID")

	ride, err := rh.rideUsecase.GetByID(ctx, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusFound, ride, Indent)
}

// Update updates a specific ride.
func (rh *RideHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("rideID")

	ride := &models.Ride{}
	ride.ID = rideID

	err := c.Bind(ride)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = rh.rideUsecase.Update(ctx, ride)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, ride, Indent)
}

// Delete deletes a specific ride.
func (rh *RideHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("rideID")

	err := rh.rideUsecase.Delete(ctx, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, "", Indent)
}

// Maintenance fetches all maintenance jobs for the given ride.
func (rh *RideHandler) Maintenance(c echo.Context) error {
	return nil
}
