package models

import (
	"gorm.io/gorm"
)

type WatchList struct {
	gorm.Model
	UserID   uint `gorm:"primaryKey;autoIncrement:false"`
	FilterId uint `gorm:"primaryKey;autoIncrement:false"`
	Time     int
	Filter   Filters `gorm:"foreignKey:FilterId"`
	User     Users   `gorm:"foreignKey:UserID"`
}
