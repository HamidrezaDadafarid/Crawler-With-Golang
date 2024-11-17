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
	"github.com/shirou/gopsutil/cpu"
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
				metric := models.Metrics{}
				divarCrawler := NewDivarCrawler(&waitGroup, rod, settings)
				t := time.Now()
				divarCrawler.Start()
				metric.TimeSpent = time.Since(t).Seconds()
			}

			if s {
				percentChanel := make(chan bool, 1)
				result := 0.0
				counter := 0
				go func() {
					for {
						percent, _ := cpu.Percent(0, true)
						temp := 0.0
						for _, item := range percent {
							temp += item
						}
						result += temp / float64(len(percent))
						counter++
						time.Sleep(time.Millisecond * 100)
						select {
						case <-percentChanel:
							return
						}
					}
				}()
				metric := models.Metrics{}
				sheypoorCrawler := NewSheypoorCrawler(&waitGroup, *&rod, settings)
				t := time.Now()
				sheypoorCrawler.Start()
				metric.TimeSpent = time.Since(t).Seconds()
				percentChanel <- true
				metric.CpuUsage = result / float64(counter)

			}

			time.Sleep(settings.Ticker * time.Minute)
		}
	}

}
