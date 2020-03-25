package impl

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repositories "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

// RideUsecaseImpl implements the RideUsecase interface.
type RideUsecaseImpl struct {
	rideRepo    repositories.RideRepository
	reviewRepo  repositories.ReviewRepository
	pictureRepo repositories.PictureRepository
	timeout     time.Duration
}

// NewRideUsecaseImpl returns a new RideUsecaseImpl instance. The timeout
// parameter specifies a duration for each request before throwing and error.
func NewRideUsecaseImpl(
	rideRepo repositories.RideRepository,
	reviewRepo repositories.ReviewRepository,
	pictureRepo repositories.PictureRepository,
	timeout time.Duration) *RideUsecaseImpl {

	return &RideUsecaseImpl{
		rideRepo,
		reviewRepo,
		pictureRepo,
		timeout,
	}
}

// GetByID fetches ride from the repositories using the given ID.
func (ru *RideUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.Ride, error) {

	ride, err := ru.rideRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride: %s", err)
	}

	reviews, err := ru.reviewRepo.FetchForRideSortedByDate(ride.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride reviews: %s", err)
	}

	pictures, err := ru.pictureRepo.FetchByCollectionID(ride.ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching ride pictures: %s", err)
	}

	ride.Reviews = reviews
	ride.Pictures = pictures
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

	// the following section fetches reviews and pictures in parallel

	type RideReviews struct {
		rideID  string
		reviews []*models.Review
	}
	type RidePictures struct {
		rideID   string
		pictures []*models.Picture
	}

	chanReviews := make(chan RideReviews)
	chanPictures := make(chan RidePictures)

	for _, ride := range rides {
		rideID := ride.ID

		eg.Go(func() error {
			reviews, err := ru.reviewRepo.FetchForRideSortedByDate(rideID)
			if err != nil {
				return err
			}

			chanReviews <- RideReviews{rideID, reviews}
			return nil
		})

		eg.Go(func() error {
			pictures, err := ru.pictureRepo.FetchByCollectionID(rideID)
			if err != nil {
				return err
			}

			chanPictures <- RidePictures{rideID, pictures}
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
		close(chanReviews)
		close(chanPictures)
	}()

	err = eg.Wait()
	if err != nil {
		return nil, err
	}

	// iterates over channels and marges rides, reviews channel, and pictures channel

	ridesMap := make(map[string]*models.Ride)
	for _, ride := range rides {
		ridesMap[ride.ID] = ride
	}

	reviewsOk := true
	picturesOk := true

	for reviewsOk || picturesOk {
		select {
		case rideReviews, reviewsOk := <-chanReviews:
			if !reviewsOk {
				continue
			}

			if ride, ok := ridesMap[rideReviews.rideID]; ok {
				ride.Reviews = rideReviews.reviews
			}

		case ridePictures, picturesOk := <-chanPictures:
			if !picturesOk {
				continue
			}

			if ride, ok := ridesMap[ridePictures.rideID]; ok {
				ride.Pictures = ridePictures.pictures
			}
		}
	}

	return rides, nil
}

// Store creates a new ride in the repository if a ride with the same ID
// doesn't exists already.
func (ru *RideUsecaseImpl) Store(ctx context.Context, ride *models.Ride) error {
	return nil
}

// Update updates an existing ride in the repository.
func (ru *RideUsecaseImpl) Update(ctx context.Context, ride *models.Ride) error {
	return nil
}

// Delete deletes an existing ride from the repository.
func (ru *RideUsecaseImpl) Delete(ctx context.Context, ID string) error {
	return nil
}
