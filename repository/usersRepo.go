package repository

import "main/models"

type User interface {
	Add(user models.Users) error
	Delete(id uint) error
	Get(ids []uint) ([]models.Users, error)
	update(user models.Users) error
}
