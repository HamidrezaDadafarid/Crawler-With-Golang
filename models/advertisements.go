package models

type Advertisement struct {
	Source         string
	UniqueID       string
	Title          string
	Desc           string
	Price          string
	City           string
	Neighbourhood  string
	Surface        string
	TypeOfAd       string
	RoomsCount     int
	Age            int
	TypeOfProperty string
	Warehouse      bool
	Elevator       bool
	Latitude       float64
	Longitude      float64
	// CreatedAt      time.
}
