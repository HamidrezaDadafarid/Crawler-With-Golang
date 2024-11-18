package crawler

import (
	"fmt"
	"log"
	"main/database"
	"main/models"
	"main/repository"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type Ads = models.Ads

type Crawler interface {
	GetTargets(page int, collector *rod.Browser) []*Ads
	GetDetails(*Ads, *rod.Browser, *sync.WaitGroup)
}

type CrawlerAbstract struct {
	Wg        *sync.WaitGroup
	Collector *rod.Browser
	Crawler   Crawler
	Settings  *Settings
}

type Settings struct {
	Items   uint          `json:"CRAWLER_MAX_ITEMS"`
	Page    int           `json:"PAGE"`
	Timeout time.Duration `json:"CRAWLER_TIMEOUT"`
	Ticker  time.Duration `json:"CRAWLER_TICKER"`
}

func (c *CrawlerAbstract) Start() {

	Ads := c.Crawler.GetTargets(c.Settings.Page, c.Collector)

	c.iterateThroughAds(Ads)

	Ads = c.validateItems(Ads)
	fmt.Println(len(Ads))
	// c.sendDataToDB(Ads)
}

func (c *CrawlerAbstract) iterateThroughAds(Ads []*Ads) {
	defer c.Wg.Wait()

	timeout := time.NewTimer(c.Settings.Timeout * time.Second)

	for i := 0; i < len(Ads); i++ {

		select {
		case <-timeout.C:
			log.Println("crawler timeout [SHEYPOOR]")
			return
		case <-time.After(time.Millisecond * 2000):
			c.Wg.Add(1)
			go c.Crawler.GetDetails(Ads[i], c.Collector, c.Wg)

		}

	}

}

func (c *CrawlerAbstract) sendDataToDB(a []*Ads) {
	g := repository.NewGormAd(database.GetInstnace().Db)
	for _, ad := range a {
		g.Add(*ad)
	}
}

func (b *CrawlerAbstract) validateItems(adList []*Ads) []*Ads {
	var res []*Ads
	for i := range adList {
		if adList[i].CategoryAV != 2 {
			res = append(res, adList[i])
		}
	}
	return res

}
