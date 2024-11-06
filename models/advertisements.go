package models

type Advertisement struct {
	Source         string
	UniqueID       string
	Title          string
	Desc           string
	TypeOfAd       string
	Price          string
	City           string
	Neighbourhood  string
	Surface        int
	RoomsCount     int
	YearOfBuild    int
	TypeOfProperty string
	Warehouse      bool
	Elevator       bool
	Latitude       float64
	Longitude      float64
	// CreatedAt      time.
}
