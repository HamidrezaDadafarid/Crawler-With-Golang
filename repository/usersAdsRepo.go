package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type User_Ad interface {
	Add(user models.Users_Ads) error
	Delete(id uint) error
	GetByUserId(userIds []uint) ([]models.Users_Ads, error)
	GetByAdId(adIds []uint) ([]models.Users_Ads, error)
	Update(user models.Users_Ads) error
}

type gormUser_Ad struct {
	Db *gorm.DB
}

func (g *gormUser_Ad) Add(userAds models.Users_Ads) error {
	result := g.Db.Create(&userAds)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormUser_Ad) Delete(id uint) error {
	result := g.Db.Delete(&models.Users_Ads{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (g *gormUser_Ad) GetByUserId(userIds []uint) ([]models.Users_Ads, error) {
	var usersAds []models.Users_Ads
	result := g.Db.Where("UserId IN ?", userIds).Find(&usersAds)
	if result.Error != nil {
		return nil, result.Error
	}

	return usersAds, nil
}

func (g *gormUser_Ad) GetByAdId(adIds []uint) ([]models.Users_Ads, error) {
	var usersAds []models.Users_Ads
	result := g.Db.Where("AdId IN ?", adIds).Find(&usersAds)
	if result.Error != nil {
		return nil, result.Error
	}

	return usersAds, nil
}

func (g *gormUser_Ad) Update(user models.Users_Ads) error {
	result := g.Db.Save(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

var _ User_Ad = (*gormUser_Ad)(nil)
