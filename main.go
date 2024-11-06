package main

import (
	crawler "main/crawler"
)

func main() {
	// For testing db review for final codes
	// dbManager := database.GetInstnace()
	// dbManager.Db.AutoMigrate(&models.Advertisements{})
	crawler.StartCrawler(12)

}
