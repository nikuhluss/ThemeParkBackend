package usecases

import (
	"context"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// UserUsecase is the usecase for interacting with users.
type UserUsecase interface {
	GetByID(context.Context, string) (*models.User, error)
	Fetch(context.Context) ([]*models.User, error)
	FetchCustomers(context.Context) ([]*models.User, error)
	FetchEmployees(context.Context) ([]*models.User, error)
	Store(context.Context, *models.User) error
	Update(context.Context, *models.User) error
}