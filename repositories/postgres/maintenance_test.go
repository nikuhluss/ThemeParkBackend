package postgres

import (
	//"database/sql"
	"fmt"
	"testing"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	//"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func MaintenanceRepositoryFixture() (*MaintenanceRepository, *sqlx.DB, func()) {
	dbconfig := testutil.NewDatabaseConnectionConfig()
	db, err := testutil.NewDatabaseConnection(dbconfig)
	if err != nil {
		panic(fmt.Sprintf("maintenance_test.go: MaintenanceRepositoryFixture: %s", err))
	}

	maintenanceRepository := NewMaintenanceRepository(db)
	return maintenanceRepository, db, func() {
		db.Close()
	}
}

func setupTestMaintenance(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE rides CASCADE")

	rideInsertQuery := `
	INSERT INTO rides (ID, picture_collection_id, name, description, min_age, min_height, longitude, latitude)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	db.MustExec(rideInsertQuery, "ride--A", "ride--A--picId", "ride--A--name", "ride--A--description", 1, 1, 1, 1)
	db.MustExec(rideInsertQuery, "ride--B", "ride--B--picId", "ride--B--name", "ride--B--description", 2, 2, 2, 2)
	db.MustExec(rideInsertQuery, "ride--C", "ride--C--picId", "ride--C--name", "ride--C--description", 3, 3, 3, 3)
	
	db.MustExec("TRUNCATE TABLE maintenance_types CASCADE")
	MtypeInsertQuery :=`
	INSERT INTO maintenance_types (ID, maintenance_type)
	VALUES ($1, $2)
	`
	db.MustExec(MtypeInsertQuery,"type--A", "tune_up")
	db.MustExec(MtypeInsertQuery,"type--B", "replacement")
	db.MustExec(MtypeInsertQuery,"type--C", "fixed")

	db.MustExec("TRUNCATE TABLE rides_maintenance CASCADE")

	MaintenanceInsertQuery := `
	INSERT INTO rides_maintenance (ID, ride_id, maintenance_type_id, description, cost, start_datetime, end_datetime)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	db.MustExec(MaintenanceInsertQuery, "maintenance--A", "ride--A", "type--A", "maintenance--A--description", 1, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--B", "ride--B", "type--B", "maintenance--B--description", 2, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--C", "ride--C", "type--C", "maintenance--C--description", 3, time.Now(), time.Now())

}

// Tests
// --------------------------------
func TestGetMaintenanceByIDSucceeds(t *testing.T){
	maintenanceRepository, db, teardown := MaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)
	tests := []struct {
		maintenanceID string
	}{
		{"maintenance--A"},
		{"maintenance--B"},
		{"maintenance--C"},
	}
	for idx, tt := range tests {
		maintenance, err := maintenanceRepository.GetByID(tt.maintenanceID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.maintenanceID, maintenance.ID)
		assert.Equal(t, tt.maintenanceID+"--ride_id", maintenance.RideID)
		assert.Equal(t, tt.maintenanceID+"--type_id", maintenance.MaintenanceTypeID)
		assert.Equal(t, tt.maintenanceID+"--type_id", maintenance.Description)
		assert.Equal(t, int(idx+1), int(maintenance.Cost))
	}
}