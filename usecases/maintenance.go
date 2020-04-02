package usecases

import (
	"context"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// MaintenanceUsecase is the usecase for interacting with maintenance jobs.
type MaintenanceUsecase interface {
	GetByID(context.Context, string) (*models.Maintenance, error)
	Fetch(context.Context) ([]*models.Maintenance, error)
	FetchForRide(context.Context, string) ([]*models.Maintenance, error)
	Begin(context.Context, *models.Maintenance) error
	Update(context.Context, *models.Maintenance) error
	Close(context.Context, string) (*models.Maintenance, error)
	Delete(context.Context, string) error
}
