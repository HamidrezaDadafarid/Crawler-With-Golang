package main

import (
	"fmt"
	"log"
	"main/models"
	"github.com/gocolly/colly"
)

const url string = "https://divar.ir/s/iran/real-estate?page=%d"

func getTargets(crawler colly.Collector, current_pg int) []*Advertisement {

	var Ads []*Advertisement

	crawler.OnHTML("a[href]", func(h *colly.HTMLElement) {
		// Checks if the attribute class is same with the given class
		// This class is for post link
		if h.Attr("class") == "kt-post-card__action" {
			link := h.Request.AbsoluteURL(h.Attr("href"))
			Ads = append(Ads, &Advertisement{Link: link}) // Creates an Ad and adds a link to it
		}
	})

	crawler.OnRequest(func(r *colly.Request) {
		log.Println("GRABBED TARGETS FROM: ", r.URL) // logs the request url
	})

	crawler.Visit(fmt.Sprintf(url, current_pg)) // Starts sending request

	return Ads
}
