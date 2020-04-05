package models

import (
	"time"
)

// Maintenance is a struct that contains maintenance details for a ride.
type Maintenance struct {
	ID              string    `json:"id"`
	RideID          string    `db:"ride_id" json:"rideId"`
	RideName        string    `db:"ride_name" json:"rideName"`
	MaintenanceType string    `db:"maintenance_type" json:"maintenanceType"`
	Description     string    `json:"description"`
	Cost            float64   `json:"cost"`
	Start           time.Time `db:"start_datetime" json:"start"`
	End             NullTime  `db:"end_datetime" json:"end"`
	Assignees       []*User   `json:"assignees"`
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
