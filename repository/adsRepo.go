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
	if result.Error != nil {
		return ad, result.Error
	}

	return ad, nil
}

func (g *gormAd) Delete(id uint) error {
	result := g.Db.Delete(&models.Ads{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormAd) Get(filter models.Filters) ([]models.Ads, error) {
	var ads []models.Ads
	result := g.Db.Where("SellPrice BETWEEN ? AND ? OR RentPrice BETWEEN ? AND ? OR MortagePrice BETWEEN ? AND ? OR City LIKE ? OR Mahale LIKE ? OR NumberOfRooms BETWEEN ? AND ? OR CategoryPMR = ? OR CategoryAV = ? OR Age BETWEEN ? AND ? OR FloorNumber BETWEEN ? AND ? OR Anbary = ? OR ELevator = ?", filter.StartPrice, filter.EndPrice, filter.StartPrice, filter.EndPrice, filter.StartPrice, filter.EndPrice, filter.City, filter.Neighborhood, filter.StartNumberOfRooms, filter.EndNumberOfRooms, filter.CategoryAV, filter.CategoryPMR, filter.StartAge, filter.EndAge, filter.StartFloorNumber, filter.EndFloorNumber, filter.Anbary, filter.Elevator).Find(&ads)
	if result.Error != nil {
		return nil, result.Error
	}

	return ads, nil
}

func (g *gormAd) GetById(ids []uint) ([]models.Ads, error) {
	var ads []models.Ads
	result := g.Db.Find(&ads, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return ads, nil
}

func (g *gormAd) Update(ad models.Ads) error {
	result := g.Db.Save(&ad)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

var _ Ad = (*gormAd)(nil)
