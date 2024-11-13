package tests

import (
	"fmt"
	"main/database"
	"main/models"
	"main/repository"
	"math/rand"
	"testing"
)

func TestUserAdsRepo(t *testing.T) {
	dbManager := database.GetInstnace()
	userRepo := repository.NewGormUser(dbManager.Db)
	adsRepo := repository.NewGormAd(dbManager.Db)
	userAdsRepo := repository.NewGormUser_Ad(dbManager.Db)

	addedUser, err := userRepo.Add(models.Users{
		TelegramId:       "telegram id",
		Role:             "User",
		Username:         "username",
		MaxSearchedItems: 10,
		TimeLimit:        100,
	})
	if err != nil {
		t.Error("failed to add user")
	}

	randomLink := fmt.Sprintf("link %d", rand.Int())
	ad, err := adsRepo.Add(models.Ads{
		Link:          randomLink,
		UniqueId:      "id",
		Longitude:     10,
		Latitude:      11,
		Description:   "some description",
		NumberOfViews: 11,
		SellPrice:     100,
		City:          "Shiraz",
		Neighborhood:  "mahale",
		Meters:        1000,
		NumberOfRooms: 10,
		CategoryPR:    1,
		Age:           10,
		CategoryAV:    2,
		FloorNumber:   1,
		Storage:       true,
		Elevator:      true,
		Title:         "Best",
	})
	if err != nil {
		t.Error("failed to add advertisement")
	}

	err = userAdsRepo.Add(models.Users_Ads{
		UserId:     addedUser.ID,
		AdId:       ad.ID,
		IsBookmark: false,
	})

	if err != nil {
		t.Error("failed to add userAdds")
	}

	_, err = userAdsRepo.GetByAdId([]uint{ad.ID})

	if err != nil {
		t.Error("failed to get userAdds")
	}

	err = userAdsRepo.Update(models.Users_Ads{UserId: addedUser.ID, AdId: ad.ID, IsBookmark: true})

	if err != nil {
		t.Error("failed to udpate userAdds")
	}

	_, err = userAdsRepo.GetByUserId([]uint{addedUser.ID})
	if err != nil {
		t.Error("failed to get userAdds")
	}

	err = userAdsRepo.Delete(addedUser.ID, ad.ID)
	if err != nil {
		t.Error("failed to delete userAdds")
	}
}
