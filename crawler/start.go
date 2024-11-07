package crawler

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type rodinstance struct {
}

func StartCrawler(page int) {
	rodInstance := make(chan *rod.Browser, 1)

	var (
		waitGroup sync.WaitGroup
	)

	go func() {
		rodInstance <- rod.New().MustConnect().Trace(true)
	}()

	select {
	case <-time.After(time.Second * 5):
		log.Println("TIMEOUT WHEN STARTING CRAWLER")
	case rod := <-rodInstance:

		go func() {
			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, os.Interrupt)
			<-sigchan

			log.Println("MANUAL INTERRUPTION / PROGRAM DEATH")
			rod.Close()
			os.Exit(0)
		}()

		divarCrawler := NewDivarCrawler(page, &waitGroup, rod)
		divarCrawler.Start()
	}

}
