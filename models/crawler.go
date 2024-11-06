package models

import (
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type Crawler interface {
	GetTargets(page int, collector *rod.Browser) []*Advertisement
	GetDetails(*Advertisement, *rod.Browser, *sync.WaitGroup)
} // TODO types should be added after ad structs are finished

type CrawlerAbstract struct {
	Page      int
	Wg        *sync.WaitGroup
	Collector *rod.Browser
	Crawler   Crawler
}

func (c *CrawlerAbstract) Start() {
	defer c.Collector.Close()
	Ads := c.Crawler.GetTargets(c.Page, c.Collector)

	for i := 0; i < len(Ads); i++ {
		c.Wg.Add(1)
		go c.Crawler.GetDetails(Ads[i], c.Collector, c.Wg)

		// randomSleep := rand.Intn(50) + 2 // To prevent rate-limits
		time.Sleep(time.Second * 2)
	}
	c.validateJSON()
	c.sendDataToDB()
}

func (c *CrawlerAbstract) sendDataToDB() {

}

func (b *CrawlerAbstract) validateJSON() {

}
