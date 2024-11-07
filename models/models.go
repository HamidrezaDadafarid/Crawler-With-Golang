package models

import (
	"time"

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
	CategoryPMR   uint
	Age           uint
	CategoryAV    uint
	FloorNumber   int
	Anbary        bool
	Elevator      bool
	Title         string
	Pictures      []*Pictures `gorm:"many2many:Pictures"`
	Users         []*Users    `gorm:"many2many:Users_Ads"`
}

type Pictures struct {
	gorm.Model
	ID          uint `gorm:"primaryKey;autoIncrement"`
	PictureLink string
	AdId        uint
	Ad          Ads `gorm:"foreignKey:AdId"`
}

type Users_Ads struct {
	gorm.Model
	UserId     uint `gorm:"primaryKey;autoIncrement:false"`
	AdId       uint `gorm:"primaryKey;autoIncrement:false"`
	IsBookmark bool
	Ad         Ads   `gorm:"foreignKey:AdId"`
	User       Users `gorm:"foreignKey:UserId"`
}

type Users struct {
	gorm.Model
	ID               uint `gorm:"primaryKey;autoIncrement"`
	TelegramId       string
	Role             string
	MaxSearchedItems uint
	TimeLimit        uint
	Ads              []*Ads       `gorm:"many2many:Users_Ads"`
	WatchLists       []*WatchList `gorm:"many2many:WatchList"`
}

type Filters struct {
	gorm.Model
	ID                 uint `gorm:"primaryKey;autoIncrement"`
	NumberOfRequests   uint
	StartPrice         uint
	EndPrice           uint
	City               string
	Neighborhood       string
	SartArea           uint
	EndArea            uint
	StartNumberOfRooms uint
	EndNumberOfRooms   uint
	CategoryPMR        uint
	StartAge           uint
	EndAge             uint
	CategoryAV         uint
	StartFloorNumber   int
	EndFloorNumber     int
	Anbary             bool
	Elevator           bool
	StartDate          time.Time
	EndDate            time.Time
	User               []*Users `gorm:"many2many:WatchList"`
}

type WatchList struct {
	gorm.Model
	UserID   uint `gorm:"primaryKey;autoIncrement:false"`
	FilterId uint `gorm:"primaryKey;autoIncrement:false"`
	Time     time.Time
	Filter   Filters `gorm:"foreignKey:FilterId"`
	User     Users   `gorm:"foreignKey:UserID"`
}
