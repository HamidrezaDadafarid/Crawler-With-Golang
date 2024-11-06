package crawler

import (
	"sync"

	"github.com/go-rod/rod"
)

func StartCrawler(page int) {
	var (
		waitGroup   sync.WaitGroup
		rodInstance *rod.Browser = rod.New().MustConnect().Trace(true)
	)

	divarCrawler := NewDivarCrawler(page, &waitGroup, rodInstance)
	divarCrawler.Start()

}
