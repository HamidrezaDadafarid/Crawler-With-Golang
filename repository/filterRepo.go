package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type Filter interface {
	Add(filter models.Filters) (models.Filters, error)
	Delete(id uint) error
	Get(ids []uint) ([]models.Filters, error)
	Update(filter models.Filters) error
}

type gormFilter struct {
	Db *gorm.DB
}

func NewGormFilter(Db *gorm.DB) Filter {
	return &gormFilter{
		Db: Db,
	}
}

func (g *gormFilter) Add(filter models.Filters) (models.Filters, error) {
	result := g.Db.Create(&filter)
	return filter, result.Error
}

func (g *gormFilter) Delete(id uint) error {
	result := g.Db.Delete(&models.Filters{}, id)
	return result.Error
}

func (g *gormFilter) Get(ids []uint) ([]models.Filters, error) {
	var filters []models.Filters
	result := g.Db.Find(&filters, ids)
	return filters, result.Error
}

func (g *gormFilter) Update(filter models.Filters) error {
	result := g.Db.Save(&filter)
	return result.Error
}

var _ Filter = (*gormFilter)(nil)
