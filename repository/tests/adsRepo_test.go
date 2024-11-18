package tests

import (
	"main/database"
	"main/models"
	"main/repository"
	"testing"
)

func TestAdsRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	adRepo := repository.NewGormAd(dbManager.Db)
	ad, err := adRepo.Add(models.Ads{
		Link:          "divar link2",
		UniqueId:      "id",
		Longitude:     10,
		Latitude:      11,
		Description:   "some description",
		NumberOfViews: 11,
		SellPrice:     100,
		City:          "Shiraz",
		Mahale:        "mahale",
		Meters:        1000,
		NumberOfRooms: 10,
		CategoryPR:    1,
		Age:           10,
		CategoryAV:    2,
		FloorNumber:   1,
		Anbary:        true,
		Elevator:      true,
		Title:         "Best",
	})
	if err != nil {
		t.Errorf("failed to add advertisement")
	}

	updateErr := adRepo.Update(models.Ads{
		ID:            ad.ID,
		Link:          "divar link",
		UniqueId:      "id",
		Longitude:     10,
		Latitude:      11,
		Description:   "some description updated",
		NumberOfViews: 11,
		SellPrice:     100,
		City:          "Shiraz",
		Mahale:        "mahale",
		Meters:        1000,
		NumberOfRooms: 10,
		CategoryPR:    1,
		Age:           10,
		CategoryAV:    2,
		FloorNumber:   1,
		Anbary:        true,
		Elevator:      true,
		Title:         "Best",
	})

	if updateErr != nil {
		t.Errorf("failed to update advertisement")
	}

	byIdAd, err := adRepo.GetById([]uint{ad.ID})
	if err != nil || byIdAd[0].ID != ad.ID {
		t.Errorf("failed to get by id advertisement")
	}
	if byIdAd[0].Description != "some description updated" {
		t.Errorf("failed to update advertisement")
	}

	byFilter, err := adRepo.Get(models.Filters{
		Elevator: false,
	})

	if err != nil {
		t.Error("Failed to get by filter")
	}

	for _, item := range byFilter {
		if item.Elevator == true {
			t.Error("Failed to get by filter")
		}
	}
}
