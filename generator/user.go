package generator

import (
	"time"

	"github.com/brianvoe/gofakeit/v4"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/crypto"
)

var (
	year1970        = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	year2000        = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	password        = "password"
	passwordHash, _ = crypto.HashFromPassword(password)
)

// InsertGender inserts the given gender and returns the generated ID.
func InsertGender(execer Execer, gender string) (string, error) {
	ID := gofakeit.UUID()

	genderInsertQuery := `
	INSERT INTO genders (ID, gender)
	VALUES ($1, $2)
	`

	_, err := execer.Exec(genderInsertQuery, ID, gender)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// InsertRole inserts the given role and returns the generated ID.
// Hourly rate is always 20.0.
func InsertRole(execer Execer, role string) (string, error) {
	ID := gofakeit.UUID()

	roleInsertQuery := `
	INSERT INTO roles (ID, role, hourly_rate)
	VALUES ($1, $2, $3)
	`

	_, err := execer.Exec(roleInsertQuery, ID, role, 20.0)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// InsertCustomer inserts a new customer with the given email and username.
// Password is always "password".
func InsertCustomer(execer Execer, username, email string) (string, error) {

	ID := gofakeit.UUID()
	registeredOn := gofakeit.DateRange(year2000, year2000.Add(time.Hour*24*365*10))

	userInsertQuery := `
	INSERT INTO users (ID, username, email, password_salt, password_hash, registered_on)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	customerInsertQuery := `
	INSERT INTO customers (user_id)
	VALUES ($1)
	`

	_, err := execer.Exec(userInsertQuery, ID, username, email, "", passwordHash, registeredOn)
	if err != nil {
		return "", err
	}

	_, err = execer.Exec(customerInsertQuery, ID)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// InsertEmployee inserts a new employee with the given email, username, and role.
// Password is always "password".
func InsertEmployee(execer Execer, username, email, roleID string) (string, error) {

	ID := gofakeit.UUID()
	registeredOn := gofakeit.DateRange(year2000, year2000.Add(time.Hour*24*365*10))

	userInsertQuery := `
	INSERT INTO users (ID, username, email, password_salt, password_hash, registered_on)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	customerInsertQuery := `
	INSERT INTO customers (user_id)
	VALUES ($1)
	`

	employeeInsertQuery := `
	INSERT INTO employees (ID, user_ID, role_ID)
	VALUES ($1, $2, $3)
	`

	_, err := execer.Exec(userInsertQuery, ID, username, email, "", passwordHash, registeredOn)
	if err != nil {
		return "", err
	}

	_, err = execer.Exec(customerInsertQuery, ID)
	if err != nil {
		return "", err
	}

	_, err = execer.Exec(employeeInsertQuery, ID, ID, roleID)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// InsertUserDetailsWithName inserts user details for the given user using the
// given first name and last name. Other details are random.
func InsertUserDetailsWithName(execer Execer, userID, genderID, firstName, lastName string) error {
	dateOfBirth := gofakeit.DateRange(year1970, year1970.Add(time.Hour*24*365*10))
	phone := gofakeit.Phone()
	address := gofakeit.Address().Street

	detailsInsertQuery := `
	INSERT INTO user_details (user_ID, gender_ID, first_name, last_name, date_of_birth, phone, address)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := execer.Exec(detailsInsertQuery, userID, genderID, firstName, lastName, dateOfBirth, phone, address)
	if err != nil {
		return err
	}

	return nil
}

// InsertUserDetails inserts random user details for the given user.
func InsertUserDetails(execer Execer, userID string, genderID string) error {
	firstName := gofakeit.FirstName()
	lastName := gofakeit.LastName()
	return InsertUserDetailsWithName(execer, userID, genderID, firstName, lastName)
}

// MustInsertGender is like InsertGender but panics on error.
func MustInsertGender(mustExecer MustExecer, gender string) string {
	return MustInsert(InsertGender(&AsExecer{mustExecer}, gender))
}

// MustInsertRole is like InsertRole but panics on error.
func MustInsertRole(mustExecer MustExecer, role string) string {
	return MustInsert(InsertRole(&AsExecer{mustExecer}, role))
}

// MustInsertCustomer is like InsertCustomer but panics on error.
func MustInsertCustomer(mustExecer MustExecer, email, username string) string {
	return MustInsert(InsertCustomer(&AsExecer{mustExecer}, email, username))
}

// MustInsertEmployee is like InsertEmployee but panics on error.
func MustInsertEmployee(mustExecer MustExecer, email, username, roleID string) string {
	return MustInsert(InsertEmployee(&AsExecer{mustExecer}, email, username, roleID))
}

// MustInsertUserDetails is like InsertUserDetails but panics on error.
func MustInsertUserDetails(mustExecer MustExecer, userID, genderID string) {
	err := InsertUserDetails(&AsExecer{mustExecer}, userID, genderID)
	if err != nil {
		panic(err)
	}
}
