package models

import (
	"gorm.io/gorm"
)

type Filter struct {
	FilterID          int    `json:"filter_id" gorm:"primaryKey"`
	NumberOfRequests  uint   `json:"number_of_requests"`
	StartPrice        uint   `json:"start_price"`
	EndPrice          uint   `json:"end_price"`
	City              string `json:"city"`
	Neighborhood      string `json:"neighborhood"`
	StartArea         uint   `json:"start_area"`
	EndArea           uint   `json:"end_area"`
	StartNumberOfRooms uint  `json:"start_number_of_rooms"`
	EndNumberOfRooms  uint   `json:"end_number_of_rooms"`
	CategoryPMR       uint   `json:"category_pmr"`
	StartAge          uint   `json:"start_age"`
	EndAge            uint   `json:"end_age"`
	CategoryAV        uint   `json:"category_av"`
	StartFloorNumber  int    `json:"start_floor_number"`
	EndFloorNumber    int    `json:"end_floor_number"`
	Anbary            bool   `json:"anbary"`
	Elevator          bool   `json:"elevator"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	Category2         string `json:"category2"`
}

func (f *Filter) CreateFilter(db *gorm.DB) error {
	return db.Create(f).Error
}

func (f *Filter) GetFilter(db *gorm.DB, id int) error {
	return db.First(f, id).Error
}

func (f *Filter) UpdateFilter(db *gorm.DB) error {
	return db.Save(f).Error
}

func (f *Filter) DeleteFilter(db *gorm.DB) error {
	return db.Delete(f).Error
}

func ListFilters(db *gorm.DB) ([]Filter, error) {
	var filters []Filter
	err := db.Find(&filters).Error
	return filters, err
}

func ListFiltersByCity(db *gorm.DB, city string) ([]Filter, error) {
	var filters []Filter
	err := db.Where("city = ?", city).Find(&filters).Error
	return filters, err
}

func ListFiltersByPriceRange(db *gorm.DB, minPrice, maxPrice uint) ([]Filter, error) {
	var filters []Filter
	err := db.Where("start_price >= ? AND end_price <= ?", minPrice, maxPrice).Find(&filters).Error
	return filters, err
}

func ListFiltersByRooms(db *gorm.DB, minRooms, maxRooms uint) ([]Filter, error) {
	var filters []Filter
	err := db.Where("start_number_of_rooms >= ? AND end_number_of_rooms <= ?", minRooms, maxRooms).Find(&filters).Error
	return filters, err
}

func (f *Filter) Normalize() {
	f.City = normalizeText(f.City)
	f.Neighborhood = normalizeText(f.Neighborhood)
}

func normalizeText(text string) string {
	return text
}
