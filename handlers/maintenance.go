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
	//e.PUT("/rides/:rideID", rh.Update)
	//e.DELETE("/rides/:rideID", rh.Delete)
	//e.GET("/rides/:rideID/maintenance", rh.FetchMaintenance)
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
func (rh *MaintenanceHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	maintenance := &models.Maintenance{}

	err := c.Bind(maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = rh.maintenanceUsecase.Begin(ctx, maintenance)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, maintenance, Indent)
}