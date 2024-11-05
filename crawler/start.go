package main

import (
	"fmt"
	"main/models"
	"sync"
)

func main() {
	var (
		page      int
		waitGroup sync.WaitGroup
	)

	fmt.Scanf("%d", page)

	divarCrawler := models.NewCrawler(page, &waitGroup)
	divarCrawler.Start()

}
