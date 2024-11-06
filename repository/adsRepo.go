package repository

import "main/models"

type Ad interface {
	Add(ad models.Ads) error
	Delete(id uint) error
	Get(filter models.Filters) (models.Ads, error)
	GetById(ids []uint) ([]models.Ads, error)
	update(ad models.Ads) error
}
