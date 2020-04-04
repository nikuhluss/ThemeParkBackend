package postgres_test

import (
	"database/sql"
	"testing"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"

	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func setupTestMaintenance(db *sqlx.DB) ([]string, []string) {

	tx := db.MustBegin()

	tx.MustExec("TRUNCATE TABLE rides CASCADE")
	tx.MustExec("TRUNCATE TABLE maintenance_types CASCADE")
	tx.MustExec("TRUNCATE TABLE rides_maintenance CASCADE")

	ride0 := generator.MustInsertRide(tx)
	ride1 := generator.MustInsertRide(tx)

	maintenanceTypeTuneUp := generator.MustInsertMaintenanceType(tx, "Tune Up")
	maintenanceTypeReplacement := generator.MustInsertMaintenanceType(tx, "Replacement")
	maintenanceTypeFixed := generator.MustInsertMaintenanceType(tx, "Fixed")

	maintenance0 := generator.MustInsertMaintenance(tx, ride0, maintenanceTypeTuneUp)
	maintenance1 := generator.MustInsertMaintenance(tx, ride1, maintenanceTypeReplacement)
	maintenance2 := generator.MustInsertMaintenance(tx, ride1, maintenanceTypeFixed)

	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return []string{ride0, ride1}, []string{maintenance0, maintenance1, maintenance2}
}

// Tests
// --------------------------------
func TestMaintenanceGetByIDSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	_, maintenanceIDs := setupTestMaintenance(db)

	for _, maintenanceID := range maintenanceIDs {
		maintenance, err := maintenanceRepository.GetByID(maintenanceID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, maintenanceID, maintenance.ID)
		assert.NotEmpty(t, maintenance.RideID)
		assert.NotEmpty(t, maintenance.RideName)
		assert.NotEmpty(t, maintenance.MaintenanceType)
		assert.NotEmpty(t, maintenance.Description)
		assert.Equal(t, float64(1000), maintenance.Cost)
	}
}

func TestMaintenanceFetchSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	maintenance, err := maintenanceRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, maintenance, 3)
}

func TestMaintenanceFetchForRideSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	rideIDs, _ := setupTestMaintenance(db)
	rideID := rideIDs[0]

	maintenance, err := maintenanceRepository.FetchForRide(rideID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, maintenance, 1)
}

func TestMaintenanceStoreSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	rideIDs, _ := setupTestMaintenance(db)
	rideID := rideIDs[0]

	users := []*models.User{
		models.NewEmployee("user--D", "user--D--email", "user--D--passS", "user--D--passH", "Ride Manager", 22),
	}

	maintenance := models.NewMaintenance("maintenance--ID", rideID, "Ride name", "Tune Up", "description", 60, time.Now(), users)
	err := maintenanceRepository.Store(maintenance)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	maintenanceOut, err := maintenanceRepository.GetByID(maintenance.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, maintenance.ID, maintenanceOut.ID)
	assert.Equal(t, maintenance.RideID, maintenanceOut.RideID)
	assert.Equal(t, maintenance.MaintenanceType, maintenanceOut.MaintenanceType)
	assert.Equal(t, maintenance.Description, maintenance.Description)
}

func TestMaintenanceUpdateSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	rideIDs, maintenanceIDs := setupTestMaintenance(db)
	rideID := rideIDs[len(rideIDs)-1]
	maintenanceID := maintenanceIDs[0]

	maintenance, err := maintenanceRepository.GetByID(maintenanceID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	users := []*models.User{
		models.NewEmployee("user--D", "user--D--email", "user--D--passS", "user--D--passH", "Ride Manager", 22),
	}

	expectedMaintenance := models.NewMaintenance(maintenance.ID, rideID, "new name", "Replacement", "new description", 70, maintenance.Start, users)
	expectedMaintenance.End = sql.NullTime{Time: time.Now(), Valid: true}

	err = maintenanceRepository.Update(expectedMaintenance)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedMaintenance, err := maintenanceRepository.GetByID(maintenanceID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedMaintenance.ID, updatedMaintenance.ID)
	assert.Equal(t, expectedMaintenance.RideID, updatedMaintenance.RideID)
	assert.Equal(t, expectedMaintenance.MaintenanceType, updatedMaintenance.MaintenanceType)
	assert.Equal(t, expectedMaintenance.Description, updatedMaintenance.Description)
	assert.Equal(t, expectedMaintenance.Cost, updatedMaintenance.Cost)

}

func TestMaintenanceDeleteSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	_, maintenanceIDs := setupTestMaintenance(db)
	maintenanceID := maintenanceIDs[0]

	err := maintenanceRepository.Delete(maintenanceID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	maintenance, err := maintenanceRepository.GetByID(maintenanceID)
	assert.Nil(t, maintenance)
	assert.NotNil(t, err)
}

func TestMaintenanceGetAllMaintenanceTypesSucceeds(t *testing.T) {
	maintenanceRepository, db, teardown := testutil.MakeMaintenanceRepositoryFixture()
	defer teardown()

	setupTestMaintenance(db)

	maintenanceTypes, err := maintenanceRepository.AvailableMaintenanceTypes()

	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, "Fixed", maintenanceTypes[0])
	assert.Equal(t, "Replacement", maintenanceTypes[1])
	assert.Equal(t, "Tune Up", maintenanceTypes[2])
}
