package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"gitlab.com/uh-spring-2020/cosc-3380-team-14/backend/models"
)

// PictureRepository implements the PictureRepository interface for postgres.
type PictureRepository struct {
	db *sqlx.DB
}

// NewPictureRepository returns a new PictureRepository
func NewPictureRepository(db *sqlx.DB) *PictureRepository {
	return &PictureRepository{db}
}

// GetByID fetches a single picture using the given ID.
func (pr *PictureRepository) GetByID(ID string) (*models.Picture, error) {
	db := pr.db
	udb := db.Unsafe()

	query, _ := psql.Select("pictures.*").From("pictures").Where("pictures.ID = ?").MustSql()

	picture := models.Picture{}
	err := udb.Get(&picture, query, ID)
	if err != nil {
		return nil, err
	}

	return &picture, nil
}

// Delete deletes a single picture using the given ID.
func (pr *PictureRepository) Delete(ID string) error {
	db := pr.db

	deletePicture, _, _ := psql.Delete("pictures").Where("ID = ?").ToSql()

	_, err := db.Exec(deletePicture, ID)
	if err != nil {
		return fmt.Errorf("deletePicture: %s", err)
	}

	return nil
}

// FetchByCollectionID returns a collection of pictures.
func (pr *PictureRepository) FetchByCollectionID(collectionID string) ([]*models.Picture, error) {
	db := pr.db
	udb := db.Unsafe()

	query, _ := psql.
		Select("pictures.*").
		From("pictures").
		LeftJoin("pictures_in_collection ON pictures_in_collection.picture_ID = pictures.ID").
		Where("pictures_in_collection.collection_ID = ?").
		MustSql()

	pictures := []*models.Picture{}
	err := udb.Select(&pictures, query, collectionID)
	if err != nil {
		return nil, err
	}

	return pictures, nil
}

// Store stores the given picture under the given collection ID.
func (pr *PictureRepository) Store(collectionID string, picture *models.Picture) error {
	db := pr.db

	ensureCollection, _, _ := psql.Insert("picture_collections").Columns("ID").Values("?").Suffix("ON CONFLICT DO NOTHING").ToSql()
	_, err := db.Exec(ensureCollection, collectionID)
	if err != nil {
		return fmt.Errorf("ensureCollection: %s", err)
	}

	insertPicture, _, _ := psql.
		Insert("pictures").
		Columns("ID", "format", "blob").
		Values("?", "?", "?").
		ToSql()

	insertInCollection, _, _ := psql.
		Insert("pictures_in_collection").
		Columns("collection_ID", "picture_ID", "picture_sequence").
		Select(psql.Select("?, ?, COUNT(*)").From("pictures_in_collection").Where("collection_ID = ?")).
		ToSql()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	{
		_, err = tx.Exec(ensureCollection, collectionID)
		if err != nil {
			return fmt.Errorf("ensureCollection: %s", err)
		}

		_, err = tx.Exec(insertPicture, picture.ID, picture.Format, picture.Data)
		if err != nil {
			return fmt.Errorf("insertPicture: %s", err)
		}

		_, err = tx.Exec(insertInCollection, collectionID, picture.ID, collectionID)
		if err != nil {
			return fmt.Errorf("insertInColleciton: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateCollectionOrdering updates the ordering of pictures under the given collection ID.
func (pr *PictureRepository) UpdateCollectionOrdering(collectionID string, fromIndex, toIndex int) error {
	return nil
}

// DeleteCollection deletes all pictures under the given collection ID.
func (pr *PictureRepository) DeleteCollection(collectionID string) error {
	db := pr.db

	deletePicturesInCollection, _, _ := psql.
		Delete("pictures_in_collection").
		Where("collection_ID = ?").
		ToSql()

	deleteCollection, _, _ := psql.
		Delete("picture_collections").
		Where("ID = ?").
		ToSql()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	{
		_, err = tx.Exec(deletePicturesInCollection, collectionID)
		if err != nil {
			return fmt.Errorf("deletePicturesInCollection: %s", err)
		}

		_, err := tx.Exec(deleteCollection, collectionID)
		if err != nil {
			return fmt.Errorf("deleteCollection: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
