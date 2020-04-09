package impl

import (
	"context"
	"fmt"
	"strings"
	"time"

	//"golang.org/x/sync/errgroup"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories"
)

var (
	errEventExists        = fmt.Errorf("event with the given ID already exists")
	errEventDoesNotExists = fmt.Errorf("event with he given ID does not exists")
)

// EventUsecaseImpl implements the EventUsecase interface.
type EventUsecaseImpl struct {
	eventRepo repos.EventRepository
	timeout   time.Duration
}

// NewEventUsecaseImpl returns a new EventUsecaseImpl instance. The timeout
// parameter specifies a duration for each request before throwing and error.
func NewEventUsecaseImpl(
	eventRepo repos.EventRepository,
	timeout time.Duration) *EventUsecaseImpl {

	return &EventUsecaseImpl{
		eventRepo,
		timeout,
	}
}

// GetByID fetches event from the repositories using the given ID.
func (eu *EventUsecaseImpl) GetByID(ctx context.Context, ID string) (*models.Event, error) {

	event, err := eu.eventRepo.GetByID(ID)
	if err != nil {
		return nil, fmt.Errorf("error fetching event: %s", err)
	}

	return event, nil
}

// Fetch fetches all Events from the repositories.
func (eu *EventUsecaseImpl) Fetch(ctx context.Context) ([]*models.Event, error) {
	allEvent, err := eu.eventRepo.Fetch()
	if err != nil {
		return nil, err
	}

	return allEvent, nil
}

// FetchSince is like fetch, but fetches all events since a specific time.
func (eu *EventUsecaseImpl) FetchSince(ctx context.Context, day time.Time) ([]*models.Event, error) {
	event, err := eu.eventRepo.FetchSince(day)
	if err != nil {
		return nil, err
	}

	return event, err
}

// Store creates a new event in the repository if a event with the same ID
// doesn't exists already.
func (eu *EventUsecaseImpl) Store(ctx context.Context, event *models.Event) error {
	_, err := eu.eventRepo.GetByID(event.ID)
	if err == nil {
		return errEventExists
	}

	uuid, err := GenerateUUID()
	if err != nil {
		return err
	}

	event.ID = uuid
	event.PostedOn = time.Now().UTC()
	cleanEvent(event)

	err = validateEvent(event)
	if err != nil {
		return err
	}

	err = eu.eventRepo.Store(event)
	if err != nil {
		return err
	}

	return nil
}

// Update updates a specific Event job in the repositories.
func (eu *EventUsecaseImpl) Update(ctx context.Context, event *models.Event) error {
	_, err := eu.eventRepo.GetByID(event.ID)
	if err != nil {
		return errEventDoesNotExists
	}

	cleanEvent(event)
	err = validateEvent(event)
	if err != nil {
		return err
	}

	event.PostedOn = time.Time{}
	err = eu.eventRepo.Update(event)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a specific event from the repositories.
func (eu *EventUsecaseImpl) Delete(ctx context.Context, ID string) error {
	_, err := eu.eventRepo.GetByID(ID)
	if err != nil {
		return errEventDoesNotExists
	}

	err = eu.eventRepo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}

// AvailableEventTypes returns the available event types.
func (eu *EventUsecaseImpl) AvailableEventTypes(ctx context.Context) ([]*models.EventType, error) {
	etypes, err := eu.eventRepo.AvailableEventTypes()
	if err != nil {
		return nil, err
	}
	return etypes, nil
}

func cleanEvent(event *models.Event) {
	event.ID = strings.TrimSpace(event.ID)
	event.EventType = strings.TrimSpace(event.EventType)
	event.Title = strings.TrimSpace(event.Title)
	event.Description = strings.TrimSpace(event.Description)
	event.EmployeeID.String = strings.TrimSpace(event.EmployeeID.String)

	// NOTE: No need to clean values derived from other tables.
	// TODO: Documenting which are derived values and which aren't.
	// event.Email = strings.TrimSpace(event.Email)
	// event.FirstName = strings.TrimSpace(event.FirstName)
	// event.LastName = strings.TrimSpace(event.LastName)

}

func validateEvent(event *models.Event) error {

	if len(event.ID) <= 0 {
		return fmt.Errorf("validateEvent: ID must be non-empty")
	}

	if len(event.EventType) <= 0 {
		return fmt.Errorf("validateEvent: event type must be non-empty")
	}

	if len(event.Title) <= 0 {
		return fmt.Errorf("validateEvent: title must be non-empty")
	}

	if event.EmployeeID.Valid && len(event.EmployeeID.String) <= 0 {
		return fmt.Errorf("validateEvent: employee ID must be non-empty")
	}

	return nil
}
