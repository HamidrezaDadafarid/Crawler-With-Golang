package models

import (
	"gorm.io/gorm"
)

type Ads struct {
	gorm.Model             // 0 , 1,2,3
	ID            uint     `gorm:"primaryKey;autoIncrement"`
	Link          string   // 4
	UniqueId      string   `gorm:"unique;not null"` //5
	Longitude     float64  //6
	Latitude      float64  //7
	Description   string   //8
	NumberOfViews uint     //9
	SellPrice     uint     //10
	RentPrice     uint     //11
	MortgagePrice uint     //12
	City          string   //13
	Neighborhood  string   //14
	Meters        uint     //15
	NumberOfRooms uint     //16
	CategoryPR    uint     //17
	Age           uint     //18
	CategoryAV    uint     //19
	FloorNumber   uint     //20
	Storage       bool     //21
	Elevator      bool     //22
	Parking       bool     //23
	Title         string   //24
	PictureLink   string   //25
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
