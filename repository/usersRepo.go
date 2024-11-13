package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type User interface {
	Add(user models.Users) (models.Users, error)
	Delete(id uint) error
	Get(ids []uint) ([]models.Users, error)
	GetByTelegramId(telegramId string) (models.Users, error)
	GetByUsername(username string) (models.Users, error)
	Update(user models.Users) error
}

type gormUser struct {
	Db *gorm.DB
}

func NewGormUser(Db *gorm.DB) User {
	return &gormUser{
		Db: Db,
	}
}

func (g *gormUser) Add(user models.Users) (models.Users, error) {
	result := g.Db.Create(&user)
	return user, result.Error
}

func (g *gormUser) Delete(id uint) error {
	result := g.Db.Delete(&models.Users{}, id)
	return result.Error
}

func (g *gormUser) Get(ids []uint) ([]models.Users, error) {
	var users []models.Users
	result := g.Db.Find(&users, ids)
	return users, result.Error
}

func (g *gormUser) GetByTelegramId(telegramId string) (models.Users, error) {
	var user models.Users
	result := g.Db.First(&user, "telegram_id = ?", telegramId)
	return user, result.Error
}

func (g *gormUser) GetByUsername(username string) (models.Users, error) {
	var user models.Users
	result := g.Db.First(&user, "username = ?", username)
	return user, result.Error
}

func (g *gormUser) Update(user models.Users) error {
	result := g.Db.Save(&user)
	return result.Error
}

var _ User = (*gormUser)(nil)
