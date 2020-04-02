package impl

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errMaintenanceExists        = fmt.Errorf("maintenance job with the given ID already exists")
	errMaintenanceDoesNotExists = fmt.Errorf("maintenance job with he given ID does not exists")
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

// GetByID fetches a specific maintenance job from the repositories.
func (mu *MaintenanceUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.Maintenance, error) {
	maintenance, err := mu.maintenanceRepo.GetByID(ID)
	if err != nil {
		return nil, err
	}

	return maintenance, nil
}

// Fetch fetches all maintenance jobs from the repositories.
func (mu *MaintenanceUsecaseImpl) Fetch(ctx context.Context) ([]*models.Maintenance, error) {
	allMaintenance, err := mu.maintenanceRepo.Fetch()
	if err != nil {
		return nil, err
	}

	return allMaintenance, err
}

// FetchForRide is like fetch, but fetches all maintenance jobs for a specific ride.
func (mu *MaintenanceUsecaseImpl) FetchForRide(ctx context.Context, rideID string) ([]*models.Maintenance, error) {
	maintenance, err := mu.maintenanceRepo.FetchForRide(rideID)
	if err != nil {
		return nil, err
	}

	return maintenance, err
}

// Begin creates a new maintenance job in the repositories.
func (mu *MaintenanceUsecaseImpl) Begin(ctx context.Context, maintenance *models.Maintenance) error {
	_, err := mu.maintenanceRepo.GetByID(maintenance.ID)
	if err == nil {
		return errMaintenanceExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	maintenance.ID = uuid
	maintenance.End = sql.NullTime{}
	cleanMaintenance(maintenance)
	err = validateMaintenance(maintenance)
	if err != nil {
		return err
	}

	err = mu.maintenanceRepo.Store(maintenance)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a specific maintenance job in the repositories.
func (mu *MaintenanceUsecaseImpl) Update(ctx context.Context, maintenance *models.Maintenance) error {
	_, err := mu.maintenanceRepo.GetByID(maintenance.ID)
	if err != nil {
		return errMaintenanceDoesNotExists
	}

	cleanMaintenance(maintenance)
	err = validateMaintenance(maintenance)
	if err != nil {
		return err
	}

	maintenance.End = sql.NullTime{}
	err = mu.maintenanceRepo.Update(maintenance)
	if err != nil {
		return err
	}

	return nil
}

// Close closes an existing maintenance job in the repositories.
func (mu *MaintenanceUsecaseImpl) Close(ctx context.Context, ID string) (*models.Maintenance, error) {
	maintenance, err := mu.maintenanceRepo.GetByID(ID)
	if err != nil {
		return nil, errMaintenanceDoesNotExists
	}

	maintenance.End = sql.NullTime{Time: time.Now(), Valid: true}
	err = mu.maintenanceRepo.Update(maintenance)
	if err != nil {
		return nil, err
	}

	return maintenance, nil
}

// Delete deletes a specific maintenance job from the repositories.
func (mu *MaintenanceUsecaseImpl) Delete(ctx context.Context, ID string) error {
	_, err := mu.maintenanceRepo.GetByID(ID)
	if err != nil {
		return errMaintenanceDoesNotExists
	}

	err = mu.maintenanceRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

func cleanMaintenance(maintenance *models.Maintenance) {
	maintenance.ID = strings.TrimSpace(maintenance.ID)
	maintenance.Description = strings.TrimSpace(maintenance.Description)
}

func validateMaintenance(maintenance *models.Maintenance) error {

	if len(maintenance.ID) <= 0 {
		return fmt.Errorf("validateMaintenance: ID must be non-empty")
	}

	if len(maintenance.Description) <= 0 {
		return fmt.Errorf("validateMaintenance: description must be non-empty")
	}

	if maintenance.Cost <= 0 {
		return fmt.Errorf("validateMaintenance: cost must be positive")
	}

	if maintenance.End.Valid && maintenance.Start.After(maintenance.End.Time) {
		return fmt.Errorf("validateMaintenance: start date must be before end date")
	}

	return nil
}
