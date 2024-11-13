package crawler

import (
	"fmt"
	"log"
	"main/models"
	"sync"

	"github.com/go-rod/rod"
)

type sheypoor struct {
	*models.CrawlerAbstract
}

const urlSheypoor = "https://www.sheypoor.com/s/iran/real-estate?page=%d"

func NewSheypoorCrawler(page int, wg *sync.WaitGroup, col *rod.Browser) *models.CrawlerAbstract {
	d := models.CrawlerAbstract{
		Crawler:   &sheypoor{},
		Page:      page,
		Wg:        wg,
		Collector: col,
	}
	return &d
}

func (s *sheypoor) GetTargets(page int, bInstance *rod.Browser) []*Advertisement {
	var ads []*Advertisement

	collector := bInstance.MustPage(fmt.Sprintf(urlSheypoor, page))

	log.Println("GRABBING TARGETS | [SHEYPOOR]")

	collector.MustWaitElementsMoreThan(`div.pt-4`, 8)

	listOfSections := collector.MustElements(`div.pt-4`)

	for _, elem := range listOfSections {
		links := elem.MustElements(`a`)

		for _, l := range links {
			unique := getUniqueID(l.MustProperty(`href`).Str())
			if unique != "" {
				ad := &Advertisement{Link: "sheypoor", UniqueId: unique}
				ads = append(ads, ad)
			}
		}
	}
	log.Println("SUCCESS GRABBING [SHEYPOOR]")
	collector.Close()

	return ads
}

func (s *sheypoor) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup) {

	defer wg.Done()

	collector := bInstance.MustPage("https://www.sheypoor.com/v/" + ad.UniqueId)

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

	fmt.Println(ad.Title)

	ad.Description = collector.MustElement(`div.MQJ5W`).MustText()

}
