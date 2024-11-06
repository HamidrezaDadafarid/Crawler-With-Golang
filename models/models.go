package models

import (
	"time"

	"gorm.io/gorm"
)

type Ads struct {
	gorm.Model
	ID            uint
	Link          string `gorm:"unique;not null"`
	OwnerNumber   string
	Description   string
	NumberOfViews uint
	Price         uint
	City          string
	Mahale        string
	Meters        uint
	NumberOfRooms uint
	Category1     string
	Category2     string
	Age           uint
	FloorNumber   int
	Anbary        bool
	Elevator      bool
	AdDate        time.Time
	Pictures      []*Pictures `gorm:"many2many:Pictures"`
	Users         []*Users    `gorm:"many2many:Users_Ads"`
}

type Pictures struct {
	gorm.Model
	ID          uint
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
	ID               uint
	TelegramId       string
	Role             string
	MaxSearchedItems uint
	TimeLimit        uint
	Ads              []*Ads       `gorm:"many2many:Users_Ads"`
	WatchLists       []*WatchList `gorm:"many2many:WatchList"`
}

type Filters struct {
	gorm.Model
	ID                 uint
	NumberOfRequests   uint
	StartPrice         uint
	EndPrice           uint
	City               string
	Mahale             string
	SartArea           uint
	EndArea            uint
	StartNumberOfRooms uint
	EndNumberOfRooms   uint
	Category1          string
	Category2          string
	StartAge           uint
	EndAge             uint
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
