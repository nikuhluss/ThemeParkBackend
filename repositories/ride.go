package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// RideRepository defines the interface for working with rides.
type RideRepository interface {
	Find(ID string) (*models.Ride, error)
	List() ([]*models.Ride, error)

	Create(name, description string, minAge, minHeight int, longitude, latitude float32) (*models.Ride, error)
	UpdateName(ID, name string) error
	UpdateDescription(ID, description string) error
	UpdateMinAge(ID string, age int) error
	UpdateMinHeight(ID string, height int) error
	UpdateLongitude(ID string, longitude float32) error
	UpdateLatitude(ID string, latitude float32) error
	Delete(ID string) error

	AddPicture(ID string, format models.PictureFormat, data []byte) (int, error)
	UpdatePictureIndex(ID string, fromIndex, toIndex int) error
	DeletePicture(ID string, pictureIndex int) error
}
