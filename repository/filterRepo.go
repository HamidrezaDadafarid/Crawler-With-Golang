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
	if result.Error != nil {
		return filter, result.Error
	}

	return filter, nil
}

func (g *gormFilter) Delete(id uint) error {
	result := g.Db.Delete(&models.Filters{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormFilter) Get(ids []uint) ([]models.Filters, error) {
	var filters []models.Filters
	result := g.Db.Find(&filters, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return filters, nil
}

func (g *gormFilter) Update(filter models.Filters) error {
	result := g.Db.Save(&filter)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

var _ Filter = (*gormFilter)(nil)
