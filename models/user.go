package models

import (
	"database/sql"
	"time"
)

// User struct represents an user in the system. Note that this struct might
// be the composition of one or more tables.
type User struct {
	ID           string
	Email        string
	RegisteredOn time.Time

	Gender      sql.NullString
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	DateOfBirth sql.NullTime   `db:"date_of_birth"`
	Phone       sql.NullString
	Address     sql.NullString

	IsEmployee bool `db:"is_employee"`
}

// NewCustomer returns a new User instance that is a customer.
func NewCustomer(ID, email string) *User {
	return &User{
		ID:           ID,
		Email:        email,
		RegisteredOn: time.Now(),
		IsEmployee:   false,
	}
}

// NewEmployee returns a new User instance that is an employee.
func NewEmployee(ID, email string) *User {
	return &User{
		ID:           ID,
		Email:        email,
		RegisteredOn: time.Now(),
		IsEmployee:   true,
	}
}
