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
	"github.com/shirou/gopsutil/mem"
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

			log.Println("CRAWLER ON SLEEP")
			time.Sleep(settings.Ticker * time.Minute)
		}
	}

}
