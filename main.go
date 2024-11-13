package main

import (
	crawler "main/crawler"
)

func main() {
	// For testing db review for final codes
<<<<<<< HEAD
	dbManager := database.GetInstnace()
	dbManager.Db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{}, &models.Users_Ads{}, &models.WatchList{})
=======
	// dbManager := database.GetInstnace()
	// dbManager.Db.AutoMigrate(&models.Advertisements{})
	crawler.StartCrawler(12)

>>>>>>> b298c8cb38b55b95c70579f48c8905df52180ae6
}
