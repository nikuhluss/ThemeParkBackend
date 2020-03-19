package models

// Ride is a struct that represents a ride in the park.
type Ride struct {
	ID          string
	Name        string
	Description string
	MinAge      int `db:"min_age"`
	MinHeight   int `db:"min_height"`
	Longitude   int
	Latitude    int
	Pictures    []Picture
	Reviews     []Review
}

func NewRide(ID, name, description string, minAge, minHeight, longitude, latitude int, picture []Picture, review []Review) *Ride {
	return &Ride{
		ID:           ID,
		Name:         name,
		Description:  description,
		MinAge:       minAge,
		MinHeight:    minHeight,
		Longitude:    longitude,
		Latitude:     latitude,
		Pictures:     picture,
		Reviews:       review,  
	}
}