package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

const url string = "https://divar.ir/s/iran/real-estate?page=%d"

func (c *divar) getTargets() ([]*Advertisement, error) {
	fmt.Println("hi nigger")
	var Ads []*Advertisement

	c.Collector.OnHTML("a[href]", func(h *colly.HTMLElement) {
		// Checks if the attribute class is same with the given class
		// This class is for post link
		if h.Attr("class") == "kt-post-card__action" {
			link := h.Request.AbsoluteURL(h.Attr("href"))
			Ads = append(Ads, &Advertisement{Link: link}) // Creates an Ad and adds a link to it
		}
	})
	c.Collector.OnRequest(func(r *colly.Request) {
		log.Println("GRABBED TARGETS FROM: ", r.URL) // logs the request url
	})

	c.Collector.Visit(fmt.Sprintf(url, c.Page)) // Starts sending request

	return Ads, nil
}

func (c *divar) getDetails(ad *Advertisement) {
	fmt.Println(ad.Link)
}
