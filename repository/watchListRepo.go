package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type WatchList interface {
	Add(watchList models.WatchList) error
	Delete(userId uint, filterId uint) error
	GetByUserId(userIds []uint) ([]models.WatchList, error)
	GetByFilterId(filterIds []uint) ([]models.WatchList, error)
	Update(watchList models.WatchList) error
}

type gormWatchList struct {
	Db *gorm.DB
}

func NewWatchList(Db *gorm.DB) WatchList {
	return &gormWatchList{
		Db: Db,
	}
}

func (g *gormWatchList) Add(watchList models.WatchList) error {
	result := g.Db.Create(&watchList)
	return result.Error
}

func (g *gormWatchList) Delete(userId uint, filterId uint) error {
	result := g.Db.Delete(&models.WatchList{
		UserID:   userId,
		FilterId: filterId,
	})
	return result.Error
}

func (g *gormWatchList) GetByUserId(userIds []uint) ([]models.WatchList, error) {
	var watchLists []models.WatchList
	result := g.Db.Where("user_id IN ?", userIds).Find(&watchLists)
	return watchLists, result.Error
}

func (g *gormWatchList) GetByFilterId(filterIds []uint) ([]models.WatchList, error) {
	var watchLists []models.WatchList
	result := g.Db.Where("filter_id IN ?", filterIds).Find(&watchLists)
	return watchLists, result.Error
}

func (g *gormWatchList) Update(watchList models.WatchList) error {
	result := g.Db.Save(&watchList)
	return result.Error
}

var _ WatchList = (*gormWatchList)(nil)
