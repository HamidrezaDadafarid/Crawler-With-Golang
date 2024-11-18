package tests

import (
	"main/database"
	"main/models"
	"main/repository"
	"testing"
	"time"
)

func TestFilterRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	repo := repository.NewGormFilter(dbManager.Db)

	addedFilter, err := repo.Add(models.Filters{
		StartPurchasePrice: nil,
		EndPurchasePrice:   nil,
		StartRentPrice:     models.Ptr(uint(12)),
		EndRentPrice:       models.Ptr(uint(13)),
		StartMortgagePrice: models.Ptr(uint(14)),
		EndMortgagePrice:   models.Ptr(uint(15)),
		City:               models.Ptr("Test"),
		Neighborhood:       models.Ptr("TestNeighborhood"),
		StartArea:          models.Ptr(uint(10)),
		EndArea:            models.Ptr(uint(11)),
		StartNumberOfRooms: models.Ptr(uint(1)),
		EndNumberOfRooms:   models.Ptr(uint(1000)),
		CategoryPR:         models.Ptr(uint(1)),
		StartAge:           models.Ptr(uint(0)),
		EndAge:             models.Ptr(uint(100)),
		CategoryAV:         models.Ptr(uint(1)),
		StartFloorNumber:   models.Ptr(uint(2)),
		EndFloorNumber:     models.Ptr(uint(10)),
		Storage:            models.Ptr(true),
		Elevator:           models.Ptr(true),
		StartDate:          models.Ptr(time.Now()),
		EndDate:            models.Ptr(time.Now().Add(time.Hour * 48)),
	})
	if err != nil {
		t.Error("failed to add filter")
	}

	addedFilter.City = models.Ptr("Updated city")
	err = repo.Update(addedFilter)
	if err != nil {
		t.Error("failed to update filter")
	}

	updatedFilter, err := repo.Get([]uint{addedFilter.ID})

	if err != nil {
		t.Error("Failed to get filter")
	}

	if *updatedFilter[0].City != "Updated city" {
		t.Error("failed to update filter")
	}

	err = repo.Delete(addedFilter.ID)

	if err != nil {
		t.Error("failed to delete filter")
	}
}
