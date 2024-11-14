package models

import (
    "errors"
    "strings"
    "gorm.io/gorm"
)

type User struct {
    UserID           uint   `json:"user_id" gorm:"primaryKey"`
    TelegramID       int64  `json:"telegram_id" gorm:"uniqueIndex"`
    Role             string `json:"role"`
    MaxSearchedItems uint   `json:"max_searched_items"`
    TimeLimit        uint   `json:"time_limit"`
    CountErrors      uint   `json:"count_errors"`
}

func (u *User) Normalize() {
    u.Role = strings.ToLower(u.Role)
}

func (u *User) Validate() error {
    if u.TelegramID == 0 {
        return errors.New("telegram_id cannot be zero")
    }
    if u.Role != "admin" && u.Role != "super_user" && u.Role != "user" {
        return errors.New("invalid role")
    }
    return nil
}

func (u *User) Create(db *gorm.DB) error {
    u.Normalize()
    if err := u.Validate(); err != nil {
        return err
    }
    return db.Create(u).Error
}

func (u *User) GetByID(db *gorm.DB, id uint) error {
    return db.First(u, id).Error
}

func (u *User) GetByTelegramID(db *gorm.DB, telegramID int64) error {
    return db.Where("telegram_id = ?", telegramID).First(u).Error
}

func (u *User) Update(db *gorm.DB) error {
    u.Normalize()
    if err := u.Validate(); err != nil {
        return err
    }
    return db.Save(u).Error
}

func (u *User) Delete(db *gorm.DB) error {
    return db.Delete(u).Error
}

func (u *User) DecreaseSearchLimit(db *gorm.DB) error {
    if u.MaxSearchedItems == 0 {
        return errors.New("search limit reached")
    }
    u.MaxSearchedItems--
    return db.Save(u).Error
}
