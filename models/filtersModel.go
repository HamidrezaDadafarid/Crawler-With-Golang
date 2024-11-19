package models

import (
	"time"

	"gorm.io/gorm"
)

type Filters struct {
	gorm.Model

	NumberOfRequests   uint
	StartPurchasePrice *uint
	EndPurchasePrice   *uint
	StartRentPrice     *uint
	EndRentPrice       *uint
	StartMortgagePrice *uint
	EndMortgagePrice   *uint
	City               *string
	Neighborhood       *string
	StartArea          *uint
	EndArea            *uint
	StartNumberOfRooms *uint
	EndNumberOfRooms   *uint
	CategoryPR         *uint // 0 for purchase, 1 for rents
	StartAge           *uint
	EndAge             *uint
	CategoryAV         *uint // 0 for villa, 1 for apartment
	StartFloorNumber   *uint
	EndFloorNumber     *uint
	Storage            *bool
	Elevator           *bool
	Parking            *bool
	StartDate          *time.Time
	EndDate            *time.Time
	User               []*Users `gorm:"many2many:WatchList"`
}
