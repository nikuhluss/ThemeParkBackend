package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/mathutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errRideExists        = fmt.Errorf("ride with the given ID already exists")
	errRideDoesNotExists = fmt.Errorf("ride with he given ID does not exists")
)

// RideUsecaseImpl implements the RideUsecase interface.
type RideUsecaseImpl struct {
	rideRepo    repos.RideRepository
	pictureRepo repos.PictureRepository
	reviewRepo  repos.ReviewRepository
	timeout     time.Duration
}

// NewRideUsecaseImpl returns a new RideUsecaseImpl instance. The timeout
// parameter specifies a duration for each request before throwing and error.
func NewRideUsecaseImpl(
	rideRepo repos.RideRepository,
	pictureRepo repos.PictureRepository,
	reviewRepo repos.ReviewRepository,
	timeout time.Duration) *RideUsecaseImpl {

	return &RideUsecaseImpl{
		rideRepo,
		pictureRepo,
		reviewRepo,
		timeout,
	}
}

// GetByID fetches ride from the repositories using the given ID.
func (ru *RideUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.Ride, error) {

	ride, err := ru.rideRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride: %s", err)
	}

	pictures, err := ru.pictureRepo.FetchByCollectionID(ride.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride pictures: %s", err)
	}

	reviews, err := ru.reviewRepo.FetchForRideSortedByDate(ride.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride reviews: %s", err)
	}

	reviewsTotal := 0
	for _, review := range reviews {
		reviewsTotal += review.Rating
	}

	ride.Pictures = pictures
	ride.Reviews = reviews
	ride.ReviewsAverage = 0
	if len(reviews) > 0 {
		ride.ReviewsAverage = reviewsTotal / len(reviews)
	}

	return ride, nil
}

// Fetch fetches all rides from the repositories.
func (ru *RideUsecaseImpl) Fetch(ctx context.Context) ([]*models.Ride, error) {

	rides, err := ru.rideRepo.Fetch()
	if err != nil {
		return nil, fmt.Errorf("error fetching rides: %s", err)
	}

	timeoutContext, cancel := context.WithTimeout(ctx, ru.timeout)
	defer cancel()
	eg, _ := errgroup.WithContext(timeoutContext)

	// the following section fetches reviews to calculate review averages

	type RideReviewsAverage struct {
		rideID         string
		reviewsAverage int
	}

	chanReviewAverages := make(chan RideReviewsAverage)

	for _, ride := range rides {
		rideID := ride.ID

		eg.Go(func() error {
			reviews, err := ru.reviewRepo.FetchForRideSortedByDate(rideID)
			if err != nil {
				return err
			}

			reviewsTotal := 0
			for _, review := range reviews {
				reviewsTotal += review.Rating
			}

			reviewsAverage := 0
			if len(reviews) > 0 {
				reviewsAverage = reviewsTotal / len(reviews)
			}

			chanReviewAverages <- RideReviewsAverage{rideID, reviewsAverage}
			return nil
		})

	}

	// close channels if error group returns from Wait

	go func() {
		err := eg.Wait()
		if err != nil {
			// TODO: log error
			return
		}
		close(chanReviewAverages)
	}()

	// iterates over channels and merges with rides

	ridesMap := make(map[string]*models.Ride)
	for _, ride := range rides {
		ridesMap[ride.ID] = ride
	}

	for rideReviewsAverage := range chanReviewAverages {
		if ride, ok := ridesMap[rideReviewsAverage.rideID]; ok {
			ride.ReviewsAverage = rideReviewsAverage.reviewsAverage
		}
	}

	// checks if the group returned error

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	return rides, nil
}

// Store creates a new ride in the repository if a ride with the same ID
// doesn't exists already.
func (ru *RideUsecaseImpl) Store(ctx context.Context, ride *models.Ride) error {
	_, err := ru.rideRepo.GetByID(ride.ID)
	if err == nil {
		return errRideExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	ride.ID = uuid
	cleanRide(ride)
	err = validateRide(ride)
	if err != nil {
		return err
	}

	err = ru.rideRepo.Store(ride)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing ride in the repository.
func (ru *RideUsecaseImpl) Update(ctx context.Context, ride *models.Ride) error {
	_, err := ru.rideRepo.GetByID(ride.ID)
	if err != nil {
		return errRideDoesNotExists
	}

	cleanRide(ride)
	err = validateRide(ride)
	if err != nil {
		return err
	}

	err = ru.rideRepo.Update(ride)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an existing ride from the repository.
func (ru *RideUsecaseImpl) Delete(ctx context.Context, ID string) error {
	_, err := ru.rideRepo.GetByID(ID)
	if err != nil {
		return errRideDoesNotExists
	}

	err = ru.rideRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

func cleanRide(ride *models.Ride) {
	ride.ID = strings.TrimSpace(ride.ID)
	ride.Name = strings.TrimSpace(ride.Name)
	ride.Description = strings.TrimSpace(ride.Description)
	ride.MinAge = mathutil.ClampInt(ride.MinAge, 0, 200)
	ride.MinHeight = mathutil.ClampInt(ride.MinHeight, 0, 400)
	ride.Longitude = mathutil.ClampFloat64(ride.Longitude, -180, 180)
	ride.Latitude = mathutil.ClampFloat64(ride.Latitude, -90, 90)
}

func validateRide(ride *models.Ride) error {
	if len(ride.ID) <= 0 {
		return fmt.Errorf("validateRide: ID must be non-empty")
	}

	if len(ride.Name) <= 0 {
		return fmt.Errorf("validateRide: name must be non-empty")
	}

	return nil
}
