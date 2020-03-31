package models

// Ride is a struct that represents a ride in the park.
type Ride struct {
	ID             string
	Name           string
	Description    string
	MinAge         int `db:"min_age"`
	MinHeight      int `db:"min_height"`
	Longitude      int
	Latitude       int
	Pictures       []*Picture
	Reviews        []*Review
	ReviewsAverage int
}

// NewRide creates a new Ride instance.
func NewRide(ID, name, description string, minAge, minHeight, longitude, latitude int) *Ride {
	return &Ride{
		ID:             ID,
		Name:           name,
		Description:    description,
		MinAge:         minAge,
		MinHeight:      minHeight,
		Longitude:      longitude,
		Latitude:       latitude,
		Pictures:       nil,
		Reviews:        nil,
		ReviewsAverage: 0,
	}
}
