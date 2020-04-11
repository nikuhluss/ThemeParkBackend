package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/mathutil"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errReviewExists        = fmt.Errorf("review with the given ID alredy exists")
	errReviewDoesNotExists = fmt.Errorf("review with the given ID does not exists")
)

// ReviewUsecaseImpl implements the ReviewUsecase interface.
type ReviewUsecaseImpl struct {
	reviewRepo repos.ReviewRepository
	rideRepo   repos.RideRepository
	timeout    time.Duration
}

// NewReviewUsecaseImpl returns a new ReviewUsecaseImpl instance.
func NewReviewUsecaseImpl(reviewRepo repos.ReviewRepository, rideRepo repos.RideRepository, timeout time.Duration) *ReviewUsecaseImpl {
	return &ReviewUsecaseImpl{reviewRepo, rideRepo, timeout}
}

// GetByID returns a spcific review using the given ID.
func (ru *ReviewUsecaseImpl) GetByID(ctx context.Context, reviewID string) (*models.Review, error) {
	return ru.reviewRepo.GetByID(reviewID)
}

// Fetch fetches all the reviews from the repository.
func (ru *ReviewUsecaseImpl) Fetch(ctx context.Context) ([]*models.Review, error) {
	return ru.reviewRepo.Fetch()
}

// FetchForRide fetches all reviews for the given ride.
func (ru *ReviewUsecaseImpl) FetchForRide(ctx context.Context, rideID string) ([]*models.Review, error) {
	_, err := ru.rideRepo.GetByID(rideID)
	if err != nil {
		return nil, errRideDoesNotExists
	}
	return ru.reviewRepo.FetchForRideSortedByDate(rideID)
}

// Store creates a new review.
func (ru *ReviewUsecaseImpl) Store(ctx context.Context, review *models.Review) error {
	_, err := ru.reviewRepo.GetByID(review.ID)
	if err == nil {
		return errReviewExists
	}

	ID, err := GenerateUUID()
	if err != nil {
		return err
	}

	review.ID = ID
	review.PostedOn = time.Now().UTC()
	cleanReview(review)
	err = validateReview(review)
	if err != nil {
		return err
	}

	err = ru.reviewRepo.Store(review)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing review.
func (ru *ReviewUsecaseImpl) Update(ctx context.Context, review *models.Review) error {
	_, err := ru.reviewRepo.GetByID(review.ID)
	if err != nil {
		return errReviewDoesNotExists
	}

	cleanReview(review)
	err = validateReview(review)
	if err != nil {
		return err
	}

	err = ru.reviewRepo.Update(review)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a specific review.
func (ru *ReviewUsecaseImpl) Delete(ctx context.Context, reviewID string) error {
	_, err := ru.reviewRepo.GetByID(reviewID)
	if err != nil {
		return errReviewDoesNotExists
	}
	return ru.reviewRepo.Delete(reviewID)
}

func cleanReview(review *models.Review) {
	review.ID = strings.TrimSpace(review.ID)
	review.RideID = strings.TrimSpace(review.RideID)
	review.UserID = strings.TrimSpace(review.UserID)
	review.Rating = mathutil.ClampInt(review.Rating, 1, 5)
	review.Title = strings.TrimSpace(review.Title)
	review.Content = strings.TrimSpace(review.Content)
}

func validateReview(review *models.Review) error {

	if len(review.ID) <= 0 {
		return fmt.Errorf("validateReview: ID must be non-empty")
	}

	if len(review.RideID) <= 0 {
		return fmt.Errorf("validateReview: RideID must be non-empty")
	}

	if len(review.UserID) <= 0 {
		return fmt.Errorf("validateReview: UserID must be non-empty")
	}

	if !(review.Rating >= 1 && review.Rating <= 5) {
		return fmt.Errorf("validateReview: Rating must be in the range [1, 5]")
	}

	if len(review.Title) <= 0 {
		return fmt.Errorf("validateReview: Title must be non-empty")
	}

	if len(review.Content) <= 0 {
		return fmt.Errorf("validateReview: Content must be non-empty")
	}

	return nil
}
