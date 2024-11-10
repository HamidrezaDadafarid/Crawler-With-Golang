package models

import (
	"time"

	"gorm.io/gorm"
)

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
