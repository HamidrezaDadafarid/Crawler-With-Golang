package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	UserID           int    `json:"user_id" gorm:"primaryKey"`
	TelegramID       int64  `json:"telegram_id"`
	Role             string `json:"role"`
	MaxSearchedItems uint   `json:"max_searched_items"`
	TimeLimit        uint   `json:"time_limit"`
}

func (u *User) CreateUser(db *gorm.DB) error {
	return db.Create(u).Error
}

func (u *User) GetUser(db *gorm.DB, id int) error {
	return db.First(u, id).Error
}

func (u *User) UpdateUser(db *gorm.DB) error {
	return db.Save(u).Error
}

func (u *User) DeleteUser(db *gorm.DB) error {
	return db.Delete(u).Error
}

func ListUsers(db *gorm.DB) ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserByTelegramID(db *gorm.DB, telegramID int64) (User, error) {
	var user User
	err := db.Where("telegram_id = ?", telegramID).First(&user).Error
	return user, err
}

func ListUsersByRole(db *gorm.DB, role string) ([]User, error) {
	var users []User
	err := db.Where("role = ?", role).Find(&users).Error
	return users, err
}

func (u *User) UpdateTimeLimit(db *gorm.DB, newTimeLimit uint) error {
	u.TimeLimit = newTimeLimit
	return db.Save(u).Error
}

func (u *User) CheckSearchLimit() error {
	if u.MaxSearchedItems == 0 {
		return errors.New("search limit reached")
	}
	u.MaxSearchedItems--
	return nil
}
