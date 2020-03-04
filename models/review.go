package models

// Review is a struct that represents a review for a ride, written by an user.
type Review struct {
	ID      string
	RideID  string
	UserID  string
	Rating  int
	Title   string
	Content string
}
