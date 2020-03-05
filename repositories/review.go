package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// ReviewRepository defines the interface for working with reviews.
type ReviewRepository interface {
	GetByID(ID string) (*models.Review, error)

	Fetch() ([]*models.Review, error)
	FetchForRideSortedByRating(rideID string) ([]*models.Review, error)
	FetchForRideSortedByDate(rideID string) ([]*models.Review, error)

	Store(*models.Review) error
	Update(*models.Review) error
	Delete(ID string) error
}
