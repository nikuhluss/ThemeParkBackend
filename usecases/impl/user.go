package impl

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	//"golang.org/x/sync/errgroup"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/mathutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errUserExists        = fmt.Errorf("user with the given ID already exists")
	errUserDoesNotExists = fmt.Errorf("user with he given ID does not exists")
)

// UserUsecaseImpl implements the UserUsecase interface.
type UserUsecaseImpl struct {
	userRepo repos.UserRepository
	timeout  time.Duration
}

// NewUserUsecaseImpl returns a new UserUsecaseImpl instance. The timeout
// parameter specifies a duration for each request before throwing and error.
func NewUserUsecaseImpl(
	userRepo repos.UserRepository,
	timeout time.Duration) *UserUsecaseImpl {

	return &UserUsecaseImpl{
		userRepo,
		timeout,
	}
}

// GetByID fetches user from the repositories using the given ID.
func (uu *UserUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.User, error) {
	user, err := uu.userRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %s", err)
	}

	return user, nil
}

// GetByEmail fetches user from the repositories using the given email.
func (uu *UserUsecaseImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := uu.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %s", err)
	}

	return user, nil
}

// Fetch fetches all Users from the repositories.
func (uu *UserUsecaseImpl) Fetch(ctx context.Context) ([]*models.User, error) {
	allUser, err := uu.userRepo.Fetch()
	if err != nil {
		return nil, err
	}

	return allUser, nil
}

// FetchCustomers fetches all Customers from the repositories.
func (uu *UserUsecaseImpl) FetchCustomers(ctx context.Context) ([]*models.User, error) {
	allCustomer, err := uu.userRepo.FetchCustomers()
	if err != nil {
		return nil, err
	}

	return allCustomer, nil
}

// FetchEmployees fetches all Employees from the repositories.
func (uu *UserUsecaseImpl) FetchEmployees(ctx context.Context) ([]*models.User, error) {
	allEmployee, err := uu.userRepo.FetchEmployees()
	if err != nil {
		return nil, err
	}

	return allEmployee, nil
}

// Store creates a new user in the repository if a user with the same ID
// doesn't exists already.
func (uu *UserUsecaseImpl) Store(ctx context.Context, user *models.User) error {
	_, err := uu.userRepo.GetByID(user.ID)
	if err == nil {
		return errUserExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	user.ID = uuid
	cleanUser(user)
	if user.IsEmployee {
		cleanEmployee(user)
	}

	err = validateUser(user)
	if err != nil {
		return err
	}

	err = uu.userRepo.Store(user)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing user in the repository.
func (uu *UserUsecaseImpl) Update(ctx context.Context, user *models.User) error {
	_, err := uu.userRepo.GetByID(user.ID)
	if err != nil {
		return errUserDoesNotExists
	}

	cleanUser(user)
	if user.IsEmployee {
		cleanEmployee(user)
	}

	err = validateUser(user)
	if err != nil {
		return err
	}

	err = uu.userRepo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func cleanUser(user *models.User) {
	user.ID = strings.TrimSpace(user.ID)
	user.Email = strings.TrimSpace(user.Email)
	user.PasswordSalt = strings.TrimSpace(user.PasswordSalt)
	user.PasswordHash = strings.TrimSpace(user.PasswordHash)
	user.Gender.String = strings.TrimSpace(user.Gender.String)
	user.FirstName.String = strings.TrimSpace(user.FirstName.String)
	user.LastName.String = strings.TrimSpace(user.LastName.String)
	user.Phone.String = strings.TrimSpace(user.Phone.String)
	user.Address.String = strings.TrimSpace(user.Address.String)
}

func cleanEmployee(employee *models.User) {
	employee.Role.String = strings.TrimSpace(employee.Role.String)
	employee.HourlyRate = mathutil.ClampFloat32(employee.HourlyRate, 0, 500)
}

func validateUser(user *models.User) error {
	var validEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var validName = regexp.MustCompile(`^[A-z]`)
	var validPhone = regexp.MustCompile(`^[0-9]+$`)

	if len(user.ID) <= 0 {
		return fmt.Errorf("validateUser: ID must be non-empty")
	}

	if len(user.Email) <= 0 {
		return fmt.Errorf("validateUser: email must be non-empty")
	}

	if user.FirstName.Valid && len(user.FirstName.String) <= 0 {
		return fmt.Errorf("validateUser: first name must be non-empty")
	}

	if user.LastName.Valid && len(user.LastName.String) <= 0 {
		return fmt.Errorf("validateUser: last name must be non-empty")
	}

	if user.Phone.Valid && len(user.Phone.String) <= 0 {
		return fmt.Errorf("validateUser: phone number must be non-empty")
	}

	if user.Address.Valid && len(user.Address.String) <= 0 {
		return fmt.Errorf("validateUser: address must be non-empty")
	}

	if !validEmail.MatchString(strings.ToLower(user.Email)) {
		return fmt.Errorf("validateUser: invalid email address format")
	}

	if user.FirstName.Valid && !validName.MatchString(strings.ToLower(user.FirstName.String)) {
		return fmt.Errorf("validateUser: invalid first name format")
	}

	if user.LastName.Valid && !validName.MatchString(strings.ToLower(user.LastName.String)) {
		return fmt.Errorf("validateUser: invalid last name format")
	}

	if user.Phone.Valid && !validPhone.MatchString(strings.ToLower(user.Phone.String)) {
		return fmt.Errorf("validateUser: invalid phone number format")
	}

	return nil
}
