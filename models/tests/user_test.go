package models_test

import (
    "testing"
    "project_name/models"
)

func TestCreateUser(t *testing.T) {
    db := setupTestDB(t)
    user := models.User{
        TelegramID: 12345678,
        Role:       "admin",
    }

    if err := user.Create(db); err != nil {
        t.Errorf("failed to create user: %v", err)
    }
}

func TestGetUser(t *testing.T) {
    db := setupTestDB(t)
    user := models.User{TelegramID: 12345678, Role: "admin"}
    db.Create(&user)

    var fetchedUser models.User
    if err := fetchedUser.Get(db, user.UserID); err != nil {
        t.Errorf("failed to get user: %v", err)
    }

    if fetchedUser.TelegramID != 12345678 {
        t.Errorf("expected TelegramID '12345678', got %v", fetchedUser.TelegramID)
    }
}

func TestUpdateUser(t *testing.T) {
    db := setupTestDB(t)
    user := models.User{TelegramID: 12345678, Role: "admin"}
    db.Create(&user)

    user.Role = "user"
    if err := user.Update(db); err != nil {
        t.Errorf("failed to update user: %v", err)
    }

    var updatedUser models.User
    db.First(&updatedUser, user.UserID)
    if updatedUser.Role != "user" {
        t.Errorf("expected role 'user', got %v", updatedUser.Role)
    }
}

func TestDeleteUser(t *testing.T) {
    db := setupTestDB(t)
    user := models.User{TelegramID: 12345678, Role: "admin"}
    db.Create(&user)

    if err := user.Delete(db); err != nil {
        t.Errorf("failed to delete user: %v", err)
    }

    var fetchedUser models.User
    if err := fetchedUser.Get(db, user.UserID); err == nil {
        t.Errorf("expected error, got nil")
    }
}
