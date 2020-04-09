package generator

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v4"
)

func InsertReviewWithID(execer Execer, ID string) (string, error) {
	rating := 1
	title := fmt.Sprintf("%s -- title", ID)
	content := fmt.Sprintf("%s -- content", ID)
	posted_on := 2

	query := `
	INSERT INTO reviews (ID, ride_id, customer_id, rating, title, content, posted_on)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := execer.Exec(query, ID, nil, rating, title, content, posted_on)

	if err != nil {
		return "", err
	}

	return ID, nil
}

func InsertReview(execer Execer) (string, error) {
	IDTemplate := fmt.Sprintf("%s - ####", gofakeit.BeerName())
	ID := gofakeit.Numerify(IDTemplate)
	return InsertRideWithID(execer, ID)
}

func MustInsertReviewWithID(mustExecer MustExecer, ID string) string {
	return MustInsert(InsertReviewWithID(&AsExecer{mustExecer}, ID))
}

func MustInsertReview(mustExecer MustExecer) string {
	return MustInsert(InsertReview(&AsExecer{mustExecer}))
}

/* func InsertReview(execer Execer, rideID, customerID string) (string, error) {
	// wordsInTitle := int(gofakeit.Float32Range(3, 5))
	// paragraphsInContent := int(gofakeit.Float32Range(1, 3))

	// rating := int(gofakeit.Float32Range(2, 5))
	// title := gofakeit.Sentence(wordsInTitle)
	// content := gofakeit.Paragraph(paragraphsInContent, 3, 10, ".")
	// postedOn := gofakeit.DateRange(year1970, year2000)

	return "", nil
} */

