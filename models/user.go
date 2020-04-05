package models

import (
	"database/sql"
	"time"
)

// User struct represents an user in the system. Note that this struct might
// be the composition of one or more tables.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordSalt string    `db:"password_salt" json:"passwordSalt"`
	PasswordHash string    `db:"password_hash" json:"passwordHash"`
	RegisteredOn time.Time `db:"registered_on" json:"registeredOn"`

	Gender      NullString `json:"gender"`
	FirstName   NullString `db:"first_name" json:"firstName"`
	LastName    NullString `db:"last_name" json:"lastName"`
	DateOfBirth NullTime   `db:"date_of_birth" json:"dateOfBirth"`
	Phone       NullString `json:"phone"`
	Address     NullString `json:"address"`

	IsEmployee bool           `db:"is_employee" json:"isEmployee"`
	Role       sql.NullString `json:"role"`
	HourlyRate float32        `db:"hourly_rate" json:"hourlyRate"`
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
