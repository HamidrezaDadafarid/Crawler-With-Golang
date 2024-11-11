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
		StartPrice:         10,
		EndPrice:           11,
		City:               "Test",
		Neighborhood:       "TestNeighbor",
		SartArea:           10,
		EndArea:            10,
		StartNumberOfRooms: 1,
		EndNumberOfRooms:   1000,
		CategoryPMR:        1,
		StartAge:           0,
		EndAge:             100,
		CategoryAV:         2,
		StartFloorNumber:   0,
		EndFloorNumber:     10,
		Anbary:             true,
		Elevator:           true,
		StartDate:          time.Now(),
		EndDate:            time.Now().Add(time.Hour * 48),
	})
	if err != nil {
		t.Error("failed to add filter")
	}

	addedFilter.City = "Updated city"
	err = repo.Update(addedFilter)
	if err != nil {
		t.Error("failed to update filter")
	}

	updatedFilter, err := repo.Get([]uint{addedFilter.ID})

	if err != nil {
		t.Error("Failed to get filter")
	}

	if updatedFilter[0].City != "Updated city" {
		t.Error("failed to update filter")
	}

	err = repo.Delete(addedFilter.ID)

	if err != nil {
		t.Error("failed to delete filter")
	}
}
