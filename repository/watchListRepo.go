package repository

import (
	"main/models"

	"gorm.io/gorm"
)

type WatchList interface {
	Add(watchList models.WatchList) error
	Delete(id uint) error
	GetByUserId(userIds []uint) ([]models.WatchList, error)
	GetByFilterId(filterIds []uint) ([]models.WatchList, error)
	Update(watchList models.WatchList) error
}

type gormWatchList struct {
	Db *gorm.DB
}

func (g *gormWatchList) Add(watchList models.WatchList) error {
	result := g.Db.Create(&watchList)
	return result.Error
}

func (g *gormWatchList) Delete(id uint) error {
	result := g.Db.Delete(&models.WatchList{}, id)
	return result.Error
}

func (g *gormWatchList) GetByUserId(userIds []uint) ([]models.WatchList, error) {
	var watchLists []models.WatchList
	result := g.Db.Where("UserId IN ?", userIds).Find(&watchLists)
	return watchLists, result.Error
}

func (g *gormWatchList) GetByFilterId(filterIds []uint) ([]models.WatchList, error) {
	var watchLists []models.WatchList
	result := g.Db.Where("FilterId IN ?", filterIds).Find(&watchLists)
	return watchLists, result.Error
}

func (g *gormWatchList) Update(watchList models.WatchList) error {
	result := g.Db.Save(&watchList)
	return result.Error
}

var _ WatchList = (*gormWatchList)(nil)
