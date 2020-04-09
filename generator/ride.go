package generator

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertRideWithID inserts a new ride using the given ID. Note that the
// other fields of the ride use the same ID but with an added suffix. E.g.
//
// name: <ID>--name
// description: <ID>--description
// ...
//
// Non-string fields are filled with incremental values starting with 1.
func InsertRideWithID(execer Execer, ID string) (string, error) {

	nameTemplate := fmt.Sprintf("%s - ###", gofakeit.BeerName())
	name := gofakeit.Numerify(nameTemplate)
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

// InsertRide is similar to InsertRideWithID but generates the ID instead.
func InsertRide(execer Execer) (string, error) {
	ID := gofakeit.UUID()
	return InsertRideWithID(execer, ID)
}

// MustInsertRideWithID is similar to InsertRideWithID but panics on error.
func MustInsertRideWithID(mustExecer MustExecer, ID string) string {
	return MustInsert(InsertRideWithID(&AsExecer{mustExecer}, ID))
}

// MustInsertRide is similar to InsertRide but panics on error.
func MustInsertRide(mustExecer MustExecer) string {
	return MustInsert(InsertRide(&AsExecer{mustExecer}))
}
