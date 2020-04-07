package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

var selectRides = psql.Select("rides.*").From("rides").OrderBy("rides.name ASC")

// RideRepository implements the RideRepository interface for postgres.
type RideRepository struct {
	db *sqlx.DB
}

// NewRideRepository creates a new RideRepository instance using the given
// database instance.
func NewRideRepository(db *sqlx.DB) *RideRepository {
	return &RideRepository{db}
}

// GetByID fetches a ride from the database using the given ID.
func (rr *RideRepository) GetByID(ID string) (*models.Ride, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectRides.Where(sq.Eq{"rides.ID": ID}).MustSql()

	ride := models.Ride{}
	err := udb.Get(&ride, query, ID)
	if err != nil {
		return nil, err
	}

	return &ride, nil
}

// Fetch fetches all rides from the database.
func (rr *RideRepository) Fetch() ([]*models.Ride, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectRides.MustSql()

	rides := []*models.Ride{}
	err := udb.Select(&rides, query)
	if err != nil {
		return nil, err
	}

	return rides, err
}

// Store creates an entry for the given ride model in the database.
func (rr *RideRepository) Store(ride *models.Ride) error {
	db := rr.db

	insertRide, _, _ := psql.
		Insert("rides").
		Columns("ID", "name", "description", "min_age", "min_height", "longitude", "latitude").
		Values("?", "?", "?", "?", "?", "?", "?").
		ToSql()

	_, err := db.Exec(insertRide, ride.ID, ride.Name, ride.Description, ride.MinAge, ride.MinHeight, ride.Longitude, ride.Latitude)
	if err != nil {
		return fmt.Errorf("inserRide: %s", err)
	}

	return nil
}

// Update updates an existing entry in the database for the given ride model.
func (rr *RideRepository) Update(ride *models.Ride) error {
	db := rr.db

	updateRide, _, _ := psql.
		Update("rides").
		Set("name", "?").
		Set("description", "?").
		Set("min_age", "?").
		Set("min_height", "?").
		Set("longitude", "?").
		Set("latitude", "?").
		Where("id = ?").
		ToSql()

	_, err := db.Exec(updateRide, ride.Name, ride.Description, ride.MinAge, ride.MinHeight, ride.Longitude, ride.Latitude, ride.ID)
	if err != nil {
		return fmt.Errorf("updateRide: %s", err)
	}

	return nil
}

// Delete deletes an existing entry in the database for the given ride ID.
func (rr *RideRepository) Delete(ID string) error {
	db := rr.db

	deleteRide, _, _ := psql.Delete("rides").Where("id = ?").ToSql()

	_, err := db.Exec(deleteRide, ID)
	if err != nil {
		return fmt.Errorf("deleteRide: %s", err)
	}

	return nil
}
