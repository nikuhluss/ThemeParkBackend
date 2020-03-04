package models

// Ride is a struct that represents a ride in the park.
type Ride struct {
	ID          string
	Name        string
	Description string
	MinAge      int
	MinHeight   int
	Longitude   float32
	Latitude    float32
	Pictures    []Picture
	Reviews     []Review
}
