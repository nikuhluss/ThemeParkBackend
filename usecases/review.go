package usecases

import (
	"context"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// ReviewUsecase is the usecase for interacting with reviews.
type ReviewUsecase interface {
	GetByID(ctx context.Context, reviewID string) (*models.Review, error)
	Fetch(ctx context.Context) ([]*models.Review, error)
	FetchForRide(ctx context.Context, rideID string) ([]*models.Review, error)
	Store(ctx context.Context, review *models.Review) error
	Update(ctx context.Context, review *models.Review) error
	Delete(ctx context.Context, reviewID string) error
}
