package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// MaintenanceHandler handles HTTP requests for maintenance jobs.
type MaintenanceHandler struct {
	maintenanceUsecase           usecases.MaintenanceUsecase
}

func NewMaintenanceHandler(maintenanceUsecase usecases.MaintenanceUsecase) *MaintenanceHandler {
	return &MaintenanceHandler{
		maintenanceUsecase,
	}
}

// Bind sets up the routes for the handler.
func (mh *MaintenanceHandler) Bind(e *echo.Echo) error {
	e.GET("/maintenance/:maintenanceID", mh.GetByID)
	e.GET("/maintenance", mh.Fetch)
	e.GET("/maintenance/:rideID", mh.FetchForRide)
	e.POST("/maintenance/begin", mh.Store)
	e.PUT("/maintenance/:maintenanceID", mh.Update)
	e.POST("/maintenance/:maintenanceID/close", mh.Close)
	e.DELETE("/maintenance/:maintenanceID", mh.Delete)
	return nil
}

// GetByID gets a specific maintenance.
func (mh *MaintenanceHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	maintenanceID := c.Param("maintenanceID")

	maintenance, err := mh.maintenanceUsecase.GetByID(ctx, maintenanceID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusFound, maintenance, Indent)
}

// Fetch fetches all maintenance.
func (mh *MaintenanceHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	maintenance, err := mh.maintenanceUsecase.Fetch(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, maintenance, Indent)
}

// FetchForRide maintenance for a ride.
func (mh *MaintenanceHandler) FetchForRide(c echo.Context) error {
	ctx := c.Request().Context()
	rideID := c.Param("maintenanceID")

	maintenance, err := mh.maintenanceUsecase.FetchForRide(ctx, rideID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusFound, maintenance, Indent)
}

// Store creates a new maintenance.
func (mh *MaintenanceHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	maintenance := &models.Maintenance{}

	err := c.Bind(maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = mh.maintenanceUsecase.Begin(ctx, maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, maintenance, Indent)
}

// Update updates a specific maintenance.
func (mh *MaintenanceHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	maintenanceID := c.Param("maintenanceID")

	maintenance := &models.Maintenance{}
	maintenance.ID = maintenanceID

	err := c.Bind(maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = mh.maintenanceUsecase.Update(ctx, maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, maintenance, Indent)
}

// Close closes a specific maintenance.
func (mh *MaintenanceHandler) Close(c echo.Context) error {
	ctx := c.Request().Context()
	maintenanceID := c.Param("maintenanceID")

	maintenance := &models.Maintenance{}
	maintenance.ID = maintenanceID
	maintenance.MaintenanceType = "closed"

	err := c.Bind(maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	maintenance, err = mh.maintenanceUsecase.Close(ctx, maintenance.ID)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, maintenance, Indent)
}

// Delete deletes a specific maintenance.
func (mh *MaintenanceHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	maintenanceID := c.Param("maintenanceID")

	err := mh.maintenanceUsecase.Delete(ctx, maintenanceID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, "", Indent)
}