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

func (u *User) Create(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Create(u).Error
}

func (u *User) Get(db *gorm.DB, id int) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.First(u, id).Error
}

func (u *User) Update(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Save(u).Error
}

func (u *User) Delete(db *gorm.DB) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
    return db.Delete(u).Error
}

func ListUsers(db *gorm.DB) ([]User, error) {
    var users []User
    if db == nil {
        return users, errors.New("database connection is nil")
    }
    err := db.Find(&users).Error
    return users, err
}

func GetUserByTelegramID(db *gorm.DB, telegramID int64) (User, error) {
    var user User
    if db == nil {
        return user, errors.New("database connection is nil")
    }
    err := db.Where("telegram_id = ?", telegramID).First(&user).Error
    return user, err
}

func (u *User) UpdateTimeLimit(db *gorm.DB, newTimeLimit uint) error {
    if db == nil {
        return errors.New("database connection is nil")
    }
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
