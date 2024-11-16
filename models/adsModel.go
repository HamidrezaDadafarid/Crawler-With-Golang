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
	Neighborhood  string
	Meters        uint
	NumberOfRooms uint
	CategoryPR    uint
	Age           uint
	CategoryAV    uint
	FloorNumber   uint
	Storage       bool
	Elevator      bool
	Title         string
	PictureLink   string
	Users         []*Users `gorm:"many2many:Users_Ads"`
}

func AddAd(db *gorm.DB, ad *Ads) error {
    return db.Create(ad).Error
}

func DeleteAd(db *gorm.DB, adID uint) error {
    return db.Delete(&Ads{}, adID).Error
}

func EditAd(db *gorm.DB, ad *Ads) error {
    return db.Save(ad).Error
}

func GetAds(db *gorm.DB, adID uint) ([]Ads, error) {
    var ads []Ads
    var result *gorm.DB
    if adID != 0 {
        result = db.Where("id = ?", adID).Find(&ads)
    } else {
        result = db.Find(&ads)
    }
    return ads, result.Error
}

