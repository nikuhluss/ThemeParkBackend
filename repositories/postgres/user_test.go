package postgres

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
)

// Fixtures
// --------------------------------

func UserRepositoryFixture() (*UserRepository, *sqlx.DB, func()) {
	dbconfig := testutil.NewDatabaseConnectionConfig()
	db, err := testutil.NewDatabaseConnection(dbconfig)
	if err != nil {
		panic(fmt.Sprintf("user_test.go: UserRepositoryFixture: %s", err))
	}

	userRepository := NewUserRepository(db)
	return userRepository, db, func() {
		db.Close()
	}
}

func setupTestUsers(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE users CASCADE")
	db.MustExec("TRUNCATE TABLE genders CASCADE")
	db.MustExec("TRUNCATE TABLE user_details CASCADE")
	db.MustExec("TRUNCATE TABLE roles CASCADE")
	db.MustExec("TRUNCATE TABLE employees CASCADE")

	userInsertQuery := `
	INSERT INTO users (ID, username, email, password_salt, password_hash, registered_on)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	genderInsertQuery := `
	INSERT INTO genders (ID, gender)
	VALUES ($1, $2)
	`

	detailsInsertQuery := `
	INSERT INTO user_details (user_ID, gender_ID, first_name, last_name, date_of_birth, phone, address)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	roleInsertQuery := `
	INSERT INTO roles (ID, role, hourly_rate)
	VALUES ($1, $2, $3)
	`

	employeeInsertQuery := `
	INSERT INTO employees (ID, user_ID, role_ID)
	VALUES ($1, $2, $3)
	`

	db.MustExec(userInsertQuery, "user--A", "user--A--username", "user--A--email", "user--A--password_salt", "user--A--password_hash", time.Now())
	db.MustExec(userInsertQuery, "user--B", "user--B--username", "user--B--email", "user--B--password_salt", "user--B--password_hash", time.Now())
	db.MustExec(userInsertQuery, "user--C", "user--C--username", "user--C--email", "user--C--password_salt", "user--C--password_hash", time.Now())

	db.MustExec(genderInsertQuery, "gender--A", "Male")
	db.MustExec(genderInsertQuery, "gender--B", "Female")
	db.MustExec(genderInsertQuery, "gender--C", "Other")

	db.MustExec(detailsInsertQuery, "user--A", "gender--A", "user--A--first_name", "user--A--last_name", time.Now(), "user--A--phone", "user--A--address")
	db.MustExec(detailsInsertQuery, "user--B", "gender--B", "user--B--first_name", "user--B--last_name", time.Now(), "user--B--phone", "user--B--address")
	db.MustExec(detailsInsertQuery, "user--C", "gender--C", "user--C--first_name", "user--C--last_name", time.Now(), "user--C--phone", "user--C--address")

	db.MustExec(roleInsertQuery, "role--C", "role--C", 15.50)

	db.MustExec(employeeInsertQuery, "employee--C", "user--C", "role--C")
}

// Tests
// --------------------------------

func TestFindSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	tests := []struct {
		userID     string
		gender     string
		isEmployee bool
	}{
		{"user--A", "Male", false},
		{"user--B", "Female", false},
		{"user--C", "Other", true},
	}

	for _, tt := range tests {
		user, err := userRepository.Find(tt.userID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.userID, user.ID)
		assert.Equal(t, tt.userID+"--email", user.Email)
		assert.Equal(t, tt.gender, user.Gender.String)
		assert.Equal(t, tt.userID+"--first_name", user.FirstName.String)
		assert.Equal(t, tt.userID+"--last_name", user.LastName.String)
		assert.Equal(t, tt.userID+"--phone", user.Phone.String)
		assert.Equal(t, tt.userID+"--address", user.Address.String)
		assert.Equal(t, tt.isEmployee, user.IsEmployee)
	}
}

func TestFindNoMatchFails(t *testing.T) {
	UserRepository, _, teardown := UserRepositoryFixture()
	defer teardown()

	user, err := UserRepository.Find("some-unknown-ID")
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestListSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	users, err := userRepository.List()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, users, 3)
}

func TestListCustomersSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	customers, err := userRepository.ListCustomers()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, customers, 2)
}

func TestListEmployeesSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	employees, err := userRepository.ListEmployees()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, employees, 1)
}
