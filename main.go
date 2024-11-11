package main

import (
	"log"
	database "main/database"
	models "main/models"
	"main/telegram"
	"os"
)

func main() {
	// For testing db review for final codes
	dbManager := database.GetInstnace()
	dbManager.Db.AutoMigrate(&models.Ads{}, &models.Filters{}, &models.Users{}, &models.Users_Ads{}, &models.WatchList{})

	telegramConfig := &telegram.TelegramConfig{
		Token: os.Getenv("TELEGRAM_TOKEN"),
	}

	telegramBot, err := telegram.NewTelegramBot(telegramConfig)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	telegramBot.Start()
}
