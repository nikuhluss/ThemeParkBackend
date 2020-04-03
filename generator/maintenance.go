package generator

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertMaintenanceType inserts the given maintenance type and returns the
// generated ID.
func InsertMaintenanceType(execer Execer, maintenanceType string) (string, error) {

	ID := gofakeit.UUID()

	maintenanceTypeInsertQuery := `
	INSERT INTO maintenance_types (ID, maintenance_type)
	VALUES ($1, $2)
	`

	_, err := execer.Exec(maintenanceTypeInsertQuery, ID, maintenanceType)
	if err != nil {
		return "", nil
	}

	return ID, nil
}

// InsertMaintenanceWithStartAndEnd inserts the given maintenance job using the
// given start/end times and returns the generated ID.
func InsertMaintenanceWithStartAndEnd(execer Execer, rideID, maintenanceTypeID string, start time.Time, end sql.NullTime) (string, error) {

	ID := gofakeit.UUID()
	description := fmt.Sprintf("%s -- description", ID)
	cost := 1000.0

	maintenanceInsertQuery := `
	INSERT INTO rides_maintenance (ID, ride_id, maintenance_type_id, description, cost, start_datetime, end_datetime)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := execer.Exec(maintenanceInsertQuery, ID, rideID, maintenanceTypeID, description, cost, start, end)
	if err != nil {
		return "", nil
	}

	return ID, nil
}

// InsertMaintenance is similar to InsertMaintenanceWithStartAndEnd but generates a random
// start (end is nil).
func InsertMaintenance(execer Execer, rideID, maintenanceTypeID string) (string, error) {
	start := gofakeit.DateRange(year2000, year2000.Add(time.Hour*24*365*10))
	end := sql.NullTime{}
	return InsertMaintenanceWithStartAndEnd(execer, rideID, maintenanceTypeID, start, end)
}

// MustInsertMaintenanceType is similar to InsertMaintenanceType but panics on error.
func MustInsertMaintenanceType(mustExecer MustExecer, maintenanceType string) string {
	return MustInsert(InsertMaintenanceType(&AsExecer{mustExecer}, maintenanceType))
}

// MustInsertMaintenanceWithStartAndEnd is similar to InsterMaintenance but panics on error.
func MustInsertMaintenanceWithStartAndEnd(mustExecer MustExecer, rideID, maintenanceTypeID string, start time.Time, end sql.NullTime) string {
	return MustInsert(InsertMaintenanceWithStartAndEnd(&AsExecer{mustExecer}, rideID, maintenanceTypeID, start, end))
}

// MustInsertMaintenance is similar to InsterMaintenance but panics on error.
func MustInsertMaintenance(mustExecer MustExecer, rideID, maintenanceTypeID string) string {
	return MustInsert(InsertMaintenance(&AsExecer{mustExecer}, rideID, maintenanceTypeID))
}
