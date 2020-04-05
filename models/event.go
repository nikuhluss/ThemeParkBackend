package models

import "time"

// Event struct represents an event (rainout, maintenance, etc).
type Event struct {
	ID          string    `json:"id"`
	EmployeeID  string    `db:"employee_id" json:"employeeId"`
	EventTypeID EventType `db:"event_type_id" json:"eventTypeId"`
	EventType   string    `db:"event_type" json:"eventType"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PostedOn    time.Time `db:"posted_on" json:"postedOn"`

	Email     string `json:"email"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
}

// EventType struct represents an event type.
type EventType struct {
	ID     string `json:"id"`
	String string `db:"event_type" json:"string"`
}
