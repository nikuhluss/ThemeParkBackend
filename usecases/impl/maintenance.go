package impl

import (
	"context"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

// MaintenanceUsecaseImpl implements the MaintenanceUsecase interface.
type MaintenanceUsecaseImpl struct {
	maintenanceRepo repos.MaintenanceRepository
	timeout         time.Duration
}

// NewMaintenanceUsecaseImpl returns a new MaintenanceUsecaseImpl instance.
func NewMaintenanceUsecaseImpl(maintenanceRepo repos.MaintenanceRepository, timeout time.Duration) *MaintenanceUsecaseImpl {
	return &MaintenanceUsecaseImpl{
		maintenanceRepo,
		timeout,
	}
}

func (mu *MaintenanceUsecaseImpl) GetByID(context.Context, string) (*models.Maintenance, error) {
	return nil, nil
}

func (mu *MaintenanceUsecaseImpl) Fetch(context.Context) ([]*models.Maintenance, error) {
	return nil, nil
}

func (mu *MaintenanceUsecaseImpl) FetchForRide(context.Context, string) ([]*models.Maintenance, error) {
	return nil, nil
}

func (mu *MaintenanceUsecaseImpl) Store(context.Context, *models.Maintenance) error {
	return nil
}

func (mu *MaintenanceUsecaseImpl) Update(context.Context, *models.Maintenance) error {
	return nil
}

func (mr *MaintenanceRepositoryImpl) Delete(context.Context, string) error {
	return nil
}
