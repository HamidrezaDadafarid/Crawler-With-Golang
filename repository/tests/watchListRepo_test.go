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
		TelegramId: "telegram id",
		Role:       "User",
	})
	if err != nil {
		t.Error("failed to add user")
	}

	addedFilter, err := filterRepo.Add(models.Filters{
		StartPurchasePrice: models.Ptr(uint(10)),
		EndPurchasePrice:   models.Ptr(uint(11)),
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
		CategoryAV:         models.Ptr(uint(2)),
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

	err = watchListRepo.Add(models.WatchList{
		UserID:   addedUser.ID,
		FilterId: addedFilter.ID,
		Time:     time.Duration(time.Second * 5),
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
