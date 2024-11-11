package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type User_Ad interface {
	Add(user models.Users_Ads) error
	Delete(userId uint, adId uint) error
	GetByUserId(userIds []uint) ([]models.Users_Ads, error)
	GetByAdId(adIds []uint) ([]models.Users_Ads, error)
	Update(user models.Users_Ads) error
}

type gormUser_Ad struct {
	Db *gorm.DB
}

func NewGormUser_Ad(Db *gorm.DB) User_Ad {
	return &gormUser_Ad{
		Db: Db,
	}
}

func (g *gormUser_Ad) Add(userAds models.Users_Ads) error {
	result := g.Db.Create(&userAds)
	return result.Error
}

func (g *gormUser_Ad) Delete(userId uint, adId uint) error {
	result := g.Db.Delete(&models.Users_Ads{
		UserId: userId,
		AdId:   adId,
	})
	return result.Error
}

func (g *gormUser_Ad) GetByUserId(userIds []uint) ([]models.Users_Ads, error) {
	var usersAds []models.Users_Ads
	result := g.Db.Where("user_id IN ?", userIds).Find(&usersAds)
	return usersAds, result.Error
}

func (g *gormUser_Ad) GetByAdId(adIds []uint) ([]models.Users_Ads, error) {
	var usersAds []models.Users_Ads
	result := g.Db.Where("ad_id IN ?", adIds).Find(&usersAds)
	return usersAds, result.Error
}

func (g *gormUser_Ad) Update(user models.Users_Ads) error {
	result := g.Db.Save(&user)
	return result.Error
}

var _ User_Ad = (*gormUser_Ad)(nil)
