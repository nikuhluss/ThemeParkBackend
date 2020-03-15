package postgres

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
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

	db.MustExec(genderInsertQuery, "gender--male", "Male")
	db.MustExec(genderInsertQuery, "gender--female", "Female")
	db.MustExec(genderInsertQuery, "gender--other", "Other")

	db.MustExec(detailsInsertQuery, "user--A", "gender--male", "user--A--first_name", "user--A--last_name", time.Now(), "user--A--phone", "user--A--address")
	db.MustExec(detailsInsertQuery, "user--B", "gender--female", "user--B--first_name", "user--B--last_name", time.Now(), "user--B--phone", "user--B--address")
	db.MustExec(detailsInsertQuery, "user--C", "gender--other", "user--C--first_name", "user--C--last_name", time.Now(), "user--C--phone", "user--C--address")

	db.MustExec(roleInsertQuery, "role--C", "role--C", 15.50)

	db.MustExec(employeeInsertQuery, "employee--C", "user--C", "role--C")
}

// Utility
// --------------------------------

func assertUserLoginEqual(t *testing.T, expected *models.User, actual *models.User) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.PasswordSalt, actual.PasswordSalt)
	assert.Equal(t, expected.PasswordHash, actual.PasswordHash)
}

func assertUserDetailsEqual(t *testing.T, expected *models.User, actual *models.User) {

}

// Tests
// --------------------------------

func TestGetByIDSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	tests := []struct {
		userID     string
		gender     string
		isEmployee bool
		hourlyRate float32
	}{
		{"user--A", "Male", false, 0.0},
		{"user--B", "Female", false, 0.0},
		{"user--C", "Other", true, 15.50},
	}

	for _, tt := range tests {
		user, err := userRepository.GetByID(tt.userID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.userID, user.ID)
		assert.Equal(t, tt.userID+"--email", user.Email)
		assert.Equal(t, tt.userID+"--password_salt", user.PasswordSalt)
		assert.Equal(t, tt.userID+"--password_hash", user.PasswordHash)
		assert.Equal(t, tt.gender, user.Gender.String)
		assert.Equal(t, tt.userID+"--first_name", user.FirstName.String)
		assert.Equal(t, tt.userID+"--last_name", user.LastName.String)
		assert.Equal(t, tt.userID+"--phone", user.Phone.String)
		assert.Equal(t, tt.userID+"--address", user.Address.String)
		assert.Equal(t, tt.isEmployee, user.IsEmployee)
		assert.Equal(t, tt.hourlyRate, user.HourlyRate)
	}
}

func TestGetByIDNoMatchFails(t *testing.T) {
	UserRepository, _, teardown := UserRepositoryFixture()
	defer teardown()

	user, err := UserRepository.GetByID("some-unknown-ID")
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestFetchSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	users, err := userRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, users, 3)
}

func TestFetchCustomersSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	customers, err := userRepository.FetchCustomers()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, customers, 2)
}

func TestFetchEmployeesSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	employees, err := userRepository.FetchEmployees()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, employees, 1)
}

func TestStoreCustomerSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	customer := models.NewCustomer("customer--A", "customer--A--email", "customer--A--password_salt", "customer--A--password_hash")
	err := userRepository.Store(customer)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	user, err := userRepository.GetByID(customer.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, user)
	assert.False(t, user.IsEmployee)
}

func TestStoreEmployeeSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	employee := models.NewEmployee("customer--A", "customer--A--email", "customer--A--password_salt", "customer--A--password_hash", "role--C", 15.50)
	err := userRepository.Store(employee)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	user, err := userRepository.GetByID(employee.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, user)
	assert.True(t, user.IsEmployee)
	assert.Equal(t, user.Role.String, "role--C")
	assert.Equal(t, user.HourlyRate, float32(15.50))
}

func TestUpdateCustomerSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	user, err := userRepository.GetByID("user--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	// create expected user. Note that not all values ca be updated (password, registered_on, etc)
	expectedUser := models.NewCustomer(user.ID, "expected--Email", user.PasswordSalt, user.PasswordHash)
	expectedUser.RegisteredOn = user.RegisteredOn
	expectedUser.Gender = sql.NullString{String: "Other", Valid: true}
	expectedUser.FirstName = sql.NullString{String: "expected--first_name", Valid: true}
	expectedUser.LastName = sql.NullString{String: "expected--last_name", Valid: true}
	expectedUser.Phone = sql.NullString{String: "expected--phone", Valid: true}
	expectedUser.Address = sql.NullString{String: "expected--address", Valid: true}

	err = userRepository.Update(expectedUser)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedUser, err := userRepository.GetByID("user--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedUser, updatedUser)
}

func TestDeleteUserSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	//delete a user, check if he isn't there
	err := userRepository.Delete("user--C")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	user, err := userRepository.GetByID("user--C")
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestGetAllGendersSucceeds(t *testing.T) {
	userRepository, db, teardown := UserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	user, err := userRepository.AvailableGenders()

	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "Other", *user[0])
	assert.Equal(t, "Female", *user[1])
	assert.Equal(t, "Male", *user[2])
}