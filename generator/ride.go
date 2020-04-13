package generator

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertRideWithName inserts a ride with the given name and ID.
func InsertRideWithName(execer Execer, name string) (string, error) {
	ID := gofakeit.UUID()
	description := gofakeit.Sentence(16)
	minAge := 1
	minHeight := 2
	longitude := 3.0
	latitude := 4.0

	query := `
	INSERT INTO rides (ID, picture_collection_id, name, description, min_age, min_height, longitude, latitude)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := execer.Exec(query, ID, nil, name, description, minAge, minHeight, longitude, latitude)
	if err != nil {
		return "", err
	}

	return ID, nil

}

// InsertRide is similar to InsertRideWithName but generates the name instead.
func InsertRide(execer Execer) (string, error) {
	nameTemplate := fmt.Sprintf("%s - ###", gofakeit.BeerName())
	name := gofakeit.Numerify(nameTemplate)
	return InsertRideWithName(execer, name)
}

// MustInsertRideWithName is similar to InsertRideWithID but panics on error.
func MustInsertRideWithName(mustExecer MustExecer, name string) string {
	return MustInsert(InsertRideWithName(&AsExecer{mustExecer}, name))
}

// MustInsertRide is similar to InsertRide but panics on error.
func MustInsertRide(mustExecer MustExecer) string {
	return MustInsert(InsertRide(&AsExecer{mustExecer}))
}
