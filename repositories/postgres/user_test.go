package postgres_test

import (
	"database/sql"
	"testing"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func setupTestUsers(db *sqlx.DB) ([]string, []string) {

	tx := db.MustBegin()

	tx.MustExec("TRUNCATE TABLE users CASCADE")
	tx.MustExec("TRUNCATE TABLE genders CASCADE")
	tx.MustExec("TRUNCATE TABLE user_details CASCADE")
	tx.MustExec("TRUNCATE TABLE roles CASCADE")
	tx.MustExec("TRUNCATE TABLE employees CASCADE")

	maleGender := generator.MustInsertGender(tx, "Male")
	femaleGender := generator.MustInsertGender(tx, "Female")
	otherGender := generator.MustInsertGender(tx, "Other")

	workerRole := generator.MustInsertRole(tx, "Worker")

	customer0 := generator.MustInsertCustomer(tx, "customer0", "customer0@email.com")
	customer1 := generator.MustInsertCustomer(tx, "customer1", "customer1@email.com")
	employee0 := generator.MustInsertEmployee(tx, "employee0", "employee0@email.com", workerRole)

	generator.MustInsertUserDetails(tx, customer0, maleGender)
	generator.MustInsertUserDetails(tx, customer1, femaleGender)
	generator.MustInsertUserDetails(tx, employee0, otherGender)

	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return []string{customer0, customer1}, []string{employee0}
}

// Tests
// --------------------------------

func TestUserGetByIDSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	customerIDs, employeeIDs := setupTestUsers(db)

	for _, customerID := range customerIDs {
		customer, err := userRepository.GetByID(customerID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, customerID, customer.ID)
		assert.False(t, customer.IsEmployee)
	}

	for _, employeeID := range employeeIDs {
		employee, err := userRepository.GetByID(employeeID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}
		assert.Equal(t, employeeID, employee.ID)
		assert.True(t, employee.IsEmployee)
	}
}

func TestUserGetByIDNoMatchFails(t *testing.T) {
	userRepository, _, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	user, err := userRepository.GetByID("some-unknown-ID")
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestUserFetchSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	users, err := userRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, users, 3)
}

func TestUserFetchCustomersSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	customers, err := userRepository.FetchCustomers()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, customers, 2)
}

func TestUserFetchEmployeesSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	employees, err := userRepository.FetchEmployees()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, employees, 1)
}

func TestUserStoreCustomerSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
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

func TestUserStoreEmployeeSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	role := "Worker"
	employee := models.NewEmployee("customer--A", "customer--A--email", "customer--A--password_salt", "customer--A--password_hash", role, 0)
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
	assert.Equal(t, sql.NullString{String: role, Valid: true}, user.Role)
	assert.Greater(t, user.HourlyRate, float32(0))
}

func TestUserUpdateCustomerSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	customerIDs, _ := setupTestUsers(db)
	userID := customerIDs[0]

	user, err := userRepository.GetByID(userID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	// create expected user. Note that not all values can be updated (password, registered_on, etc)
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

	updatedUser, err := userRepository.GetByID(user.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedUser, updatedUser)
}

func TestDeleteUserSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	customerIDs, _ := setupTestUsers(db)
	userID := customerIDs[0]

	// delete a user, check if he isn't there
	err := userRepository.Delete(userID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	user, err := userRepository.GetByID(userID)
	assert.Nil(t, user)
	assert.NotNil(t, err)
}

func TestUserUpdatePasswordSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	customerIDs, _ := setupTestUsers(db)
	userID := customerIDs[0]

	err := userRepository.UpdatePassword(userID, "abc", "new password")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedUser, err := userRepository.GetByID(userID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "abc", updatedUser.PasswordSalt)
	assert.Equal(t, "new password", updatedUser.PasswordHash)

}

func TestUserGetAllGendersSucceeds(t *testing.T) {
	userRepository, db, teardown := testutil.MakeUserRepositoryFixture()
	defer teardown()

	setupTestUsers(db)

	genders, err := userRepository.AvailableGenders()

	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "Female", genders[0])
	assert.Equal(t, "Male", genders[1])
	assert.Equal(t, "Other", genders[2])
}
