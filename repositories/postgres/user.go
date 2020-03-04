package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// UserRepository implements the UserRepository interface for postgres.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance using the given
// database instance.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

// Find fetches an user from the postgres `users` and `user_details` tables.
func (ur *UserRepository) Find(ID string) (*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query := `
	SELECT
		users.*,
		user_details.*,
		genders.gender,
		(employees.ID IS NOT NULL) as is_employee
	FROM users
	LEFT JOIN user_details ON user_details.user_ID = users.ID
	INNER JOIN genders ON genders.ID = user_details.gender_ID
	LEFT JOIN employees ON employees.user_ID = users.ID
	WHERE users.ID = $1
	LIMIT 1
	`

	user := models.User{}
	err := udb.Get(&user, query, ID)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// List fetches all users from the postgres `users` and `user_details` tables.
func (ur *UserRepository) List() ([]*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query := `
	SELECT
		users.*,
		user_details.*,
		genders.gender,
		(employees.ID IS NOT NULL) as is_employee
	FROM users
	LEFT JOIN user_details ON user_details.user_ID = users.ID
	INNER JOIN genders ON genders.ID = user_details.gender_ID
	LEFT JOIN employees ON employees.user_ID = users.ID
	`

	users := []*models.User{}
	err := udb.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, err
}

// ListCustomers is like List, but fetches only the customers.
func (ur *UserRepository) ListCustomers() ([]*models.User, error) {
	users, err := ur.List()
	if err != nil {
		return nil, err
	}

	customers := make([]*models.User, 0, len(users))
	for _, user := range users {
		if user.IsEmployee {
			continue
		}
		customers = append(customers, user)
	}

	return customers, nil
}

// ListEmployees is like List, but fetches only the employees.
func (ur *UserRepository) ListEmployees() ([]*models.User, error) {
	users, err := ur.List()
	if err != nil {
		return nil, err
	}

	employees := make([]*models.User, 0, len(users))
	for _, user := range users {
		if !user.IsEmployee {
			continue
		}
		employees = append(employees, user)
	}

	return employees, nil
}

// CreateCustomer creates a new customer using the given information. To add personal
// details, use the update functions after creation.
func CreateCustomer(email, passwordSalt, passwordHash string) (*models.User, error) {
	return nil, nil
}

// CreateEmployee creates a new employee using the given information. To add personal
// details, use the update functions after creation.
func CreateEmployee(email, passwordSalt, passwordHash string) (*models.User, error) {
	return nil, nil
}

// UpdateGender changes the gender for the given user by querying available
// genders and comparing against requested gender, if there's a match, grabs
// and uses the ID of the matched gender.
func (ur *UserRepository) UpdateGender(ID, gender string) error {
	return nil
}

// UpdateFirstName updates the first name of the given user (if any).
func (ur *UserRepository) UpdateFirstName(ID, firstName string) error {
	return nil
}

// UpdateLastName updates the last name of the given user (if any).
func (ur *UserRepository) UpdateLastName(ID, lastName string) error {
	return nil
}

// UpdateDateOfBirth date of birth updates the date of birth of the given user (if any).
func (ur *UserRepository) UpdateDateOfBirth(ID string, dateOfBirth time.Time) error {
	return nil
}

// UpdatePhone phone updates the phone number of the given user (if any).
func (ur *UserRepository) UpdatePhone(ID, phone string) error {
	return nil
}

// UpdateAddress updates the address of the given user (if any).
func (ur *UserRepository) UpdateAddress(ID, address string) error {
	return nil
}

func Delete(ID string) error {
	return nil
}

// AvailableGenders returns are the valid values for gender.
func (ur *UserRepository) AvailableGenders() ([]string, error) {
	return nil, nil
}
