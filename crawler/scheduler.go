package crawler

import (
	"log"
	"main/database"
	"main/repository"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type rodinstance struct {
}

func StartCrawler() {
	rodInstance := make(chan *rod.Browser, 1)
	gormMetric := repository.NewGormUMetric(database.GetInstnace().Db)

	var (
		waitGroup sync.WaitGroup
	)

	go func() {
		rodInstance <- rod.New().MustConnect()
	}()

	select {
	case <-time.After(time.Second * 15):
		log.Println("TIMEOUT WHEN STARTING CRAWLER")
	case rod := <-rodInstance:

		go func() {
			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, os.Interrupt)
			<-sigchan

			log.Println("MANUAL INTERRUPTION / PROGRAM DEATH")
			rod.MustClose()
			os.Exit(0)
		}()

		for {
			log.Println("STARTING CRAWLER")
			settings, err := readConfig()

			if err != nil {
				log.Fatal("CRAWLER CONFIG ERROR")
			}

			divarCrawler := NewDivarCrawler(&waitGroup, rod, settings)
			divarCrawler.Start()

			sheypoorCrawler := NewSheypoorCrawler(&waitGroup, *&rod, settings)
			sheypoorCrawler.Start()

			log.Println("CRAWLER ON SLEEP")
			time.Sleep(settings.Ticker * time.Minute)
		}
	}

}
