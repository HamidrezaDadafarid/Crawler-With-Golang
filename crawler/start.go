package crawler

import (
	"fmt"
	"sync"

	"github.com/gocolly/colly"
)

func StartCrawler() {
	var (
		page      int
		waitGroup sync.WaitGroup
	)

	collyInstance := colly.NewCollector(
		colly.CacheDir("./crawler/cache"),
	)
	fmt.Scanf("%d", page)

	divarCrawler := NewDivarCrawler(page, &waitGroup, collyInstance)
	divarCrawler.Start()

}
