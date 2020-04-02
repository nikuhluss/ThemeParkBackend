package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/handlers"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories/postgres"
	usecases "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases/impl"
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

	dbconfig := testutil.NewDatabaseConnectionConfig()
	db, err := testutil.NewDatabaseConnection(dbconfig)
	if err != nil {
		return err
	}

	// repos

	// userRepository := repos.NewUserRepository(db)
	rideRepository := repos.NewRideRepository(db)
	pictureRepository := repos.NewPictureRepository(db)
	reviewRepository := repos.NewReviewRepository(db)
	maintenanceRepository := repos.NewMaintenanceRepository(db)

	// usecases

	timeout := time.Second * 2
	rideUsecase := usecases.NewRideUsecaseImpl(rideRepository, pictureRepository, reviewRepository, timeout)
	maintenanceUsecase := usecases.NewMaintenanceUsecaseImpl(maintenanceRepository, timeout)

	// handlers

	rideHandler := handlers.NewRideHandler(rideUsecase, maintenanceUsecase)
	err = rideHandler.Bind(e)
	if err != nil {
		return err
	}

	return e.Start(address)
}
