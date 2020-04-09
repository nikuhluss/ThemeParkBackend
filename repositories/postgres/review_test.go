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

func setupTestReviews(db *sqlx.DB) []string {
	reviewIDs := make([]string, 0, 3)

	tx := db.MustBegin()
	tx.MustExec("TRUNCATE TABLE reviews CASCADE")
	reviewIDs = append(reviewIDs, generator.MustInsertReview(tx))
	reviewIDs = append(reviewIDs, generator.MustInsertReview(tx))
	reviewIDs = append(reviewIDs, generator.MustInsertReview(tx))
	err := tx.Commit()
	if err != nil {
		panic(err)
	}

	return reviewIDs
}

// Tests
// --------------------------------

func TestReviewGetByIDSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	tests := setupTestReviews(db)

	for _, reviewID := range tests {
		review, err := reviewRepository.GetByID(reviewID)
		if !assert.Nil(t, err) {
			t.FailNow()
		}

		assert.Equal(t, reviewID, review.ID)
		assert.Equal(t, 1, review.Rating)
		assert.Equal(t, reviewID+" -- title", review.Title)
		assert.Equal(t, reviewID+" -- content", review.Content)
		assert.Equal(t, 2, review.PostedOn)
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

	setupTestReviews(db)

	reviews, err := reviewRepository.Fetch()
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Len(t, reviews, 3)
}

func TestReviewStoreSucceeds(t *testing.T) {
	reviewRepository, _, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	expectedReview := models.NewReview("review--ID", 1, "review--ID--title", "review--ID--content", time.Now())
	err := reviewRepository.Store(expectedReview)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	review, err := reviewRepository.GetByID(expectedReview.ID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.NotNil(t, review)
	assert.Equal(t, expectedReview.Rating, review.Rating)
	assert.Equal(t, expectedReview.Title, review.Title)
	assert.Equal(t, expectedReview.Content, review.Content)
	assert.Equal(t, expectedReview.PostedOn, review.PostedOn)
}

func TestReviewUpdateSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	tests := setupTestReviews(db)
	reviewID := tests[0]

	review, err := reviewRepository.GetByID(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	expectedReview := models.NewReview(review.ID, 2, "new title", "new content", time.Now())
	err = reviewRepository.Update(expectedReview)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	updatedReview, err := reviewRepository.GetByID(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	assert.Equal(t, expectedReview.ID, updatedReview.ID)
	assert.Equal(t, expectedReview.Rating, updatedReview.Rating)
	assert.Equal(t, expectedReview.Title, updatedReview.Title)
	assert.Equal(t, expectedReview.Content, updatedReview.Content)
	assert.Equal(t, expectedReview.PostedOn, updatedReview.PostedOn)
}

func TestReviewDeleteSucceeds(t *testing.T) {
	reviewRepository, db, teardown := testutil.MakeReviewRepositoryFixture()
	defer teardown()

	tests := setupTestReviews(db)
	reviewID := tests[0]

	err := reviewRepository.Delete(reviewID)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	review, err := reviewRepository.GetByID(reviewID)
	assert.Nil(t, review)
	assert.NotNil(t, err)
}





