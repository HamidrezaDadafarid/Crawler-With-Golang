package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type User interface {
	Add(user models.Users) (models.Users, error)
	Delete(id uint) error
	Get(ids []uint) ([]models.Users, error)
	Update(user models.Users) error
}

type gormUser struct {
	Db *gorm.DB
}

func (g *gormUser) Add(user models.Users) (models.Users, error) {
	result := g.Db.Create(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (g *gormUser) Delete(id uint) error {
	result := g.Db.Delete(&models.Users{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormUser) Get(ids []uint) ([]models.Users, error) {
	var users []models.Users
	result := g.Db.Find(&users, ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (g *gormUser) Update(user models.Users) error {
	result := g.Db.Save(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

var _ User = (*gormUser)(nil)
