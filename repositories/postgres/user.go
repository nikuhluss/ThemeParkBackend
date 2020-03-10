package postgres

import (
	"database/sql"
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
		Select("ID").
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

	var genderID sql.NullString
	_ = db.Get(&genderID, selectGenderID, user.Gender.String)

	// begin the transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// ANONYMOUS BLOCK FOR TRANSACTION
	{
		_, err = tx.Exec(insertUser, user.ID, user.Email, user.Email, user.PasswordSalt, user.PasswordHash, user.RegisteredOn)
		if err != nil {
			return fmt.Errorf("insertUser: %s", err)
		}

		_, err = tx.Exec(insertDetail, user.ID, genderID, user.FirstName, user.LastName, user.DateOfBirth, user.Phone, user.Address)
		if err != nil {
			return fmt.Errorf("insertDetail: %s", err)
		}

		if user.IsEmployee {
			var roleID string
			err = db.Get(&roleID, selectRoleID, user.Role.String)
			if err != nil {
				return fmt.Errorf("selectRoleID: %s", err)
			}

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
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing user in the database.
func (ur *UserRepository) Update(user *models.User) error {

	updateUser, _, _ := psql.Update("users").
		Set("email", "?").
		Where("id = ?").
		ToSql()
	
	updateDetails, _, _ := psql.Update("user_details").
		Set("gender_id", "?").
		Set("first_name", "?").
		Set("last_name", "?").
		Set("phone", "?").
		Set("address","?").
		Where("user_id = ?").
		ToSql()

	updateEmployee, _, _ := psql.Update("employees").
		Set("role_id","?").
		Where("id","?").
		ToSql()
	
	selectRoleID, _ := psql.
		Select("ID").
		From("roles").
		Where("role = $1").
		MustSql()

	selectGenderID, _ := psql.
		Select("ID").
		From("genders").
		Where("gender = $1").
		MustSql()

	var genderID sql.NullString
	_ = ur.db.Get(&genderID, selectGenderID, user.Gender.String)
	
	//Start transaction
	tx, err := ur.db.Begin()
	if err != nil {
		return err
	}

	{
		_, err = tx.Exec(updateUser, user.Email, user.ID)
		if err != nil {
			return fmt.Errorf("updateUser: %s", err)
		}

		_, err = tx.Exec(updateDetails, genderID, user.FirstName, user.LastName, user.Phone, user.Address, user.ID)
		if err != nil {
			return fmt.Errorf("updateDetails: %s", err)
		}

		if user.IsEmployee{
			var roleID string
			err = ur.db.Get(&roleID, selectRoleID, user.Role.String)
			if err != nil {
				return fmt.Errorf("selectRoleID: %s", err)
			}

			_, err = tx.Exec(updateEmployee, roleID, user.ID)
			if err != nil {
				return fmt.Errorf("updateRole: %s", err)
			}
		}

	}
	err = tx.Commit()
	if err != nil{
		return err
	}

	return nil
}

// Delete deletes an existing user from the database.
func (ur *UserRepository) Delete(ID string) error {

	// TODO account for deleting employees from maintanence
	db := ur.db

	deleteUser, _, _ := psql.
		Delete("users").
		Where("id = ?").
		ToSql()

	deleteDetails, _, _ := psql.
		Delete("user_details").
		Where("user_id = ?").
		ToSql()

	deleteEmployee, _, _ := psql.
		Delete("employees").
		Where("user_id = ?").
		ToSql()
	
	deleteCustomer, _, _ := psql.
		Delete("customers").
		Where("user_id").
		ToSql()
	
	user, err := ur.GetByID(ID)
	if err != nil {
		return fmt.Errorf("user does not exist to be deleted: %s", err)
	}
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}


	{
		

		_, err = tx.Exec(deleteDetails, ID)
		if err != nil {
			return fmt.Errorf("deleteDetails: %s", err)
		}
		//this is where removing from employee maintanence would happen
		if user.IsEmployee { 
			_, err = tx.Exec(deleteEmployee, ID)
			if err != nil {
				return fmt.Errorf("deleteEmployee: %s", err)
			}
		} else {
			_, err = tx.Exec(deleteCustomer, ID)
			if err != nil {
				return fmt.Errorf("deleteCustomer: %s", err)
			}
		}	

		_, err = tx.Exec(deleteUser, ID)
		if err != nil {
			return fmt.Errorf("deleteUser: %s", err)
		}
	}


	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// AvailableGenders returns are the valid values for gender.
func (ur *UserRepository) AvailableGenders() ([]string, error) {
	return nil, nil
}
