package models

import (
	"log"
	"sync"
	"time"

	"math/rand"

	"github.com/gocolly/colly"
)

type Crawler interface {
	GetTargets(page int, collector *colly.Collector) ([]*Advertisement, error)
	GetDetails(*Advertisement)
} // TODO types should be added after ad structs are finished

type CrawlerAbstract struct {
	Page      int
	Wg        *sync.WaitGroup
	Collector *colly.Collector
	Crawler   Crawler
}

func (c *CrawlerAbstract) Start() {
	Ads, err := c.Crawler.GetTargets(c.Page, c.Collector)

	if err != nil {
		log.Fatal("CRAWLER ERROR", err)
		return
	}

	for _, ad := range Ads {
		c.Wg.Add(1)
		go c.Crawler.GetDetails(ad)

		randomSleep := rand.Intn(50) + 2 // To prevent rate-limits
		time.Sleep(time.Second * time.Duration(randomSleep))
	}
	c.validateJSON()
	c.sendDataToDB()
}

func (c *CrawlerAbstract) sendDataToDB() {

}

func (b *CrawlerAbstract) validateJSON() {

}
