package repository

import "main/models"

type WatchList interface {
	Add(watchList models.WatchList) error
	Delete(id uint) error
	GetByUserId(userIds []uint) ([]models.WatchList, error)
	GetByFilterId(filterIds []uint) ([]models.WatchList, error)
	update(watchList models.WatchList) error
}
