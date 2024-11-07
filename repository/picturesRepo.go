package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type Picture interface {
	Add(picture models.Pictures) (models.Pictures, error)
	Delete(id uint) error
	Get(ids []uint) ([]models.Pictures, error)
	Update(picture models.Pictures) error
}

type gormPicture struct {
	Db *gorm.DB
}

func (g *gormPicture) Add(picture models.Pictures) (models.Pictures, error) {
	result := g.Db.Create(&picture)
	if result.Error != nil {
		return picture, result.Error
	}

	return picture, nil
}

func (g *gormPicture) Delete(id uint) error {
	result := g.Db.Delete(&models.Pictures{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormPicture) Get(ids []uint) ([]models.Pictures, error) {
	var pictures []models.Pictures
	result := g.Db.Find(&pictures, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return pictures, nil
}

func (g *gormPicture) Update(picture models.Pictures) error {
	result := g.Db.Save(&picture)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

var _ Picture = (*gormPicture)(nil)
