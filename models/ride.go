package models

// Ride is a struct that represents a ride in the park.
type Ride struct {
	ID          string
	Name        string
	Description string
	MinAge      int `db:"min_age"`
	MinHeight   int `db:"min_height"`
	Longitude   float32
	Latitude    float32
	Pictures    []Picture
	Reviews     []Review
}
