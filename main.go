package main

import (
	database "main/database"
	"main/models"
)

func main() {
	// For testing db review for final codes
	dbManager := database.GetInstnace()
	dbManager.Db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{}, &models.Users_Ads{}, &models.WatchList{})
}
