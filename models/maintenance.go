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

func NewMaintenance(ID, rideID, maintenanceType, description string, cost float64, start time.Time, assignees []*User) *Maintenance{
	return &Maintenance{
		ID:               ID,
		RideID:           rideID,
		MaintenanceType:  maintenanceType,
		Description:      description,
		Cost:             cost,
		Start:            start,
		Assignees:        assignees,
	}
}