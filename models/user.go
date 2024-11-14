package models

import (
    "errors"
    "strings"
    "gorm.io/gorm"
)

type Users struct {
    gorm.Model
    ID               uint     `gorm:"primaryKey;autoIncrement"`
    TelegramId       string
    Role             string
    MaxSearchedItems uint
    TimeLimit        uint
}

func (u *Users) Normalize() {
    u.Role = strings.ToLower(u.Role)
}

func (u *Users) Validate() error {
    if u.TelegramId == "" {
        return errors.New("telegram_id cannot be empty")
    }
    if u.Role != "admin" && u.Role != "super_user" && u.Role != "user" {
        return errors.New("invalid role")
    }
    return nil
}
