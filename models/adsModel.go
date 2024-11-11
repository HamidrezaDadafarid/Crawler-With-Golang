package models

import (
	"gorm.io/gorm"
)

type Ads struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Link          string `gorm:"unique;not null"`
	UniqueId      string
	Longitude     int
	Latitude      int
	Description   string
	NumberOfViews uint
	SellPrice     uint
	RentPrice     uint
	MortagePrice  uint
	City          string
	Mahale        string
	Meters        uint
	NumberOfRooms uint
	CategoryPR    uint
	Age           uint
	CategoryAV    uint
	FloorNumber   int
	Anbary        bool
	Elevator      bool
	Title         string
	PictureLink   string
	Users         []*Users `gorm:"many2many:Users_Ads"`
}
