package crawler

import (
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
	Page      int
	Wg        *sync.WaitGroup
	Collector *rod.Browser
	Crawler   Crawler
}

func (c *CrawlerAbstract) Start(t time.Duration) {
	defer c.Collector.Close()
	Ads := c.Crawler.GetTargets(c.Page, c.Collector)

	timeout := time.NewTimer(t)

	for i := 0; i < len(Ads); i++ {

		select {
		case <-timeout.C:
			break
		case <-time.After(time.Millisecond * 1500):
			c.Wg.Add(1)
			go c.Crawler.GetDetails(Ads[i], c.Collector, c.Wg)
		}

	}
	c.Wg.Wait()
	Ads = c.validateItems(Ads)
	c.sendDataToDB(Ads)
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
	return adList

}
