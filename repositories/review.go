package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// ReviewRepository defines the interface for working with reviews.
type ReviewRepository interface {
	Find(ID string) (*models.Review, error)
	List() ([]*models.Review, error)
	ListForRideSortedByRating(rideID string) ([]*models.Review, error)
	ListForRideSortedByDate(rideID string) ([]*models.Review, error)

	Create(ID, userID, title, content string) (*models.Review, error)
	UpdateRating(ID string, rating int) error
	UpdateTitle(ID, title string) error
	UpdateContent(ID, content string) error
	Delete(ID string) error
}
