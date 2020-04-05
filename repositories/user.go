package repositories

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// UserRepository defines the interface for working with users.
type UserRepository interface {
	GetByID(ID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)

	Fetch() ([]*models.User, error)
	FetchCustomers() ([]*models.User, error)
	FetchEmployees() ([]*models.User, error)

	Store(*models.User) error
	Update(*models.User) error
	Delete(ID string) error

	UpdatePassword(ID, passwordSalt, passwordHash string) error
	AvailableGenders() ([]string, error)
}
