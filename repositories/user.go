package repositories

import (
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// UserRepository defines the interface for working with users.
type UserRepository interface {
	Find(ID string) (*models.User, error)
	List() ([]*models.User, error)
	ListCustomers() ([]*models.User, error)
	ListEmployees() ([]*models.User, error)

	CreateCustomer(email, passwordSalt, passwordHash string) (*models.User, error)
	CreateEmployee(email, passwordSalt, passwordHash string) (*models.User, error)
	UpdateGender(ID, gender string) error
	UpdateFirstName(ID, firstName string) error
	UpdateLastName(ID, lastName string) error
	UpdateDateOfBirth(ID string, dateOfBirth time.Time) error
	UpdatePhone(ID, phone string) error
	UpdateAddress(ID, address string) error
	Delete(ID string) error

	AvailableGenders() ([]string, error)
}
