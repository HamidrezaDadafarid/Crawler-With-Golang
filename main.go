package main

import (
	"log"
	"main/database"
	"main/models"

	crawler "main/crawler"

	"github.com/joho/godotenv"
)

func main() {
	dbManager := database.GetInstnace()
	dbManager.Db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{}, &models.Users_Ads{}, &models.WatchList{})

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// telegramConfig := &telegram.TelegramConfig{
	// 	Token: os.Getenv("TELEGRAM_TOKEN"),
	// }

	// telegram, err := telegram.NewTelegramBot(telegramConfig)
	// if err != nil {
	// 	log.Fatalf("Error initializing Telegram bot: %v", err)
	// }

	// go telegram.Start()
	// For testing db review for final codes
	// dbManager := database.GetInstnace()
	// dbManager.Db.AutoMigrate(&models.Ads{})
	go crawler.StartCrawler()
	select {}

}
