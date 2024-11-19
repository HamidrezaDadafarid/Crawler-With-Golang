package crawler

import (
	"fmt"
	logg "main/log"
	"main/models"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

type sheypoor struct {
	*CrawlerAbstract
}

const urlSheypoor = "https://www.sheypoor.com/s/iran/real-estate?page=%d"

func NewSheypoorCrawler(wg *sync.WaitGroup, col *rod.Browser, s *Settings, metric *models.Metrics) *CrawlerAbstract {
	d := CrawlerAbstract{
		Wg:        wg,
		Collector: col,
		Settings:  s,
		Metric:    metric,
	}
	d.Crawler = &sheypoor{
		CrawlerAbstract: &d,
	}
	return &d
}

func (s *sheypoor) GetTargets(page int, bInstance *rod.Browser, lg logg.CrawlerLogger) []*Advertisement {
	var ads []*Advertisement

	collector := bInstance.MustPage(fmt.Sprintf(urlSheypoor, page))

	lg.InfoLogger.Println("[SHEYPOOR] fetching all targets...")

	collector.Mouse.MustScroll(0, 300)
	collector.MustWaitElementsMoreThan(`div.pt-4`, 8)

	listOfSections := collector.MustElements(`div.pt-4`)

	for _, elem := range listOfSections {
		links := elem.MustElements(`a`)

		for _, l := range links {
			unique := getUniqueID(l.MustProperty(`href`).Str())
			if unique != "" {
				ad := &Advertisement{Link: "sheypoor", UniqueId: unique, CategoryAV: 2}
				ads = append(ads, ad)
			}
		}
	}
	lg.InfoLogger.Println("[SHEYPOOR] fetched all targets")
	collector.Close()

	return ads
}

func (s *sheypoor) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup, lg logg.CrawlerLogger) {
	defer wg.Done()
	done := make(chan struct{})

	go func() {

		collector := bInstance.MustPage("https://www.sheypoor.com/v/" + ad.UniqueId)
		defer collector.MustClose()
		defer close(done)

		collector.WaitStable(15)

		if ok, _, _ := collector.HasR(`a.qL9GS`, `ویلا`); ok {
			ad.CategoryAV = 0
		} else if ok, _, _ = collector.HasR(`a.qL9GS`, `آپارتمان`); ok {
			ad.CategoryAV = 1
		} else {
			ad.CategoryAV = 2
			return
		}

		if ok, _, _ := collector.HasR(`a.qL9GS`, `رهن و اجاره`); ok {
			ad.CategoryPR = 1

		} else {
			ad.CategoryPR = 0
			sellPrice := changeFarsiToEng(collector.MustElement(`span.l29r1`).MustText())
			if sellPrice != -1 {
				ad.SellPrice = uint(sellPrice)
			}

		}

		ad.Title = collector.MustElement(`h1.mjNIv`).MustText()

		ad.Description = collector.MustElement(`div.MQJ5W`).MustText()

		a := []string{`انباری`, `آسانسور`, `پارکینگ`, `تعداد اتاق`, `متراژ`, `سن بنا`, `رهن`, `اجاره`}

		for _, key := range a {

			if ok, section, _ := collector.HasR(`div.C7Rh9`, key); ok {
				switch {

				case key == `انباری` || key == `آسانسور` || key == `پارکینگ`:

					if ok, _, _ := section.HasR(`p._874-x`, `^دارد$`); ok {

						switch key {
						case `انباری`:
							ad.Storage = ok
						case `آسانسور`:
							ad.Elevator = ok
						case `پارکینگ`:
							ad.Parking = ok
						}

					}

				default:

					uncleanedText := section.MustElement(`p._874-x`).MustText()
					cleanedText := changeFarsiToEng(cleanTexts(uncleanedText))

					if cleanedText != -1 {

						switch key {
						case `متراژ`:
							ad.Meters = uint(cleanedText)
						case `تعداد اتاق`:
							ad.NumberOfRooms = uint(cleanedText)
						case `سن بنا`:
							ad.Age = uint(cleanedText)
						case `رهن`:
							ad.MortgagePrice = uint(cleanedText)
						case `اجاره`:
							ad.RentPrice = uint(cleanedText)
						}

					}
				}
			}
		}

		patternFloor := regexp.MustCompile(`طبقه ملک: [0-9]+`)
		patternFloorNum := regexp.MustCompile(`[0-9]+`)

		flooruncleaned := patternFloor.FindString(ad.Description)
		floorcleaned := patternFloorNum.FindString(flooruncleaned)

		numeric, err := strconv.Atoi(floorcleaned)

		if err == nil {
			ad.FloorNumber = uint(numeric)
		}

		ad.City = getCity(collector, `div._3oBho`)
		ad.Neighborhood = getNeighbourhood(collector, `div._3oBho`)

		pic := collector.MustElement(`img`).MustProperty(`src`).Str()

		ad.PictureLink = pic

	}()
	select {
	case <-time.After(time.Second * 10):
		ad.CategoryAV = 2
		lg.ErrorLogger.Printf("[SHEYPOOR] failed to get advertisement %s\n", ad.UniqueId)
		s.Metric.FailRequestCount += 1
		return
	case <-done:
		lg.InfoLogger.Printf("[SHEYPOOR] successful advertisement %s\n", ad.UniqueId)
		s.Metric.SucceedRequestCount += 1
		return
	}

}
