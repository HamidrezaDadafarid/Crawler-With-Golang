package models

import (
	"gorm.io/gorm"
)

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
