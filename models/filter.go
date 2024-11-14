package models

import (
    "errors"
    "strings"
    "gorm.io/gorm"
)

type Filter struct {
    FilterID           uint   `json:"filter_id" gorm:"primaryKey"`
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
}

func (f *Filter) Normalize() {
    f.City = strings.ToLower(f.City)
    f.Neighborhood = strings.ToLower(f.Neighborhood)
}

func (f *Filter) Validate() error {
    if f.StartPrice > f.EndPrice {
        return errors.New("start_price cannot be greater than end_price")
    }
    if f.StartArea > f.EndArea {
        return errors.New("start_area cannot be greater than end_area")
    }
    if f.StartNumberOfRooms > f.EndNumberOfRooms {
        return errors.New("start_number_of_rooms cannot be greater than end_number_of_rooms")
    }
    if f.StartAge > f.EndAge {
        return errors.New("start_age cannot be greater than end_age")
    }
    if f.StartFloorNumber > f.EndFloorNumber {
        return errors.New("start_floor_number cannot be greater than end_floor_number")
    }
    return nil
}

func (f *Filter) Create(db *gorm.DB) error {
    f.Normalize()
    if err := f.Validate(); err != nil {
        return err
    }
    return db.Create(f).Error
}

func (f *Filter) GetByID(db *gorm.DB, id uint) error {
    return db.First(f, id).Error
}

func (f *Filter) Update(db *gorm.DB) error {
    f.Normalize()
    if err := f.Validate(); err != nil {
        return err
    }
    return db.Save(f).Error
}

func (f *Filter) Delete(db *gorm.DB) error {
    return db.Delete(f).Error
}
