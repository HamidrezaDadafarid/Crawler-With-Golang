package crawler

import (
	"fmt"
	"log"
	"main/models"
	"sync"

	"github.com/gocolly/colly"
)

const url string = "https://divar.ir/s/iran/real-estate?page=%d"

type Advertisement = models.Advertisement

type divar struct {
	models.CrawlerAbstract
}

func NewDivarCrawler(page int, wg *sync.WaitGroup, col *colly.Collector) *models.CrawlerAbstract {
	d := models.CrawlerAbstract{
		Crawler:   &divar{},
		Page:      page,
		Wg:        wg,
		Collector: col,
	}
	return &d
}

func (c *divar) GetTargets() ([]*Advertisement, error) {
	var Ads []*Advertisement
	fmt.Println("gg")
	c.Collector.OnHTML("a[href]", func(h *colly.HTMLElement) {
		// Checks if the attribute class is same with the given class
		// This class is for post link
		if h.Attr("class") == "kt-post-card__action" {
			link := h.Request.AbsoluteURL(h.Attr("href"))
			Ads = append(Ads, &Advertisement{Link: link}) // Creates an Ad and adds a link to it
		}
	})
	fmt.Println("g")
	c.Collector.OnRequest(func(r *colly.Request) {
		log.Println("GRABBED TARGETS FROM: ", r.URL) // logs the request url
	})
	c.Collector.Visit(fmt.Sprintf(url, c.Page)) // Starts sending request

	return Ads, nil
}

func (c *divar) GetDetails(ad *Advertisement) {
	fmt.Println(ad.Link)
}
