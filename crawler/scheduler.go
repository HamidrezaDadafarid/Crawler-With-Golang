package crawler

import (
	"log"
	"main/database"
	logg "main/log"
	"main/models"
	"main/repository"
	"os"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type rodinstance struct {
}

func StartCrawler() {
	rodInstance := make(chan *rod.Browser, 1)
	gormMetric := repository.NewGormUMetric(database.GetInstnace().Db)

	logFile, err := os.Create(`./log/crawler.log`)

	if err != nil {
		log.Fatal(`ERORR: CRAWLER LOGGER CREATION`, err)
	}

	infoLogger := log.New(logFile, "INFO:", log.LstdFlags)
	errorLogger := log.New(logFile, "ERROR:", log.LstdFlags)

	crawlerLogger := logg.CrawlerLogger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}

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

		for {
			log.Println("STARTING CRAWLER")
			settings, err := readConfig()

			if err != nil {
				log.Fatal("ERROR: CRAWLER CONFIG")
			}

			settings.Logger = crawlerLogger

			finished := false
			result := 0.0
			memoryResult := 0.0
			counter := 0
			go func() {
				for {
					percent, _ := cpu.Percent(0, true)
					v, _ := mem.VirtualMemory()
					temp := 0.0
					for _, item := range percent {
						temp += item
					}
					memoryResult += v.UsedPercent
					result += temp / float64(len(percent))
					counter++
					time.Sleep(time.Millisecond * 100)
					if finished {
						return
					}
				}
			}()
			metric := models.Metrics{}
			divarCrawler := NewDivarCrawler(&waitGroup, rod, settings, &metric)
			t := time.Now()
			divarCrawler.Start()
			metric.TimeSpent = time.Since(t).Seconds()
			finished = true
			metric.CpuUsage = result / float64(counter)
			metric.RamUsage = memoryResult / float64(counter)
			gormMetric.Add(metric)

			finished = false
			result = 0.0
			counter = 0
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
					if finished {
						return
					}
				}
			}()
			metric = models.Metrics{}
			sheypoorCrawler := NewSheypoorCrawler(&waitGroup, *&rod, settings, &metric)
			t = time.Now()
			sheypoorCrawler.Start()
			metric.TimeSpent = time.Since(t).Seconds()
			finished = true
			metric.CpuUsage = result / float64(counter)
			gormMetric.Add(metric)

			settings.Logger.InfoLogger.Printf("Crawler sleeping for %d minutes\n", settings.Ticker)

			if 22 <= time.Now().Hour() {

				settings.Logger.InfoLogger.Println("Overnight crawler started...")

				gAd := repository.NewGormAd(database.GetInstnace().Db)

				ads, _ := repository.GetAds(database.GetInstnace().Db, 0)
				for i := 0; i < len(ads); i += 2 {
					if ads[i].Link == "divar" {
						waitGroup.Add(2)
						go divarCrawler.Crawler.GetDetails(&ads[i], rod, &waitGroup, settings.Logger)
						go divarCrawler.Crawler.GetDetails(&ads[i+1], rod, &waitGroup, settings.Logger)

					} else {
						waitGroup.Add(2)
						go sheypoorCrawler.Crawler.GetDetails(&ads[i], rod, &waitGroup, settings.Logger)
						go sheypoorCrawler.Crawler.GetDetails(&ads[i+1], rod, &waitGroup, settings.Logger)
					}
					waitGroup.Wait()
					gAd.Update(ads[i])
					gAd.Update(ads[i+1])

					time.Sleep(900 * time.Millisecond)
				}

				settings.Logger.InfoLogger.Println("Overnight crawler finished!")
				time.Sleep(1 * time.Hour)

			} else {
				time.Sleep(settings.Ticker * time.Minute)
			}
		}
	}

}
