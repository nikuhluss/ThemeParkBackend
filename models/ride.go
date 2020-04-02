package models

// Ride is a struct that represents a ride in the park.
type Ride struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	MinAge         int        `db:"min_age" json:"minAge"`
	MinHeight      int        `db:"min_height" json:"minHeight"`
	Longitude      float64    `json:"longitude"`
	Latitude       float64    `json:"latitude"`
	Pictures       []*Picture `json:"pictures"`
	Reviews        []*Review  `json:"reviews"`
	ReviewsAverage int        `json:"reviewsAverage"`
}

// NewRide creates a new Ride instance.
func NewRide(ID, name, description string, minAge, minHeight int, longitude, latitude float64) *Ride {
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
