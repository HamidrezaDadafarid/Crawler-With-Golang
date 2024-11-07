package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type Ad interface {
	Add(ad models.Ads) error
	Delete(id uint) error
	Get(filter models.Filters) (models.Ads, error)
	GetById(ids []uint) ([]models.Ads, error)
	update(ad models.Ads) error
}

type GormAd struct {
	Db *gorm.DB
}

func (g *GormAd) Add(ad models.Ads) error {
	result := g.Db.Create(&ad)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *GormAd) Delete(id uint) error {

}
