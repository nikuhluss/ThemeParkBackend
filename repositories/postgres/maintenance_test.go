package postgres_test

import (
	//"database/sql"
	"fmt"
	"testing"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
)

// Fixtures
// --------------------------------

func setupTestMaintenance(db *sqlx.DB) {

	rideInsertQuery := `
	INSERT INTO rides (ID, picture_collection_id, name, description, min_age, min_height, longitude, latitude)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	db.MustExec("TRUNCATE TABLE rides CASCADE")
	db.MustExec(rideInsertQuery, "maintenance--A--ride", "ride--A--picId", "ride--A--name", "ride--A--description", 1, 1, 1, 1)
	db.MustExec(rideInsertQuery, "maintenance--B--ride", "ride--B--picId", "ride--B--name", "ride--B--description", 2, 2, 2, 2)
	db.MustExec(rideInsertQuery, "maintenance--C--ride", "ride--C--picId", "ride--C--name", "ride--C--description", 3, 3, 3, 3)

	MtypeInsertQuery := `
	INSERT INTO maintenance_types (ID, maintenance_type)
	VALUES ($1, $2)
	`

	db.MustExec("TRUNCATE TABLE maintenance_types CASCADE")
	db.MustExec(MtypeInsertQuery, "type--A", "Tune up")
	db.MustExec(MtypeInsertQuery, "type--B", "Replacement")
	db.MustExec(MtypeInsertQuery, "type--C", "Fixed")

	MaintenanceInsertQuery := `
	INSERT INTO rides_maintenance (ID, ride_id, maintenance_type_id, description, cost, start_datetime, end_datetime)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	db.MustExec("TRUNCATE TABLE rides_maintenance CASCADE")
	db.MustExec(MaintenanceInsertQuery, "maintenance--A", "maintenance--A--ride", "type--A", "maintenance--A--description", 1, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--B", "maintenance--B--ride", "type--B", "maintenance--B--description", 2, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--C", "maintenance--C--ride", "type--C", "maintenance--C--description", 3, time.Now(), time.Now())

}

// Tests
// --------------------------------
func TestGetMaintenanceByIDSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	tests := []struct {
		maintenanceID           string
		expectedMaintenanceType string
	}{
		{"maintenance--A", "Tune up"},
		{"maintenance--B", "Replacement"},
		{"maintenance--C", "Fixed"},
	}

	for idx, tt := range tests {
		maintenance, err := maintenanceRepository.GetByID(tt.maintenanceID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		fmt.Println(">", maintenance)

		assert.Equal(t, tt.maintenanceID, maintenance.ID)
		assert.Equal(t, tt.maintenanceID+"--ride", maintenance.RideID)
		assert.Equal(t, tt.expectedMaintenanceType, maintenance.MaintenanceType)
		assert.Equal(t, tt.maintenanceID+"--description", maintenance.Description)
		assert.Equal(t, int(idx+1), int(maintenance.Cost))
	}
}
