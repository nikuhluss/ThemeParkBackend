package models

import "time"

// Event struct represents an event (rainout, maintenance, etc).
type Event struct {
	ID          string     `json:"id"`
	EventTypeID string     `db:"event_type_id" json:"eventTypeId"`
	EventType   string     `db:"event_type" json:"eventType"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PostedOn    time.Time  `db:"posted_on" json:"postedOn"`
	EmployeeID  NullString `db:"employee_id" json:"employeeId"`

	Email     NullString `json:"email"`
	FirstName NullString `db:"first_name" json:"firstName"`
	LastName  NullString `db:"last_name" json:"lastName"`
}

// EventType struct represents an event type.
type EventType struct {
	ID     string `json:"id"`
	String string `db:"event_type" json:"string"`
}
