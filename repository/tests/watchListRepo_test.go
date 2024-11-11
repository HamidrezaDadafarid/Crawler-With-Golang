package tests

import (
	"main/database"
	"main/models"
	"main/repository"
	"testing"
	"time"
)

func TestWatchListRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	userRepo := repository.NewGormUser(dbManager.Db)
	filterRepo := repository.NewGormFilter(dbManager.Db)
	watchListRepo := repository.NewWatchList(dbManager.Db)

	addedUser, err := userRepo.Add(models.Users{
		TelegramId:       "telegram id",
		Role:             "User",
		MaxSearchedItems: 10,
		TimeLimit:        100,
	})
	if err != nil {
		t.Error("failed to add user")
	}

	addedFilter, err := filterRepo.Add(models.Filters{
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

	err = watchListRepo.Add(models.WatchList{
		UserID:   addedUser.ID,
		FilterId: addedFilter.ID,
		Time:     time.Now(),
	})

	if err != nil {
		t.Error("failed to add watchlist")
	}

	_, err = watchListRepo.GetByFilterId([]uint{addedFilter.ID})

	if err != nil {
		t.Error("failed to get watchList")
	}

	err = watchListRepo.Update(models.WatchList{UserID: addedUser.ID, FilterId: addedFilter.ID})

	if err != nil {
		t.Error("failed to udpate watchList")
	}

	_, err = watchListRepo.GetByUserId([]uint{addedUser.ID})
	if err != nil {
		t.Error("failed to get watchList")
	}

	err = watchListRepo.Delete(addedUser.ID, addedFilter.ID)
	if err != nil {
		t.Error("failed to delete watchList")
	}
}
