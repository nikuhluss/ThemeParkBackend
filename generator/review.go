package generator

import (
	"time"

	"github.com/brianvoe/gofakeit/v4"
)

// InsertReview inserts a review for the given ride and time.
func InsertReview(execer Execer, rideID, customerID string, postedOn time.Time) (string, error) {
	ID := gofakeit.UUID()
	rating := int(gofakeit.Float32Range(2, 6))
	title := gofakeit.Sentence(4)
	content := gofakeit.Sentence(10)

	query := `
	INSERT INTO reviews (ID, ride_id, customer_id, rating, title, content, posted_on)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := execer.Exec(query, ID, rideID, customerID, rating, title, content, postedOn)
	if err != nil {
		return "", err
	}

	return ID, nil
}

// MustInsertReview is like InsertReview but panics on error.
func MustInsertReview(mustExecer MustExecer, rideID, customerID string, postedOn time.Time) string {
	return MustInsert(InsertReview(&AsExecer{mustExecer}, rideID, customerID, postedOn))
}
