package postgres_test

import (
	"database/sql"
	"fmt"
	"testing"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func setupTestMaintenance(db *sqlx.DB) {

	rideInsertQuery := `
	INSERT INTO rides (ID, picture_collection_id, name, description, min_age, min_height, longitude, latitude)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	db.MustExec("TRUNCATE TABLE rides CASCADE")
	db.MustExec(rideInsertQuery, "maintenance--A--ride-id", "ride--A--picId", "maintenance--A--ride-name", "ride--A--description", 1, 1, 1, 1)
	db.MustExec(rideInsertQuery, "maintenance--B--ride-id", "ride--B--picId", "maintenance--B--ride-name", "ride--B--description", 2, 2, 2, 2)
	db.MustExec(rideInsertQuery, "maintenance--C--ride-id", "ride--C--picId", "maintenance--C--ride-name", "ride--C--description", 3, 3, 3, 3)

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
	db.MustExec(MaintenanceInsertQuery, "maintenance--A", "maintenance--A--ride-id", "type--A", "maintenance--A--description", 1, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--B", "maintenance--B--ride-id", "type--B", "maintenance--B--description", 2, time.Now(), time.Now())
	db.MustExec(MaintenanceInsertQuery, "maintenance--C", "maintenance--C--ride-id", "type--C", "maintenance--C--description", 3, time.Now(), time.Now())

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
		assert.Equal(t, tt.maintenanceID+"--ride-id", maintenance.RideID)
		assert.Equal(t, tt.maintenanceID+"--ride-name", maintenance.RideName)
		assert.Equal(t, tt.expectedMaintenanceType, maintenance.MaintenanceType)
		assert.Equal(t, tt.maintenanceID+"--description", maintenance.Description)
		assert.Equal(t, int(idx+1), int(maintenance.Cost))
	}
}

func TestFetchMaintenanceSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	maintenance, err := maintenanceRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, maintenance, 3)
}

func TestFetchForRideSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)
	maintenance, err := maintenanceRepository.FetchForRide("maintenance--A--ride-id")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, maintenance, 1)
}

func TestStoreMaintenanceSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	users := []*models.User{
		models.NewEmployee("user--D", "user--D--email", "user--D--passS", "user--D--passH", "Ride Manager", 22),
	}

	maintenance := models.NewMaintenance("maintenance--D", "maintenance--A--ride-id", "maintenance--A--ride-name", "Tune up", "maintenance--D--description", 60, time.Now(), users)
	err := maintenanceRepository.Store(maintenance)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	maintenanceOut, err := maintenanceRepository.GetByID(maintenance.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, maintenanceOut)

}

func TestUpdateMaintenanceSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	maintenance, err := maintenanceRepository.GetByID("maintenance--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	users := []*models.User{
		models.NewEmployee("user--D", "user--D--email", "user--D--passS", "user--D--passH", "Ride Manager", 22),
	}

	expectedMaintenance := models.NewMaintenance(maintenance.ID, "maintenance--B--ride-id", "maintenance--B--ride-name", "Replacement", "maintenance--A--new Description", 70, time.Now(), users)
	var p sql.NullTime
	p.Time = time.Now()
	expectedMaintenance.End = p

	err = maintenanceRepository.Update(expectedMaintenance)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedMaintenance, err := maintenanceRepository.GetByID("maintenance--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedMaintenance.ID, updatedMaintenance.ID)
	assert.Equal(t, expectedMaintenance.RideID, updatedMaintenance.RideID)
	assert.Equal(t, expectedMaintenance.MaintenanceType, updatedMaintenance.MaintenanceType)
	assert.Equal(t, expectedMaintenance.Description, updatedMaintenance.Description)
	assert.Equal(t, expectedMaintenance.Cost, updatedMaintenance.Cost)

}

func TestDeleteMaintenanceSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	err := maintenanceRepository.Delete("maintenance--C")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	maintenance, err := maintenanceRepository.GetByID("maintenance--C")
	assert.Nil(t, maintenance)
	assert.NotNil(t, err)
}

func TestGetAllMaintenanceTypesSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	maintenanceTypes, err := maintenanceRepository.AvailableMaintenanceTypes()

	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "Fixed", maintenanceTypes[0])
	assert.Equal(t, "Replacement", maintenanceTypes[1])
	assert.Equal(t, "Tune up", maintenanceTypes[2])
}
