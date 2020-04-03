package postgres_test

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"
	//"database/sql"

	"testing"

	//"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func setupTestRides(db *sqlx.DB) []string {

	rideIDs := make([]string, 0, 3)

	tx := db.MustBegin()
	tx.MustExec("TRUNCATE TABLE rides CASCADE")
	rideIDs = append(rideIDs, generator.MustInsertRide(tx))
	rideIDs = append(rideIDs, generator.MustInsertRide(tx))
	rideIDs = append(rideIDs, generator.MustInsertRide(tx))
	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return rideIDs
}

// Tests
// --------------------------------

func TestGetRidesByIDSucceeds(t *testing.T) {
	rideRepository, db, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	tests := setupTestRides(db)

	for _, rideID := range tests {
		ride, err := rideRepository.GetByID(rideID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, rideID, ride.ID)
		assert.Equal(t, rideID+" -- name", ride.Name)
		assert.Equal(t, rideID+" -- description", ride.Description)
		assert.Equal(t, 1, ride.MinAge)
		assert.Equal(t, 2, ride.MinHeight)
		assert.Equal(t, float64(3.0), ride.Longitude)
		assert.Equal(t, float64(4.0), ride.Latitude)
	}
}

func TestGetByRideIDNoMatchFails(t *testing.T) {
	rideRepository, _, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	ride, err := rideRepository.GetByID("some-unknown-ID")
	assert.Nil(t, ride)
	assert.NotNil(t, err)
}

func TestRideFetchSucceeds(t *testing.T) {
	rideRepository, db, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	setupTestRides(db)

	rides, err := rideRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, rides, 3)
}

func TestStoreRideSucceeds(t *testing.T) {
	rideRepository, _, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	expectedRide := models.NewRide("ride--ID", "ride--ID--name", "ride--ID--description", 1, 2, 3, 4)
	err := rideRepository.Store(expectedRide)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ride, err := rideRepository.GetByID(expectedRide.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, ride)
	assert.Equal(t, expectedRide.ID, ride.ID)
	assert.Equal(t, expectedRide.Name, ride.Name)
	assert.Equal(t, expectedRide.Description, ride.Description)
	assert.Equal(t, expectedRide.MinAge, ride.MinAge)
	assert.Equal(t, expectedRide.MinHeight, ride.MinHeight)
	assert.Equal(t, expectedRide.Longitude, ride.Longitude)
	assert.Equal(t, expectedRide.Latitude, ride.Latitude)

}

func TestUpdateRideSucceeds(t *testing.T) {
	rideRepository, db, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	tests := setupTestRides(db)
	rideID := tests[0]

	ride, err := rideRepository.GetByID(rideID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedRide := models.NewRide(ride.ID, "new name", "new description", 4, 4, 4, 4)
	err = rideRepository.Update(expectedRide)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedRide, err := rideRepository.GetByID(rideID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedRide.ID, updatedRide.ID)
	assert.Equal(t, expectedRide.Name, updatedRide.Name)
	assert.Equal(t, expectedRide.Description, updatedRide.Description)
	assert.Equal(t, expectedRide.MinAge, updatedRide.MinAge)
	assert.Equal(t, expectedRide.MinHeight, updatedRide.MinHeight)
	assert.Equal(t, expectedRide.Longitude, updatedRide.Longitude)
	assert.Equal(t, expectedRide.Latitude, updatedRide.Latitude)
}

func TestDeleteRideSucceeds(t *testing.T) {
	rideRepository, db, teardown := testutil.MakeRideRepositoryFixture()
	defer teardown()

	tests := setupTestRides(db)
	rideID := tests[0]

	err := rideRepository.Delete(rideID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ride, err := rideRepository.GetByID(rideID)
	assert.Nil(t, ride)
	assert.NotNil(t, err)
}
