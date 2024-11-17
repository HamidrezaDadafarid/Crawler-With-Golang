package crawler

import (
	"log"
	"main/database"
	"main/models"
	"main/repository"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type rodinstance struct {
}

func StartCrawler(page int, d bool, s bool) {
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

			settings, err := readConfig()

			if err != nil {
				log.Fatal("CRAWLER CONFIG ERROR")
			}

			if d {
				meric := models.Metrics{}
				divarCrawler := NewDivarCrawler(&waitGroup, rod, settings)
				t := time.Now()
				divarCrawler.Start()
				meric.TimeSpent = time.Since(t).Seconds()
			}

			if s {
				sheypoorCrawler := NewSheypoorCrawler(&waitGroup, *&rod, settings)
				sheypoorCrawler.Start()
			}

			time.Sleep(settings.Ticker * time.Minute)
		}
	}

}
