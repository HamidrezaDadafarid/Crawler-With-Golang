package models

type Advertisement struct {
	Link         string   `json:"link,omitempty"`
	Picture      string   `json:"image,omitempty"`
	SellerNumber string   `json:"number,omitempty"`
	Description  string   `json:"desc,omitempty"`
	Filters      *Filters `json:"filters"`
}

type Filters struct {
	Price          int
	City           string
	Neighbourhood  string
	Surface        int
	Rooms          int
	TypeOfAd       string
	Age            int
	TypeOfProperty string
	Floor          int
	Warehouse      bool
	Elevator       bool
	Date           string
}
