package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/usecases"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	userUsecase        usecases.UserUsecase
}

// NewUserHandler returns a new UserHandler instance.
func NewUserHandler(userUsecase usecases.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase,
	}
}

// Bind sets up the routes for the handler.
func (uh *UserHandler) Bind(e *echo.Echo) error {
	e.GET("/users", uh.Fetch)
	e.GET("/users/customers", uh.FetchCustomers)
	e.GET("/users/employees", uh.FetchEmployees)
	e.GET("/users/:userID", uh.GetByID)
	e.POST("/users", uh.Store)
	e.PUT("/users/:userID", uh.Update)
	return nil
}

// Fetch fetches all users.
func (uh *UserHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := uh.userUsecase.Fetch(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, users, Indent)
}

// FetchCustomers fetches all customers.
func (uh *UserHandler) FetchCustomers(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := uh.userUsecase.FetchCustomers(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, users, Indent)
}

// FetchEmployees fetches all employees.
func (uh *UserHandler) FetchEmployees(c echo.Context) error {
	ctx := c.Request().Context()

	users, err := uh.userUsecase.FetchEmployees(ctx)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, users, Indent)
}

// GetByID gets a specific user.
func (uh *UserHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	user, err := uh.userUsecase.GetByID(ctx, userID)
	if err != nil {
		return c.JSONPretty(http.StatusNotFound, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusFound, user, Indent)
}

// Store creates a new user.
func (uh *UserHandler) Store(c echo.Context) error {
	ctx := c.Request().Context()

	user := &models.User{}

	err := c.Bind(user)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = uh.userUsecase.Store(ctx, user)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusCreated, user, Indent)
}

// Update updates a specific user.
func (uh *UserHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	user := &models.User{}
	user.ID = userID

	err := c.Bind(user)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	err = uh.userUsecase.Update(ctx, user)
	if err != nil {
		return c.JSONPretty(http.StatusBadRequest, ResponseError{err.Error()}, Indent)
	}

	return c.JSONPretty(http.StatusOK, user, Indent)
}
