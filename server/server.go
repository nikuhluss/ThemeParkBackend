package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/handlers"
	middlew "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/middleware"
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
func Start(address string, testing bool) error {

	e := echo.New()

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
	eventRepo := repos.NewEventRepository(db)

	// usecases

	timeout := time.Second * 2
	userUsecase := usecases.NewUserUsecaseImpl(userRepo, timeout)
	rideUsecase := usecases.NewRideUsecaseImpl(rideRepo, pictureRepo, reviewRepo, timeout)
	reviewUsecase := usecases.NewReviewUsecaseImpl(reviewRepo, rideRepo, timeout)
	maintenanceUsecase := usecases.NewMaintenanceUsecaseImpl(maintenanceRepo, timeout)
	ticketUsecase := usecases.NewTicketUsecaseImpl(ticketRepo, rideRepo, userRepo)
	eventUsecase := usecases.NewEventUsecaseImpl(eventRepo, timeout)

	// middleware

	keyAuth := middlew.NewKeyAuth(userUsecase)
	keyAuthConfig := middleware.DefaultKeyAuthConfig
	keyAuthConfig.Validator = keyAuth.Validator
	keyAuthConfig.Skipper = func(c echo.Context) bool {
		if strings.HasPrefix(c.Path(), "/login") {
			return true
		}
		return false
	}

	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true

	e.Use(middleware.CORSWithConfig(corsConfig))
	if !testing {
		e.Use(middleware.KeyAuthWithConfig(keyAuthConfig))
	}
	e.Use(middleware.Logger())

	// handlers

	loginHandler := handlers.NewLoginHandler(keyAuth, userUsecase)
	err = loginHandler.Bind(e)
	if err != nil {
		return err
	}

	userHandler := handlers.NewUserHandler(userUsecase)
	err = userHandler.Bind(e)
	if err != nil {
		return err
	}

	rideHandler := handlers.NewRideHandler(rideUsecase, maintenanceUsecase)
	err = rideHandler.Bind(e)
	if err != nil {
		return err
	}

	reviewHandler := handlers.NewReviewHandler(reviewUsecase)
	err = reviewHandler.Bind(e)
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

	eventHandler := handlers.NewEventHandler(eventUsecase)
	err = eventHandler.Bind(e)
	if err != nil {
		return err
	}

	return e.Start(address)
}
