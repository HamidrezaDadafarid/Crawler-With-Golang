package repository

import "main/models"

type Fiter interface {
	Add(filter models.Filters) error
	Delete(id uint) error
	Get(ids []uint) ([]models.Filters, error)
	update(filter models.Filters) error
}
