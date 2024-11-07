package models

type Ads struct {
	Link          string // divar or sheypoor
	UniqueID      string // the unique link
	Title         string
	Description   string
	CategoryPMR   string
	City          string
	Mahale        string
	CategoryAV    string
	NumberOfViews int
	SellPrice     int
	RentPrice     int
	MortgagePrice int
	Meters        int
	NumberOfRooms int
	Age           int
	FloorNumber   int
	Pictures      []*Pictures
	Anbary        bool
	Elevator      bool
	Latitude      float64
	Longitude     float64
}

type Pictures struct {
	PictureLink string
}
