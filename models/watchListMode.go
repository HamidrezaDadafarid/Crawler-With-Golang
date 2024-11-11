package models

import (
	"time"

	"gorm.io/gorm"
)

type WatchList struct {
	gorm.Model
	UserID   uint `gorm:"primaryKey;autoIncrement:false"`
	FilterId uint `gorm:"primaryKey;autoIncrement:false"`
	Time     time.Time
	Filter   Filters `gorm:"foreignKey:FilterId"`
	User     Users   `gorm:"foreignKey:UserID"`
}
