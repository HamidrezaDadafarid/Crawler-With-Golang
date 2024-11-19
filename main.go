package main

import (
	"fmt"
	"io"
	"log"
	"main/database"
	"main/models"
	"os"
	"os/signal"
	"time"

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
	go crawler.StartCrawler()

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt, os.Kill)
		<-sigchan

		crawl_log, _ := os.Open(`./log/crawler.log`)
		defer crawl_log.Close()

		dest, _ := os.Create(fmt.Sprintf(`./log/crawler%s.log`, time.Now().Format("2006-01-02")))
		defer dest.Close()

		io.Copy(dest, crawl_log)

		crawl_log, _ = os.Open(`./log/telegram.log`)
		defer crawl_log.Close()

		dest, _ = os.Create(fmt.Sprintf(`./log/telegram%s.log`, time.Now().Format("2006-01-02")))
		defer dest.Close()

		io.Copy(dest, crawl_log)

		log.Fatal("MANUAL INTERRUPTION / PROGRAM DEATH")
	}()
	select {}

}
