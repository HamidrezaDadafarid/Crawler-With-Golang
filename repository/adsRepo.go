package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type Ad interface {
	Add(ad models.Ads) (models.Ads, error)
	Delete(id uint) error
	Get(filter models.Filters) ([]models.Ads, error)
	GetById(ids []uint) ([]models.Ads, error)
	Update(ad models.Ads) error
}

type gormAd struct {
	Db *gorm.DB
}

func NewGormAd(Db *gorm.DB) Ad {
	return &gormAd{
		Db: Db,
	}
}

func (g *gormAd) Add(ad models.Ads) (models.Ads, error) {
	result := g.Db.Create(&ad)
	return ad, result.Error
}

func (g *gormAd) Delete(id uint) error {
	result := g.Db.Delete(&models.Ads{}, id)
	return result.Error
}

func (g *gormAd) Get(filter models.Filters) ([]models.Ads, error) {
	var ads []models.Ads
	result := g.Db.Where("sell_price BETWEEN ? AND ? OR rent_price BETWEEN ? AND ? OR mortage_price BETWEEN ? AND ? OR city LIKE ? OR neighborhood LIKE ? OR number_of_rooms BETWEEN ? AND ? OR category_av = ? OR category_pr = ? OR age BETWEEN ? AND ? OR floor_number BETWEEN ? AND ? OR storage = ? OR elevator = ?", filter.StartPurchasePrice, filter.EndPurchasePrice, filter.StartRentPrice, filter.EndRentPrice, filter.StartMortgagePrice, filter.EndMortgagePrice, filter.City, filter.Neighborhood, filter.StartNumberOfRooms, filter.EndNumberOfRooms, filter.CategoryAV, filter.CategoryPR, filter.StartAge, filter.EndAge, filter.StartFloorNumber, filter.EndFloorNumber, filter.Storage, filter.Elevator).Find(&ads)
	return ads, result.Error
}

func (g *gormAd) GetById(ids []uint) ([]models.Ads, error) {
	var ads []models.Ads
	result := g.Db.Find(&ads, ids)
	return ads, result.Error
}

func (g *gormAd) Update(ad models.Ads) error {
	result := g.Db.Save(&ad)
	return result.Error
}

var _ Ad = (*gormAd)(nil)


func AddAd(db *gorm.DB, ad *models.Ads) error {
    return db.Create(ad).Error
}

func DeleteAd(db *gorm.DB, adID uint) error {
    return db.Delete(&models.Ads{}, adID).Error
}

func EditAd(db *gorm.DB, ad *models.Ads) error {
    return db.Save(ad).Error
}

func GetAds(db *gorm.DB, adID uint) ([]models.Ads, error) {
    var ads []models.Ads
    var result *gorm.DB
    if adID != 0 {
        result = db.Where("id = ?", adID).Find(&ads)
    } else {
        result = db.Find(&ads)
    }
    return ads, result.Error
}

