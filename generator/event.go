package generator

import (
	"time"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertEventType inserts the given event type.
func InsertEventType(execer Execer, eventType string) (string, error) {
	ID := gofakeit.UUID()

	insertEventTypeQuery := `
	INSERT INTO event_types (id, event_type)
	VALUES ($1, $2)
	`

	_, err := execer.Exec(insertEventTypeQuery, ID, eventType)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// MustInsertEventType is like InsertEvenType but panics on error.
func MustInsertEventType(mustExecer MustExecer, eventType string) string {
	return MustInsert(InsertEventType(&AsExecer{mustExecer}, eventType))
}

// InsertEventWithTitleAndTime inserts the event using the given information.
func InsertEventWithTitleAndTime(execer Execer, eventTypeID, title string, postedOn time.Time) (string, error) {

	ID := gofakeit.UUID()
	description := gofakeit.Sentence(16)

	insertEventQuery := `
	INSERT INTO events (id, event_type_id, title, description, posted_on)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := execer.Exec(insertEventQuery, ID, eventTypeID, title, description, postedOn)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// MustInsertEventWithTitleAndTime is like InsertEventWithTitleAndTime but panics on error.
func MustInsertEventWithTitleAndTime(mustExecer MustExecer, eventTypeID, title string, postedOn time.Time) string {
	return MustInsert(InsertEventWithTitleAndTime(&AsExecer{mustExecer}, eventTypeID, title, postedOn))
}
