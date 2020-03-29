package models

import (
	"database/sql"
	"time"
)

// Maintenance is a struct that contains maintenance details for a ride.
type Maintenance struct {
	ID              string
	RideID          string `db:"ride_id"`
	MaintenanceType string `db:"maintenance_type"`
	Description     string
	Cost            float64
	Start           time.Time    `db:"start_datetime"`
	End             sql.NullTime `db:"end_datetime"`
	Assignees       []*User
}
