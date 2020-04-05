package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	// middleware

	e.Use(middleware.Logger())

	// database

	dbconfig := testutil.NewDatabaseConnectionConfig()
	db, err := testutil.NewDatabaseConnection(dbconfig)
	if err != nil {
		return err
	}

	// repos

	userRepo := repos.NewUserRepository(db)
	rideRepo := repos.NewRideRepository(db)
	pictureRepo := repos.NewPictureRepository(db)
	reviewRepo := repos.NewReviewRepository(db)
	maintenanceRepo := repos.NewMaintenanceRepository(db)
	ticketRepo := repos.NewTicketRepository(db)

	// usecases

	timeout := time.Second * 2
	rideUsecase := usecases.NewRideUsecaseImpl(rideRepo, pictureRepo, reviewRepo, timeout)
	maintenanceUsecase := usecases.NewMaintenanceUsecaseImpl(maintenanceRepo, timeout)
	ticketUsecase := usecases.NewTicketUsecaseImpl(ticketRepo, rideRepo, userRepo)

	// handlers

	rideHandler := handlers.NewRideHandler(rideUsecase, maintenanceUsecase)
	err = rideHandler.Bind(e)
	if err != nil {
		return err
	}

	maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceUsecase)
	err = maintenanceHandler.Bind(e)
	if err != nil {
		return err
	}

	ticketHandler := handlers.NewTicketHandler(ticketUsecase)
	err = ticketHandler.Bind(e)
	if err != nil {
		return err
	}

	return e.Start(address)
}
