package usecases

import (
	"context"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// RideUsecase is the usecase for interacting with rides.
type RideUsecase interface {
	GetByID(context.Context, string) (*models.Ride, error)
	Fetch(context.Context) ([]*models.Ride, error)
	Store(context.Context, *models.Ride) error
	Update(context.Context, *models.Ride) error
	Delete(context.Context, string) error
}
