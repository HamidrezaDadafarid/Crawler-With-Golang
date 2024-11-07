package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Filter struct {
	Id                 uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	NumberOfRequests   uint      `json:"number_of_requests"`
	StartPrice         uint      `json:"start_price"`
	EndPrice           uint      `json:"end_price"`
	City               string    `json:"city"`
	Mahale             string    `json:"mahale"`
	StartArea          uint      `json:"start_area"`
	EndArea            uint      `json:"end_area"`
	StartNumberOfRooms uint      `json:"start_number_of_rooms"`
	EndNumberOfRooms   uint      `json:"end_number_of_rooms"`
	Category1          string    `json:"category1"`
	StartAge           uint      `json:"start_age"`
	EndAge             uint      `json:"end_age"`
	Category2          string    `json:"category2"`
	StartFloorNumber   int       `json:"start_floor_number"`
	EndFloorNumber     int       `json:"end_floor_number"`
	Anbary             bool      `json:"anbary"`
	Elevator           bool      `json:"elevator"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
}

type User struct {
	Id               uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	TelegramID       string    `json:"telegram_id" validate:"required"`
	Role             string    `json:"role"`
	MaxSearchedItems uint      `json:"max_searched_items"`
	TimeLimit        uint      `json:"time_limit"`
}

func (f *Filter) Normalize() {
	f.City = normalizeText(f.City)
	f.Mahale = normalizeText(f.Mahale)
	f.Category1 = normalizeText(f.Category1)
	f.Category2 = normalizeText(f.Category2)
}

func (u *User) Normalize() {
	u.Role = normalizeText(u.Role)
	u.TelegramID = normalizeTelegramID(u.TelegramID)
}

func (f *Filter) Validate() error {
	if f.StartPrice > f.EndPrice {
		return errors.New("start price cannot be greater than end price")
	}
	if f.StartArea > f.EndArea {
		return errors.New("start area cannot be greater than end area")
	}
	if f.StartNumberOfRooms > f.EndNumberOfRooms {
		return errors.New("start number of rooms cannot be greater than end number of rooms")
	}
	return nil
}

func (u *User) Validate() error {
	if u.TelegramID == "" {
		return errors.New("telegram ID is required")
	}
	return nil
}

func CreateFilter(db *gorm.DB, filter *Filter) error {
	filter.Normalize()
	if err := filter.Validate(); err != nil {
		return err
	}
	return db.Create(filter).Error
}

func GetFilter(db *gorm.DB, id uuid.UUID) (*Filter, error) {
	var filter Filter
	if err := db.First(&filter, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &filter, nil
}

func UpdateFilter(db *gorm.DB, filter *Filter) error {
	filter.Normalize()
	if err := filter.Validate(); err != nil {
		return err
	}
	return db.Save(filter).Error
}

func DeleteFilter(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&Filter{}, id).Error
}

func CreateUser(db *gorm.DB, user *User) error {
	user.Normalize()
	if err := user.Validate(); err != nil {
		return err
	}
	return db.Create(user).Error
}

func GetUser(db *gorm.DB, id uuid.UUID) (*User, error) {
	var user User
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(db *gorm.DB, user *User) error {
	user.Normalize()
	if err := user.Validate(); err != nil {
		return err
	}
	return db.Save(user).Error
}

func DeleteUser(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(&User{}, id).Error
}

func normalizeText(text string) string {
	return strings.TrimSpace(strings.ToLower(text))
}

func normalizeTelegramID(id string) string {
	return strings.TrimSpace(id)
}
