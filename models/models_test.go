package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Filter{}, &User{})
	return db, nil
}

func TestFilterValidation(t *testing.T) {
	filter := &Filter{
		StartPrice: 1000,
		EndPrice:   500,
		StartArea:  50,
		EndArea:    20,
	}
	err := filter.Validate()
	if err == nil {
		t.Error("expected validation error for invalid start and end ranges")
	}
}

func TestUserValidation(t *testing.T) {
	user := &User{}
	err := user.Validate()
	if err == nil {
		t.Error("expected validation error for missing TelegramID")
	}
}

func TestCreateAndRetrieveFilter(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	filter := &Filter{
		StartPrice:         1000,
		EndPrice:           5000,
		City:               "Test City",
		StartArea:          50,
		EndArea:            150,
		StartNumberOfRooms: 1,
		EndNumberOfRooms:   3,
		Category1:          "Residential",
	}

	if err := CreateFilter(db, filter); err != nil {
		t.Fatalf("failed to create filter: %v", err)
	}

	retrievedFilter, err := GetFilter(db, filter.Id)
	if err != nil {
		t.Fatalf("failed to retrieve filter: %v", err)
	}

	if retrievedFilter.City != filter.City {
		t.Errorf("expected city %s, got %s", filter.City, retrievedFilter.City)
	}
}

func TestCreateAndRetrieveUser(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	user := &User{
		TelegramID: "123456789",
		Role:       "Admin",
	}

	if err := CreateUser(db, user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	retrievedUser, err := GetUser(db, user.Id)
	if err != nil {
		t.Fatalf("failed to retrieve user: %v", err)
	}

	if retrievedUser.TelegramID != user.TelegramID {
		t.Errorf("expected TelegramID %s, got %s", user.TelegramID, retrievedUser.TelegramID)
	}
}
