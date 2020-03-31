package postgres_test

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/internal/testutil"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// Fixtures
// --------------------------------

func truncatePictures(db *sqlx.DB) {
	db.MustExec("TRUNCATE TABLE pictures CASCADE")
	db.MustExec("TRUNCATE TABLE picture_collections CASCADE")
}

// Tests
// --------------------------------

func TestStoreSucceeds(t *testing.T) {
	pictureRepository, db, teardown := testutil.MakePictureRepositoryFixture()
	defer teardown()

	truncatePictures(db)

	picture := &models.Picture{}
	picture.Format = models.PictureFormatPNG
	picture.Data = []byte{0, 1, 2, 3}

	// first picture

	picture.ID = "picture-0"
	err := pictureRepository.Store("coll", picture)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	coll, err := pictureRepository.FetchByCollectionID("coll")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Len(t, coll, 1)
	assert.Equal(t, picture, coll[0])

	// second picture

	picture.ID = "picture-1"
	err = pictureRepository.Store("coll", picture)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	coll, err = pictureRepository.FetchByCollectionID("coll")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Len(t, coll, 2)
	assert.Equal(t, picture, coll[1])

	// third picture

	picture.ID = "picture-2"
	err = pictureRepository.Store("coll", picture)
	if !assert.Nil(t, err) {
		t.FailNow()
	}

	coll, err = pictureRepository.FetchByCollectionID("coll")
	if !assert.Nil(t, err) {
		t.FailNow()
	}
	assert.Len(t, coll, 3)
	assert.Equal(t, picture, coll[2])
}
