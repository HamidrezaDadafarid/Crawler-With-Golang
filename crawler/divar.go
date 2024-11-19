package crawler

import (
	"fmt"
	logg "main/log"
	"main/models"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

const urlDivar string = "https://divar.ir/s/iran/real-estate?page=%d"

type Advertisement = models.Ads

type divar struct {
	*CrawlerAbstract
}

func NewDivarCrawler(wg *sync.WaitGroup, col *rod.Browser, s *Settings, metric *models.Metrics) *CrawlerAbstract {
	d := CrawlerAbstract{
		Crawler:   &divar{},
		Wg:        wg,
		Collector: col,
		Settings:  s,
		Metric:    metric,
	}
	return &d
}

// Takes uniqueID from the link
func getUniqueID(a string) string {
	r := regexp.MustCompile(`\/[a-zA-Z0-9_-]+$|[0-9]{5,}`)

	return r.FindString(a)
}

// Grabs all targets needed to scrape.
func (d *divar) GetTargets(page int, bInstance *rod.Browser, lg logg.CrawlerLogger) []*Advertisement {

	visit := fmt.Sprintf(urlDivar, page)

	lg.InfoLogger.Println("[DIVAR] fetching all targets...")

	var Ads []*Advertisement

	collector := bInstance.MustPage(visit)
	defer collector.Close()

	collector.MustWaitElementsMoreThan(`a[class=kt-post-card__action]`, 10)

	listOfLinks := collector.MustElements(`a[class=kt-post-card__action]`) // Grabs all tag <a>

	for _, aTag := range listOfLinks {
		link := getUniqueID(aTag.MustProperty("href").Str())
		ad := &Advertisement{UniqueId: link, Link: "divar", CategoryAV: 2}
		Ads = append(Ads, ad)
	}

	lg.InfoLogger.Println("[DIVAR] fetched all targets")
	return Ads
}

func (d *divar) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup, lg logg.CrawlerLogger) {

	defer wg.Done()

	done := make(chan struct{})

	go func() {
		collector := bInstance.MustPage("https://divar.ir/v" + ad.UniqueId)

		defer collector.MustClose()
		defer close(done)

		collector.WaitStable(10)

		if ok, _, _ := collector.HasR(`a.kt-breadcrumbs__action`, `\u06a9\u0644\u0646\u06af\u06cc`); ok {
			ad.CategoryAV = 2
			return
		}

		ad.Title = collector.MustElement("h1").MustWaitVisible().MustText()

		ad.Description = collector.MustElement(`p.kt-description-row__text.kt-description-row__text--primary`).MustWaitVisible().MustText()

		ok, _, _ := collector.HasR(`td`, `^\u0627\u0646\u0628\u0627\u0631\u06cc$`) // checks for warehouse

		ad.Storage = ok

		ok, _, _ = collector.HasR(`td`, `^\u0622\u0633\u0627\u0646\u0633\u0648\u0631$`) // checks for elevator

		ad.Elevator = ok

		ok, _, _ = collector.HasR(`td`, `^\u067e\u0627\u0631\u06a9\u06cc\u0646\u06af$`) // checks for parking

		ad.Parking = ok

		ok, _, _ = collector.HasR(`a.kt-breadcrumbs__action`, `\u0622\u067e\u0627\u0631\u062a\u0645\u0627\u0646`) //  checks for property type

		if ok {
			ad.CategoryAV = 1
		} else {
			ad.CategoryAV = 0
		}

		if ok, _, _ = collector.HasR(`a.kt-breadcrumbs__action`, `\u0641\u0631\u0648\u0634`); ok {
			ad.CategoryPR = 0
			if val := getSellPrice(collector); val != -1 {
				ad.SellPrice = uint(val)
			}
		} else {
			ad.CategoryPR = 1
			m, r := getRentAndMortgagePrice(collector)

			if r != -1 {
				ad.RentPrice = uint(r)
			}
			if m != -1 {
				ad.MortgagePrice = uint(m)
			}

		}

		if val := getFloor(collector); val != -1 {
			ad.FloorNumber = uint(val)
		}

		surface, year, rooms := getSurfaceAndYearAndRooms(collector)
		if surface != -1 {
			ad.Meters = uint(surface)
		}
		if year != -1 {
			ad.Age = uint(year)
		}
		if rooms != -1 {
			ad.NumberOfRooms = uint(rooms)
		}

		lat, long := getLocation(collector)

		if lat != -1 {
			ad.Latitude = lat
		}

		if long != -1 {
			ad.Longitude = long
		}

		ad.City = getCity(collector, `kt-page-title__subtitle.kt-page-title__subtitle--responsive-sized`)
		ad.Neighborhood = getNeighbourhood(collector, `kt-page-title__subtitle.kt-page-title__subtitle--responsive-sized`)

		ad.PictureLink = getPicture(collector)
	}()

	select {
	case <-time.After(time.Second * 10):
		ad.CategoryAV = 2
		lg.ErrorLogger.Printf("[DIVAR] failed to get advertisement %s\n", ad.UniqueId)
		return
	case <-done:
		lg.InfoLogger.Printf("[DIVAR] successful advertisement %s\n", ad.UniqueId)
		return
	}
}

