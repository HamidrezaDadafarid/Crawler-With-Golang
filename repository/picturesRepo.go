package repository

import "main/models"

type Picture interface {
	Add(picture models.Pictures) error
	Delete(id uint) error
	Get(ids []uint) ([]models.Pictures, error)
	update(picture models.Pictures) error
}
