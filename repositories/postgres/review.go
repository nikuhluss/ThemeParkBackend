package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// selectReviews is a query template we can reuse later
var selectReviews = psql.Select("reviews.*").From("reviews")

// ReviewRepository implements the ReviewRepository interface for postgres
type ReviewRepository struct {
	db *sqlx.DB
}

// NewReviewRepository returns a new ReviewRepository
func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db}
}

// GetByID fetches a review from the database using the given ID
func (rr *ReviewRepository) GetByID(ID string) (*models.Review, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectReviews.Where(sq.Eq{"reviews.ID": ID}).MustSql()

	review := models.Review{}
	err := udb.Get(&review, query, ID)
	if err != nil {
		return nil, err
	}

	return &review, nil

}

// Fetch fetches all reviews from the database
func (rr *ReviewRepository) Fetch() ([]*models.Review, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectReviews.MustSql()

	reviews := []*models.Review{}
	err := udb.Select(&reviews, query)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// FetchForRideSortedByRating fetches all reviews from the database for the given ride.
func (rr *ReviewRepository) FetchForRideSortedByRating(rideID string) ([]*models.Review, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectReviews.Where("ride_ID = ?").OrderBy("rating DESC").MustSql()

	reviews := []*models.Review{}
	err := udb.Select(&reviews, query, rideID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// FetchForRideSortedByDate fetches all reviews from the database for the given ride.
func (rr *ReviewRepository) FetchForRideSortedByDate(rideID string) ([]*models.Review, error) {
	db := rr.db
	udb := db.Unsafe()

	query, _ := selectReviews.Where("ride_ID = ?").OrderBy("posted_on DESC").MustSql()

	reviews := []*models.Review{}
	err := udb.Select(&reviews, query, rideID)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

// Store creates an entry for the given review model in the database
func (rr *ReviewRepository) Store(review *models.Review) error {
	db := rr.db

	insertReview, _, _ := psql.
		Insert("reviews").
		Columns("ID", "ride_ID", "customer_ID", "rating", "title", "content", "posted_on").
		Values("?", "?", "?", "?", "?", "?", "?").
		ToSql()

	_, err := db.Exec(insertReview, review.ID, review.RideID, review.UserID, review.Rating, review.Title, review.Content, review.PostedOn)
	if err != nil {
		return fmt.Errorf("insertReview: %s", err)
	}

	return nil
}

// Update updates an existing entry in the database for the given review model
func (rr *ReviewRepository) Update(review *models.Review) error {
	db := rr.db

	updateReview, _, _ := psql.
		Update("reviews").
		Set("ride_ID", "?").
		Set("customer_ID", "?").
		Set("rating", "?").
		Set("title", "?").
		Set("content", "?").
		Where("id = ?").
		ToSql()

	_, err := db.Exec(updateReview, review.RideID, review.UserID, review.Rating, review.Title, review.Content)
	if err != nil {
		return fmt.Errorf("updateReview: %s", err)
	}

	return nil
}

// Delete deletes an existing entry in the database for the given review ID
func (rr *ReviewRepository) Delete(ID string) error {
	db := rr.db

	deleteReview, _, _ := psql.Delete("reviews").Where("id = ?").ToSql()

	_, err := db.Exec(deleteReview, ID)
	if err != nil {
		return fmt.Errorf("deleteReview: %s", err)
	}

	return nil
}
