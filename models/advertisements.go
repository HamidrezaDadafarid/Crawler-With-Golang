package models

type Advertisement struct {
	Source         string
	UniqueID       string
	Title          string
	Desc           string
	TypeOfAd       string
	City           string  
	Neighbourhood  string
	TypeOfProperty string
	Price          int
	Surface        int
	RoomsCount     int
	YearOfBuild    int
	Warehouse      bool
	Elevator       bool
	Latitude       float64
	Longitude      float64
	// CreatedAt      time.
}
