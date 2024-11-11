package models

import "gorm.io/gorm"

type Users_Ads struct {
	gorm.Model
	UserId     uint `gorm:"primaryKey;autoIncrement:false"`
	AdId       uint `gorm:"primaryKey;autoIncrement:false"`
	IsBookmark bool
	Ad         Ads   `gorm:"foreignKey:AdId"`
	User       Users `gorm:"foreignKey:UserId"`
}
