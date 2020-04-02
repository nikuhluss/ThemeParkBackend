package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
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

	userRepository := repos.NewUserRepository(db)
	rideRepository := repos.NewRideRepository(db)
	pictureRepository := repos.NewPictureRepository(db)
	reviewRepository := repos.NewReviewRepository(db)
	maintenanceRepository := repos.NewMaintenanceRepository(db)

	timeout := time.Second * 2

	rideUsecase := usecases.NewRideUsecaseImpl(rideRepository, pictureRepository, reviewRepository, timeout)
	maintenanceUsecase := usecases.NewMaintenanceUsecaseImpl(maintenanceRepository, timeout)

	fmt.Println(userRepository, rideUsecase, maintenanceUsecase)

	return e.Start(address)
}
