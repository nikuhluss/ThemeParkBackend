package testutil

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	repos "gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/repositories/postgres"
)

func MakeDatabaseFixture() (*sqlx.DB, func()) {
	dbconfig := NewDatabaseConnectionConfig()

	db, err := NewDatabaseConnection(dbconfig)
	if err != nil {
		panic(fmt.Sprintf("fixtures.go: MakeDatabaseFixture: %s", err))
	}

	return db, func() {
		db.Close()
	}
}

// Make*RepositoryFixture
// --------------------------------

func MakeUserRepositoryFixture() (*repos.UserRepository, *sqlx.DB, func()) {
	db, dbTeardown := MakeDatabaseFixture()
	userRepository := repos.NewUserRepository(db)
	return userRepository, db, func() {
		dbTeardown()
	}
}

func MakeRideRepositoryFixture() (*repos.RideRepository, *sqlx.DB, func()) {
	db, dbTeardown := MakeDatabaseFixture()
	rideRepository := repos.NewRideRepository(db)
	return rideRepository, db, func() {
		dbTeardown()
	}
}

func MakePictureRepositoryFixture() (*repos.PictureRepository, *sqlx.DB, func()) {
	db, dbTeardown := MakeDatabaseFixture()
	pictureRepository := repos.NewPictureRepository(db)
	return pictureRepository, db, func() {
		dbTeardown()
	}
}

func MakeReviewRepositoryFixture() (*repos.ReviewRepository, *sqlx.DB, func()) {
	db, dbTeardown := MakeDatabaseFixture()
	reviewRepository := repos.NewReviewRepository(db)
	return reviewRepository, db, func() {
		dbTeardown()
	}
}

func MakeMaintenanceRepositoryFixture() (*repos.MaintenanceRepository, *sqlx.DB, func()) {
	db, dbTeardown := MakeDatabaseFixture()
	maintenanceRepository := repos.NewMaintenanceRepository(db)
	return maintenanceRepository, db, func() {
		dbTeardown()
	}
}

// Make*RepositoryFixtureWithDB
// --------------------------------

func MakeUserRepositoryFixtureWithDB(db *sqlx.DB) (*repos.UserRepository, func()) {
	userRepository := repos.NewUserRepository(db)
	return userRepository, func() {}
}

func MakeRideRepositoryFixtureWithDB(db *sqlx.DB) (*repos.RideRepository, func()) {
	rideRepository := repos.NewRideRepository(db)
	return rideRepository, func() {}
}

func MakePictureRepositoryFixtureWithDB(db *sqlx.DB) (*repos.PictureRepository, func()) {
	pictureRepository := repos.NewPictureRepository(db)
	return pictureRepository, func() {}
}

func MakeReviewRepositoryFixtureWithDB(db *sqlx.DB) (*repos.ReviewRepository, func()) {
	reviewRepository := repos.NewReviewRepository(db)
	return reviewRepository, func() {}
}

func MakeMaintenanceRepositoryFixtureWithDB(db *sqlx.DB) (*repos.MaintenanceRepository, func()) {
	maintenanceRepository := repos.NewMaintenanceRepository(db)
	return maintenanceRepository, func() {}
}
