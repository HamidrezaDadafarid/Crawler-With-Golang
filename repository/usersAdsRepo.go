package repository

import "main/models"

type User_Ad interface {
	Add(user models.Users_Ads) error
	Delete(id uint) error
	GetByUserId(userIds []uint) ([]models.Users_Ads, error)
	GetByAdId(adIds []uint) ([]models.Users_Ads, error)
	update(user models.Users_Ads) error
}
