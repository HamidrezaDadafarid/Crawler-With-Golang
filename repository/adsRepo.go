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
	result := g.Db.Where("SellPrice BETWEEN ? AND ? OR RentPrice BETWEEN ? AND ? OR MortagePrice BETWEEN ? AND ? OR City LIKE ? OR Neighborhood LIKE ? OR NumberOfRooms BETWEEN ? AND ? OR CategoryPR = ? OR CategoryAV = ? OR Age BETWEEN ? AND ? OR FloorNumber BETWEEN ? AND ? OR Storage = ? OR ELevator = ?", filter.StartPurchasePrice, filter.EndPurchasePrice, filter.StartRentPrice, filter.EndRentPrice, filter.StartMortgagePrice, filter.EndMortgagePrice, filter.City, filter.Neighborhood, filter.StartNumberOfRooms, filter.EndNumberOfRooms, filter.CategoryAV, filter.CategoryPR, filter.StartAge, filter.EndAge, filter.StartFloorNumber, filter.EndFloorNumber, filter.Storage, filter.Elevator).Find(&ads)
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
