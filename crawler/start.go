package main

import (
	"fmt"
	"main/models"
	"sync"

	"github.com/gocolly/colly"
)

type Advertisement models.Advertisement

type divar struct {
	models.CrawlerAbstract
}

func NewDivarCrawler(page int, wg *sync.WaitGroup, col *colly.Collector) models.CrawlerAbstract {
	d := models.CrawlerAbstract{
		Crawler:   &divar{},
		Page:      page,
		Wg:        wg,
		Collector: col,
	}
	return d
}

func main() {
	var (
		page          int
		waitGroup     sync.WaitGroup
		collyInstance *colly.Collector = colly.NewCollector()
	)

	fmt.Scanf("%d", page)

	divarCrawler := NewDivarCrawler(page, &waitGroup, collyInstance)
	divarCrawler.Start()

}
