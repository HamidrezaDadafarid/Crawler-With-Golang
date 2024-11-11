package models_test

import (
    "testing"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "project_name/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }
    // Migrate both Filter and User models
    if err := db.AutoMigrate(&models.Filter{}, &models.User{}); err != nil {
        t.Fatalf("failed to migrate database: %v", err)
    }
    return db
}

func TestCreateFilter(t *testing.T) {
    db := setupTestDB(t)
    filter := models.Filter{
        City:            "Tehran",
        StartPrice:      100000,
        EndPrice:        500000,
        StartNumberOfRooms: 1,
        EndNumberOfRooms: 3,
    }

    if err := filter.Create(db); err != nil {
        t.Errorf("failed to create filter: %v", err)
    }
}

func TestGetFilter(t *testing.T) {
    db := setupTestDB(t)
    filter := models.Filter{City: "Tehran"}
    db.Create(&filter)

    var fetchedFilter models.Filter
    if err := fetchedFilter.Get(db, filter.FilterID); err != nil {
        t.Errorf("failed to get filter: %v", err)
    }

    if fetchedFilter.City != "Tehran" {
        t.Errorf("expected city 'Tehran', got %v", fetchedFilter.City)
    }
}

func TestUpdateFilter(t *testing.T) {
    db := setupTestDB(t)
    filter := models.Filter{City: "Tehran"}
    db.Create(&filter)

    filter.City = "Mashhad"
    if err := filter.Update(db); err != nil {
        t.Errorf("failed to update filter: %v", err)
    }

    var updatedFilter models.Filter
    db.First(&updatedFilter, filter.FilterID)
    if updatedFilter.City != "Mashhad" {
        t.Errorf("expected city 'Mashhad', got %v", updatedFilter.City)
    }
}

func TestDeleteFilter(t *testing.T) {
    db := setupTestDB(t)
    filter := models.Filter{City: "Tehran"}
    db.Create(&filter)

    if err := filter.Delete(db); err != nil {
        t.Errorf("failed to delete filter: %v", err)
    }

    var fetchedFilter models.Filter
    if err := fetchedFilter.Get(db, filter.FilterID); err == nil {
        t.Errorf("expected error, got nil")
    }
}
