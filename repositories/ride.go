package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// RideRepository defines the interface for working with rides.
type RideRepository interface {
	GetByID(ID string) (*models.Ride, error)
	Fetch() ([]*models.Ride, error)
	Store(*models.Ride) error
	Update(*models.Ride) error
	Delete(ID string) error
}
