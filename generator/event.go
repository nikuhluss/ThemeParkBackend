package generator

import "github.com/brianvoe/gofakeit/v4"

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
