package models

import (
	"sync"
)

type ICrawler interface {
	getTargets() //([]*Advertisement, error)
	getDetails() // ([]*Advertisement, error)
	validateJSON()
	sendDataToDB() // []*Advertisement // OR JSON
	Start()
} // TODO types should be added after ad structs are finished

type Crawler struct {
	page int
	wg   *sync.WaitGroup
}

func (c *Crawler) getTargets() {

}

func (c *Crawler) getDetails() {

}

func (c *Crawler) sendDataToDB() {
	// THIS IS THE SAME FOR ALL CRAWLERS
}

func (c *Crawler) validateJSON() {
	// DONE LATER
}

func (c *Crawler) Start() {
	c.getTargets()
	c.getDetails()
	c.validateJSON()
	c.sendDataToDB()
}

func NewCrawler(pg int, waitGroup *sync.WaitGroup) *Crawler {
	return &Crawler{page: pg, wg: waitGroup}
}
