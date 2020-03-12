package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// PictureRepository specifies the interface for working with pictures and
// picture collections. Pictures can be individually fetched and deleted,
// however, the rest of operations are usually done on top of collections.
type PictureRepository interface {
	GetByID(ID string) (*models.Picture, error)
	Delete(ID string) error

	FetchByCollectionID(collectionID string) ([]*models.Picture, error)
	Store(collectionID string, picture *models.Picture) error
	UpdateCollectionOrdering(collectionID string, fromIndex, toIndex int) error
	DeleteCollection(collectionID string) error
}
