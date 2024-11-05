package main

import (
	"fmt"
	"main/crawler"
	"sync"

	"github.com/gocolly/colly"
)

func main() {
	var (
		page      int
		waitGroup sync.WaitGroup
	)

	collyInstance := colly.NewCollector(
		colly.CacheDir("./crawler/cache"),
	)
	fmt.Scanf("%d", page)

	divarCrawler := crawler.NewDivarCrawler(page, &waitGroup, collyInstance)
	
	divarCrawler.Start()

}
