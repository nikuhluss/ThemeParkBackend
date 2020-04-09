package models

import (
	"time"
)

// Review is a struct that represents a review for a ride, written by an user.
type Review struct {
	ID       string    `json:"id"`
	RideID   string    `db:"ride_id" json:"rideId"`
	UserID   string    `db:"customer_id" json:"userId"`
	Rating   int       `json:"rating"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	PostedOn time.Time `json:"postedOn"`
}

// NewReview creates a new Review instance
func NewReview(ID, rating int, title, content string, postedOn time.Time) {
	return &Review{
		ID: ID,
		RideID: "",
		UserID: "",
		Rating: rating,
		Title: title,
		Content: content,
		PostedOn: postedOn,
	}
}
