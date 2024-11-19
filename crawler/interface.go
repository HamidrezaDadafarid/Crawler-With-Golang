package crawler

import (
	"main/database"
	logg "main/log"
	"main/models"
	"main/repository"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type Ads = models.Ads

type Crawler interface {
	GetTargets(int, *rod.Browser, logg.CrawlerLogger) []*Ads
	GetDetails(*Ads, *rod.Browser, *sync.WaitGroup, logg.CrawlerLogger)
}

type CrawlerAbstract struct {
	Wg        *sync.WaitGroup
	Collector *rod.Browser
	Crawler   Crawler
	Settings  *Settings
	Metric    *models.Metrics
}

type Settings struct {
	Items   uint          `json:"CRAWLER_MAX_ITEMS"`
	Page    int           `json:"PAGE"`
	Timeout time.Duration `json:"CRAWLER_TIMEOUT"`
	Ticker  time.Duration `json:"CRAWLER_TICKER"`
	Logger  logg.CrawlerLogger
}

func (c *CrawlerAbstract) Start() {

	Ads := c.Crawler.GetTargets(c.Settings.Page, c.Collector, c.Settings.Logger)
	c.Metric.RequestCount = len(Ads) + 1
	c.Metric.SucceedRequestCount++
	c.iterateThroughAds(Ads)

	Ads = c.validateItems(Ads)
	c.sendDataToDB(Ads)

	c.Settings.Logger.InfoLogger.Println("added or updated items in database")
}

func (c *CrawlerAbstract) iterateThroughAds(Ads []*Ads) {
	defer c.Wg.Wait()

	timeout := time.NewTimer(c.Settings.Timeout * time.Second)

	for i := 0; i < len(Ads); i++ {

		select {
		case <-timeout.C:
			c.Settings.Logger.ErrorLogger.Println("CRAWL TIMEOUT REACHED")
			return
		case <-time.After(time.Millisecond * 2000):
			c.Wg.Add(1)
			go c.Crawler.GetDetails(Ads[i], c.Collector, c.Wg, c.Settings.Logger)

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
