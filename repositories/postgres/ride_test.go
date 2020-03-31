package postgres

import (
	//"database/sql"
	"fmt"
	"testing"

	//"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func RideRepositoryFixture() (*RideRepository, *sqlx.DB, func()) {
	dbconfig := testutil.NewDatabaseConnectionConfig()
	db, err := testutil.NewDatabaseConnection(dbconfig)
	if err != nil {
		panic(fmt.Sprintf("ride_test.go: RideRepositoryFixture: %s", err))
	}

	rideRepository := NewRideRepository(db)
	return rideRepository, db, func() {
		db.Close()
	}
}

func setupTestRides(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE rides CASCADE")

	rideInsertQuery := `
	INSERT INTO rides (ID, picture_collection_id, name, description, min_age, min_height, longitude, latitude)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	db.MustExec(rideInsertQuery, "ride--A", "ride--A--picId", "ride--A--name", "ride--A--description", 1, 1, 1, 1)
	db.MustExec(rideInsertQuery, "ride--B", "ride--B--picId", "ride--B--name", "ride--B--description", 2, 2, 2, 2)
	db.MustExec(rideInsertQuery, "ride--C", "ride--C--picId", "ride--C--name", "ride--C--description", 3, 3, 3, 3)

}

// Tests
// --------------------------------

func TestGetRidesByIDSucceeds(t *testing.T) {
	rideRepository, db, teardown := RideRepositoryFixture()
	defer teardown()

	setupTestRides(db)
	tests := []struct {
		rideID string
	}{
		{"ride--A"},
		{"ride--B"},
		{"ride--C"},
	}

	for idx, tt := range tests {
		ride, err := rideRepository.GetByID(tt.rideID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, tt.rideID, ride.ID)
		assert.Equal(t, tt.rideID+"--name", ride.Name)
		assert.Equal(t, tt.rideID+"--description", ride.Description)
		assert.Equal(t, int(idx+1), int(ride.MinAge))
		assert.Equal(t, int(idx+1), int(ride.MinHeight))
		assert.Equal(t, int(idx+1), int(ride.Longitude))
		assert.Equal(t, int(idx+1), int(ride.Latitude))
	}
}

func TestGetByRideIDNoMatchFails(t *testing.T) {
	rideRepository, _, teardown := RideRepositoryFixture()
	defer teardown()

	ride, err := rideRepository.GetByID("some-unknown-ID")
	assert.Nil(t, ride)
	assert.NotNil(t, err)
}

func TestRideFetchSucceeds(t *testing.T) {
	rideRepository, db, teardown := RideRepositoryFixture()
	defer teardown()

	setupTestRides(db)

	rides, err := rideRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, rides, 3)
}

func TestStoreRideSucceeds(t *testing.T) {
	rideRepository, _, teardown := RideRepositoryFixture()
	defer teardown()

	ride := models.NewRide("ride--D", "ride--D--name", "ride--D--description", 4, 4, 4, 4)
	err := rideRepository.Store(ride)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ride, err = rideRepository.GetByID(ride.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, ride)
	assert.Equal(t, "ride--D", ride.ID)

	assert.Equal(t, "ride--D--name", ride.Name)
	assert.Equal(t, "ride--D--description", ride.Description)

}

func TestUpdateRideSucceeds(t *testing.T) {
	rideRepository, db, teardown := RideRepositoryFixture()
	defer teardown()

	setupTestRides(db)

	ride, err := rideRepository.GetByID("ride--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedRide := models.NewRide(ride.ID, "ride--D--name", "ride--D--description", 4, 4, 4, 4)

	err = rideRepository.Update(expectedRide)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedRide, err := rideRepository.GetByID("ride--A")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedRide.ID, updatedRide.ID)
	assert.Equal(t, expectedRide.Name, updatedRide.Name)
	assert.Equal(t, expectedRide.Description, updatedRide.Description)
	assert.Equal(t, expectedRide.MinAge, updatedRide.MinAge)
	assert.Equal(t, expectedRide.MinHeight, updatedRide.MinHeight)

}

func TestDeleteRideSucceeds(t *testing.T) {
	rideRepository, db, teardown := RideRepositoryFixture()
	defer teardown()

	setupTestRides(db)

	err := rideRepository.Delete("ride--C")
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	ride, err := rideRepository.GetByID("ride--C")
	assert.Nil(t, ride)
	assert.NotNil(t, err)
}
