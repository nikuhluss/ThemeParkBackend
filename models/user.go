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
	PasswordSalt string    `db:"password_salt"`
	PasswordHash string    `db:"password_hash"`
	RegisteredOn time.Time `db:"registered_on"`

	Gender      sql.NullString
	FirstName   sql.NullString `db:"first_name"`
	LastName    sql.NullString `db:"last_name"`
	DateOfBirth sql.NullTime   `db:"date_of_birth"`
	Phone       sql.NullString
	Address     sql.NullString

	IsEmployee bool `db:"is_employee"`
	Role       sql.NullString
	HourlyRate float32 `db:"hourly_rate"`
}

// NewCustomer returns a new User instance that is a customer.
func NewCustomer(ID, email, passwordSalt, passwordHash string) *User {
	return &User{
		ID:           ID,
		Email:        email,
		PasswordSalt: passwordSalt,
		PasswordHash: passwordHash,
		RegisteredOn: time.Now(),
		IsEmployee:   false,
	}
}

// NewEmployee returns a new User instance that is an employee.
func NewEmployee(ID, email, passwordSalt, passwordHash, role string, hourlyRate float32) *User {
	return &User{
		ID:           ID,
		Email:        email,
		PasswordSalt: passwordSalt,
		PasswordHash: passwordHash,
		RegisteredOn: time.Now(),
		IsEmployee:   true,
		Role:         sql.NullString{String: role, Valid: true},
		HourlyRate:   hourlyRate,
	}
}
