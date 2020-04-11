package postgres_test

import (
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/generator"

	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func setupTestReviews(db *sqlx.DB) ([]string, []string, []string) {

	customers := make([]string, 0)
	rides := make([]string, 0)
	reviews := make([]string, 0)

	tx := db.MustBegin()

	tx.MustExec("TRUNCATE TABLE users CASCADE")
	tx.MustExec("TRUNCATE TABLE rides CASCADE")
	tx.MustExec("TRUNCATE TABLE reviews CASCADE")

	customers = append(customers, generator.MustInsertCustomer(tx, "customer0@email.com", "customer0"))
	customers = append(customers, generator.MustInsertCustomer(tx, "customer1@email.com", "customer1"))
	customers = append(customers, generator.MustInsertCustomer(tx, "customer2@email.com", "customer2"))

	rides = append(rides, generator.MustInsertRide(tx))
	rides = append(rides, generator.MustInsertRide(tx))
	rides = append(rides, generator.MustInsertRide(tx))

	// customers[0] posted no reviews, customers[1] posted one review, etc
	// rides[0] has no reviews, rides[1] has one review, etc
	reviews = append(reviews, generator.MustInsertReview(tx, rides[1], customers[1], time.Now()))
	reviews = append(reviews, generator.MustInsertReview(tx, rides[2], customers[1], time.Now()))
	reviews = append(reviews, generator.MustInsertReview(tx, rides[2], customers[2], time.Now()))

	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return customers, rides, reviews
}

// Tests
// --------------------------------

func TestReviewGetByIDSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	_, _, reviewIDs := setupTestReviews(db)

	for _, reviewID := range reviewIDs {
		review, err := reviewRepository.GetByID(reviewID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, reviewID, review.ID)
		assert.NotEmpty(t, review.RideID)
		assert.NotEmpty(t, review.UserID)
		assert.NotEmpty(t, review.Rating)
		assert.NotEmpty(t, review.Title)
		assert.NotEmpty(t, review.Content)
	}
}

func TestReviewGetByIDNoMatchFails(t *testing.T) {
	reviewRepository, _, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	review, err := reviewRepository.GetByID("some-unknown-ID")
	assert.Nil(t, review)
	assert.NotNil(t, err)
}

func TestReviewFetchSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	_, _, reviewIDs := setupTestReviews(db)

	reviews, err := reviewRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, reviews, len(reviewIDs))
}

func TestReviewStoreSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	userIDs, rideIDs, _ := setupTestReviews(db)
	userID := userIDs[0]
	rideID := rideIDs[0]

	expectedReview := models.NewReview("review--ID", 1, "review--ID--title", "review--ID--content", time.Now().UTC())
	expectedReview.RideID = rideID
	expectedReview.UserID = userID
	err := reviewRepository.Store(expectedReview)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	review, err := reviewRepository.GetByID(expectedReview.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, review)
	assert.Equal(t, expectedReview.ID, review.ID)
	assert.Equal(t, expectedReview.RideID, review.RideID)
	assert.Equal(t, expectedReview.UserID, review.UserID)
	assert.Equal(t, expectedReview.Rating, review.Rating)
	assert.Equal(t, expectedReview.Title, review.Title)
	assert.Equal(t, expectedReview.Content, review.Content)
	assert.Equal(t, expectedReview.PostedOn, review.PostedOn)
}

func TestReviewUpdateSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	_, _, reviewIDs := setupTestReviews(db)
	reviewID := reviewIDs[0]

	review, err := reviewRepository.GetByID(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedReview := models.NewReview(review.ID, 1, "new title", "new content", time.Now().UTC())
	expectedReview.RideID = review.RideID
	expectedReview.UserID = review.UserID
	err = reviewRepository.Update(expectedReview)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedReview, err := reviewRepository.GetByID(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedReview.ID, updatedReview.ID)
	assert.Equal(t, expectedReview.RideID, updatedReview.RideID)
	assert.Equal(t, expectedReview.UserID, updatedReview.UserID)
	assert.Equal(t, expectedReview.Rating, updatedReview.Rating)
	assert.Equal(t, expectedReview.Title, updatedReview.Title)
	assert.Equal(t, expectedReview.Content, updatedReview.Content)
	assert.Equal(t, expectedReview.PostedOn, updatedReview.PostedOn)
}

func TestReviewDeleteSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	_, _, reviewIds := setupTestReviews(db)
	reviewID := reviewIds[0]

	err := reviewRepository.Delete(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	review, err := reviewRepository.GetByID(reviewID)
	assert.Nil(t, review)
	assert.NotNil(t, err)
}
