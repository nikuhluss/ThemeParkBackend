package models

import (
	"database/sql"
	"time"
)

// Maintenance is a struct that contains maintenance details for a ride.
type Maintenance struct {
	ID              string
	RideID          string `db:"ride_id"`
	RideName        string `db:"ride_name"`
	MaintenanceType string `db:"maintenance_type"`
	Description     string
	Cost            float64
	Start           time.Time    `db:"start_datetime"`
	End             sql.NullTime `db:"end_datetime"`
	Assignees       []*User
}

// NewMaintenance returns a new Maintenance instance.
func NewMaintenance(ID, rideID, rideName, maintenanceType, description string, cost float64, start time.Time, assignees []*User) *Maintenance {
	return &Maintenance{
		ID:              ID,
		RideID:          rideID,
		RideName:        rideName,
		MaintenanceType: maintenanceType,
		Description:     description,
		Cost:            cost,
		Start:           start,
		Assignees:       assignees,
	}
}
