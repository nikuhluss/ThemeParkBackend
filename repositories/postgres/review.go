package postgres

import (
	"github.com/jmoiron/sqlx"
)

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{
		db,
	}
}
