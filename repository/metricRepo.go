package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type Metric interface {
	Add(user models.Metrics) (models.Metrics, error)
	Delete(id uint) error
	DeleteAll() error
	Get(ids []uint) ([]models.Metrics, error)
}

type gormMetric struct {
	Db *gorm.DB
}

func NewGormUMetric(Db *gorm.DB) Metric {
	return &gormMetric{
		Db: Db,
	}
}

func (g *gormMetric) Add(metric models.Metrics) (models.Metrics, error) {
	result := g.Db.Create(&metric)
	return metric, result.Error
}

func (g *gormMetric) Delete(id uint) error {
	result := g.Db.Delete(&models.Metrics{}, id)
	return result.Error
}

func (g *gormMetric) DeleteAll() error {
	result := g.Db.Where("TRUE").Delete(&models.Metrics{})
	return result.Error
}

func (g *gormMetric) Get(ids []uint) ([]models.Metrics, error) {
	var metrics []models.Metrics
	result := g.Db.Find(&metrics, ids)
	return metrics, result.Error
}

func (g *gormMetric) GetAll() ([]models.Metrics, error) {
	var metrics []models.Metrics
	result := g.Db.Find(&metrics)
	return metrics, result.Error
}

var _ Metric = (*gormMetric)(nil)
