package postgres

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

var selectEvents = psql.
	Select("events.*", "event_types.*", "users.email", "user_details.first_name", "user_details.last_name").
	From("events").
	Join("event_types ON event_types.ID = events.event_type_id").
	Join("users ON users.id = events.employee_id").
	LeftJoin("user_details ON user_details.user_id = events.employee_id").
	OrderBy("events.posted_on")

// EventRepository implements the EventRepository interface for postgres.
type EventRepository struct {
	db *sqlx.DB
}

// NewEventRepository creates a new EventRepository instance.
func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{db}
}

// GetByID fetches an event using the given ID.
func (er *EventRepository) GetByID(ID string) (*models.Event, error) {
	db := er.db
	udb := db.Unsafe()

	query, _ := selectEvents.Where(sq.Eq{"events.id": "$1"}).MustSql()

	event := models.Event{}
	err := udb.Get(&event, query, ID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

// Fetch fetches all events.
func (er *EventRepository) Fetch() ([]*models.Event, error) {
	db := er.db
	udb := db.Unsafe()

	query, _ := selectEvents.MustSql()

	events := []*models.Event{}
	err := udb.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// FetchSince fetches all events since given time.
func (er *EventRepository) FetchSince(since time.Time) ([]*models.Event, error) {
	db := er.db
	udb := db.Unsafe()

	query, _ := selectEvents.Where(sq.GtOrEq{"events.posted_on": "$1"}).MustSql()

	events := []*models.Event{}
	err := udb.Select(&events, query, since)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Store creates a new event.
func (er *EventRepository) Store(event *models.Event) error {
	db := er.db

	query, _, _ := psql.
		Insert("events").
		Columns("id", "employee_id", "event_type_id", "title", "description", "posted_on").
		Values("$1", "$2", "$3", "$4", "$5", "$6").
		ToSql()

	_, err := db.Exec(query, event.ID, event.EmployeeID, event.EventTypeID, event.Title, event.Description, event.PostedOn)
	if err != nil {
		return err
	}

	return nil
}

// Update updates an existing event.
func (er *EventRepository) Update(event *models.Event) error {
	db := er.db

	query, _, _ := psql.
		Update("events").
		Set("employee_id", "$1").
		Set("event_type_id", "$2").
		Set("title", "$3").
		Set("description", "$4").
		Set("posted_on", "$5").
		Where(sq.Eq{"id": "$6"}).
		ToSql()

	_, err := db.Exec(query, event.EmployeeID, event.EventTypeID, event.Title, event.Description, event.PostedOn, event.ID)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes an existing event.
func (er *EventRepository) Delete(eventID string) error {
	db := er.db

	query, _, _ := psql.Delete("events").Where(sq.Eq{"id": "$1"}).ToSql()

	_, err := db.Exec(query, eventID)
	if err != nil {
		return err
	}

	return nil
}

// AvailableEventTypes returns the available event types.
func (er *EventRepository) AvailableEventTypes() ([]*models.EventType, error) {
	db := er.db

	query, _ := psql.Select("event_types.*").From("event_types").OrderBy("event_type ASC").MustSql()

	eventTypes := []*models.EventType{}
	err := db.Select(&eventTypes, query)
	if err != nil {
		return nil, err
	}

	return eventTypes, nil
}
