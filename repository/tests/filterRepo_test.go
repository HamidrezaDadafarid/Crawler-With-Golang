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

	}
}
