package models

import (
    "errors"
    "strings"
    "time"
    "gorm.io/gorm"
)

type Filters struct {
    gorm.Model
    ID                 uint     `gorm:"primaryKey;autoIncrement"`
    NumberOfRequests   uint
    StartPrice         uint
    EndPrice           uint
    City               string
    Neighborhood       string
    StartArea          uint
    EndArea            uint
    StartNumberOfRooms uint
    EndNumberOfRooms   uint
    CategoryPMR        uint
    StartAge           uint
    EndAge             uint
    CategoryAV         uint
    StartFloorNumber   int
    EndFloorNumber     int
    Anbary             bool
    Elevator           bool
    StartDate          time.Time
    EndDate            time.Time
}

func (f *Filters) Normalize() {
    f.City = strings.ToLower(f.City)
    f.Neighborhood = strings.ToLower(f.Neighborhood)
}

func (f *Filters) Validate() error {
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
