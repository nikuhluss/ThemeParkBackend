package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// psql is a statement builder that uses Dollar format (Postgres).
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// selectUsers is a query template that we can re-use for some queries below (GetBy*).
var selectUsers = psql.
	Select(
		"users.*",
		"user_details.*",
		"genders.gender",
		"(employees.ID IS NOT NULL) AS is_employee",
		"roles.role",
		"COALESCE(roles.hourly_rate, 0.0) as hourly_rate",
	).
	From("users").
	LeftJoin("user_details ON user_details.user_ID = users.ID").
	LeftJoin("genders ON genders.ID = user_details.gender_ID").
	LeftJoin("employees ON employees.user_ID = users.ID").
	LeftJoin("roles ON roles.ID = employees.role_ID")

// UserRepository implements the UserRepository interface for postgres.
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance using the given
// database instance.
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

// GetByID fetches an user from the postgres `users` and `user_details` tables.
func (ur *UserRepository) GetByID(ID string) (*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query, _ := selectUsers.Where("users.ID = $1").Limit(1).MustSql()

	user := models.User{}
	err := udb.Get(&user, query, ID)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// GetByEmail is similar to GetByID but uses the email for finding the user.
func (ur *UserRepository) GetByEmail(email string) (*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query, _ := selectUsers.Where("users.email = $1").Limit(1).MustSql()

	user := models.User{}
	err := udb.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// GetByUsername is similar to GetByID but uses the username for finding the user.
func (ur *UserRepository) GetByUsername(username string) (*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query, _ := selectUsers.Where("users.username = $1").Limit(1).MustSql()

	user := models.User{}
	err := udb.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, err
}

// Fetch fetches all users from the postgres `users` and `user_details` tables.
func (ur *UserRepository) Fetch() ([]*models.User, error) {
	db := ur.db
	udb := db.Unsafe()

	query, _ := selectUsers.MustSql()

	users := []*models.User{}
	err := udb.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, err
}

// FetchCustomers is like Fetch, but fetches only the customers.
func (ur *UserRepository) FetchCustomers() ([]*models.User, error) {
	users, err := ur.Fetch()
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

// FetchEmployees is like Fetch, but fetches only the employees.
func (ur *UserRepository) FetchEmployees() ([]*models.User, error) {
	users, err := ur.Fetch()
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

// Store creates a new user in the database.
func (ur *UserRepository) Store(user *models.User) error {

	db := ur.db

	selectGenderID, _ := psql.
		Select("ID").
		From("genders").
		Where("gender = $1").
		MustSql()

	selectRoleID, _ := psql.
		Select("role_ID").
		From("roles").
		Where("role = $1").
		MustSql()

	insertUser, _, _ := psql.
		Insert("users").
		Columns("ID", "email", "username", "password_salt", "password_hash", "registered_on").
		Values("?", "?", "?", "?", "?", "?").
		ToSql()

	insertDetail, _, _ := psql.
		Insert("user_details").
		Columns("user_ID", "gender_ID", "first_name", "last_name", "date_of_birth", "phone", "address").
		Values("?", "?", "?", "?", "?", "?", "?").
		ToSql()

	insertEmployee, _, _ := psql.
		Insert("employees").
		Columns("ID", "user_ID", "role_ID").
		Values("?", "?", "?").
		ToSql()

	insertCustomer, _, _ := psql.
		Insert("customers").
		Columns("user_ID").
		Values("?").
		ToSql()

	var genderID *string
	_ = db.Get(genderID, selectGenderID, user.Gender.String)

	var roleID *string
	_ = db.Get(roleID, selectRoleID, user.Role.String)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(insertUser, user.ID, user.Email, user.Email, user.PasswordSalt, user.PasswordHash, user.RegisteredOn)
	if err != nil {
		return fmt.Errorf("insertUser: %s", err)
	}

	_, err = tx.Exec(insertDetail, user.ID, genderID, user.FirstName, user.LastName, user.DateOfBirth, user.Phone, user.Address)
	if err != nil {
		return fmt.Errorf("insertDetail: %s", err)
	}

	if user.IsEmployee {
		_, err = tx.Exec(insertEmployee, user.ID, user.ID, roleID)
		if err != nil {
			return fmt.Errorf("insertEmployee: %s", err)
		}
	} else {
		_, err := tx.Exec(insertCustomer, user.ID)
		if err != nil {
			return fmt.Errorf("insertCustomer: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing user in the database.
func (ur *UserRepository) Update(*models.User) error {
	return nil
}

// Delete deletes an existing user from the database.
func (ur *UserRepository) Delete(ID string) error {
	return nil
}

// AvailableGenders returns are the valid values for gender.
func (ur *UserRepository) AvailableGenders() ([]string, error) {
	return nil, nil
}
