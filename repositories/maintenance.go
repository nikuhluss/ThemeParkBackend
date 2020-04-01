package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

type MaintenanceRepository interface {
	GetByID(ID string) (*models.Maintenance, error)
	Fetch() ([]*models.Maintenance, error)
	FetchByRideID(rideID string) ([]*models.Maintenance, error)
	Store(*models.Maintenance) error
	Update(*models.Maintenance) error
	Delete(ID string) error
	AvailableMaintenanceTypes() ([]string, error)
}
