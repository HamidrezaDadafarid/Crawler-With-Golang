package models

import (
    "errors"
    "gorm.io/gorm"
)

type Filter struct {
    FilterID           int    `json:"filter_id" gorm:"primaryKey"`
    NumberOfRequests   uint   `json:"number_of_requests"`
    StartPrice         uint   `json:"start_price"`
    EndPrice           uint   `json:"end_price"`
    City               string `json:"city"`
    Neighborhood       string `json:"neighborhood"`
    StartArea          uint   `json:"start_area"`
    EndArea            uint   `json:"end_area"`
    StartNumberOfRooms uint   `json:"start_number_of_rooms"`
    EndNumberOfRooms   uint   `json:"end_number_of_rooms"`
    CategoryPMR        uint   `json:"category_pmr"`
    StartAge           uint   `json:"start_age"`
    EndAge             uint   `json:"end_age"`
    CategoryAV         uint   `json:"category_av"`
    StartFloorNumber   int    `json:"start_floor_number"`
    EndFloorNumber     int    `json:"end_floor_number"`
    Anbary             bool   `json:"anbary"`
    Elevator           bool   `json:"elevator"`
    StartDate          string `json:"start_date"`
    EndDate            string `json:"end_date"`
    Category2          string `json:"category2"`
}

func (f *Filter) Create(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Create(f).Error
}

func (f *Filter) Get(db *gorm.DB, id int) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.First(f, id).Error
}

func (f *Filter) Update(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Save(f).Error
}

func (f *Filter) Delete(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Delete(f).Error
}

func ListFilters(db *gorm.DB) ([]Filter, error) {
    var filters []Filter
    if db == nil {
        return filters, errors.New("database connection is nil")
    }
    err := db.Find(&filters).Error
    return filters, err
}

func (f *Filter) Normalize() {
    f.City = normalizeText(f.City)
    f.Neighborhood = normalizeText(f.Neighborhood)
}

func normalizeText(text string) string {
    // Implement normalization logic here if needed
    return text
}
