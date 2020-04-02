package impl

import (
	"time"

	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

// MaintenanceUsecaseImpl implements the MaintenanceUsecase interface.
type MaintenanceUsecaseImpl struct {
}

// NewMaintenanceUsecaseImpl returns a new MaintenanceUsecaseImpl instance.
func NewMaintenanceUsecaseImpl(maintenanceRepository repos.MaintenanceRepository, timeout time.Duration) *MaintenanceUsecaseImpl {
	return nil
}
