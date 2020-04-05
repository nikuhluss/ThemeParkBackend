package repositories

import (
	"time"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// EventRepository defines the interface for interacting with events.
type EventRepository interface {
	GetByID(ID string) (*models.Event, error)
	Fetch() ([]*models.Event, error)
	FetchForDay(day time.Time) ([]*models.Event, error)
	Store(event *models.Event) error
	Update(event *models.Event) error
	Delete(eventID string) error
	AvailableEventTypes() ([]*models.EventType, error)
}
