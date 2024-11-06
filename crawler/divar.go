package crawler

import (
	"errors"
	"fmt"
	"log"
	"main/models"
	"regexp"
	"strconv"
	"sync"

	"github.com/go-rod/rod"
)

const url string = "https://divar.ir/s/iran/real-estate?page=%d"

type Advertisement = models.Advertisement

type divar struct {
	*models.CrawlerAbstract
}

func NewDivarCrawler(page int, wg *sync.WaitGroup, col *rod.Browser) *models.CrawlerAbstract {
	d := models.CrawlerAbstract{
		Crawler:   &divar{},
		Page:      page,
		Wg:        wg,
		Collector: col,
	}
	return &d
}

// Takes uniqueID from the link
func getUniqueID(a string) string {
	r := regexp.MustCompile(`\/[a-zA-Z0-9_-]+$`)

	return r.FindString(a)
}

// Grabs all targets needed to scrape.
func (c *divar) GetTargets(page int, bInstance *rod.Browser) []*Advertisement {

	visit := fmt.Sprintf(url, page)

	log.Println("GRABBING TARGETS:", visit)

	var Ads []*Advertisement

	collector := bInstance.MustPage(visit)

	collector.MustWaitElementsMoreThan(`a[class=kt-post-card__action]`, 10)

	listOfLinks := collector.MustElements(`a[class=kt-post-card__action]`) // Grabs all tag <a>

	for _, aTag := range listOfLinks {
		link := getUniqueID(aTag.MustProperty("href").Str())
		ad := &Advertisement{UniqueID: link, Source: "divar"}
		Ads = append(Ads, ad)
	}

	log.Println("SUCCESS FOR GRABBING")

	return Ads
}

func getLatitudeAndLongitude(a string) (float64, float64) {
	r := regexp.MustCompile(`[0-9][0-9\.]+`)
	result := r.FindAllString(a, 2)
	lat, errlat := strconv.ParseFloat(result[0], 64)
	long, errlong := strconv.ParseFloat(result[1], 64)

	if errlat != nil || errlong != nil {
		log.Println(`FAILED TO CONVERT TO LATITUDE OR LONGTITUDE\nLINK:` + a)
		return lat, long
	}
	return lat, long
}

func getLocation(ad *Advertisement, collector *rod.Page) error {
	if ok, element, _ := collector.Has(`a.map-cm__attribution.map-cm__button`); ok {
		loc := element.MustProperty("href").Str()
		lat, long := getLatitudeAndLongitude(loc)

		ad.Latitude = lat
		ad.Longitude = long
		return nil
	}
	return errors.New("This ad has no location!")
}

func (c *divar) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup) {
	defer wg.Done()
	collector := bInstance.MustPage("https://divar.ir/v" + ad.UniqueID)

	collector.WaitStable(5)

	ad.Title = collector.MustElement("h1").MustWaitVisible().MustText()

	getLocation(ad, collector)

	fmt.Println(ad)
}
