package models

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	TelegramId string `gorm:"unique"`
	Role       string

	ErrorsCount uint
	Ads         []*Ads       `gorm:"many2many:Users_Ads"`
	WatchLists  []*WatchList `gorm:"many2many:WatchList"`
}