func getPicture(collector *rod.Page) string {
	img := collector.MustElement(`img.kt-image-block__image.kt-image-block__image--fading`)
	return img.MustProperty(`src`).Str()
}

func getChangabale(collector *rod.Page, elem *rod.Element) (int, int) {
	var (
		mortPrice int
		rentPrice int
	)

	mortelem := elem.MustElement(`td.kt-group-row-item.kt-group-row-item__value.kt-group-row-item--info-row`)
	mortPriceText := mortelem.MustText()

	if strings.Contains(mortPriceText, `میلیون`) {

		mortPrice = changeFarsiToEng(cleanTexts(strings.ReplaceAll(mortPriceText, `میلیون`, ``)))

	} else {

		mortPrice = changeFarsiToEng(cleanTexts(strings.ReplaceAll(mortPriceText, `میلیارد`, ``)))
	}

	rentPriceText := mortelem.MustNext().MustText()
	rentPrice = changeFarsiToEng(cleanTexts(strings.ReplaceAll(rentPriceText, `میلیون`, ``)))

	if rentPrice != -1 {
		rentPrice *= 1000000
	}

	return mortPrice * 1000000, rentPrice
}

func getNormal(collector *rod.Page) (int, int) {
	var (
		mort      int
		rent      int
		uncleaned string
	)

	uncleaned = getNumbersFromSections(`\u0648\u062f\u06cc\u0639\u0647`, collector) // mortgage price
	mort = changeFarsiToEng(cleanTexts(uncleaned))

	uncleaned = getNumbersFromSections(`\u0627\u062c\u0627\u0631\u0647\u0654\u0020\u0645\u0627\u0647\u0627\u0646\u0647`, collector) // rent price
	rent = changeFarsiToEng(cleanTexts(uncleaned))

	return mort, rent

}

func getRentAndMortgagePrice(collector *rod.Page) (int, int) {

	if ok, elem, _ := collector.Has(`div.convert-slider`); ok {
		mort, rent := getChangabale(collector, elem)
		return mort, rent
	}

	return getNormal(collector)

}

func getSellPrice(collector *rod.Page) int {
	uncleanedPrice := getNumbersFromSections(`\u0642\u06cc\u0645\u062a\u0020\u06a9\u0644`, collector)
	return changeFarsiToEng(cleanTexts(uncleanedPrice))
}

func getFloor(collector *rod.Page) int {
	floor := getNumbersFromSections(`\u0637\u0628\u0642\u0647`, collector)
	r := regexp.MustCompile(`[۱-۹]+`) // numbers are persian
	return changeFarsiToEng(r.FindString(floor))
}

func getSurfaceAndYearAndRooms(collector *rod.Page) (int, int, int) {

	conds := map[string]int{`متراژ`: -1, `ساخت`: -1, `اتاق`: -1}

	section := collector.MustElement(`table.kt-group-row`)
	all := strings.Split(section.MustText(), "\n")

	// covers every single condition that happens.
	for key, _ := range conds {

		for i := 0; i < len(all)/2; i++ {

			if key == all[i] {
				conds[key] = changeFarsiToEng(all[i+len(all)/2])
			}

		}
	}

	surface := conds[`متراژ`]
	rooms := conds[`اتاق`]
	year := conds[`ساخت`]

	return surface, year, rooms
}

func getLocation(collector *rod.Page) (float64, float64) {
	if ok, element, _ := collector.Has(`a.map-cm__attribution.map-cm__button`); ok {
		loc := element.MustProperty("href").Str()
		lat, long := GetLatitudeAndLongitude(loc)

		return lat, long
	}
	return -1, -1
}
