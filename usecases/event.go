package usecases

import (
	"context"
	"time"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// EventUsecase is the usecase for interacting with events.
type EventUsecase interface {
	GetByID(ctx context.Context, ID string) (*models.Event, error)
	Fetch(ctx context.Context) ([]*models.Event, error)
	FetchSince(ctx context.Context, since time.Time) ([]*models.Event, error)
	Store(ctx context.Context, event *models.Event) error
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, ID string) error
	AvailableEventTypes(ctx context.Context) ([]*models.EventType, error)
}
