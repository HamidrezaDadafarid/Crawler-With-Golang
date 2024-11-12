package crawler

import (
	"errors"
	"fmt"
	"log"
	"main/models"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/go-rod/rod"
)

const url string = "https://divar.ir/s/iran/real-estate?page=%d"

type Advertisement = models.Ads

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
func (d *divar) GetTargets(page int, bInstance *rod.Browser) []*Advertisement {

	visit := fmt.Sprintf(url, page)

	log.Println("GRABBING TARGETS:", visit)

	var Ads []*Advertisement

	collector := bInstance.MustPage(visit)

	collector.MustWaitElementsMoreThan(`a[class=kt-post-card__action]`, 10)

	listOfLinks := collector.MustElements(`a[class=kt-post-card__action]`) // Grabs all tag <a>

	for _, aTag := range listOfLinks {
		link := getUniqueID(aTag.MustProperty("href").Str())
		ad := &Advertisement{UniqueID: link, Link: "divar", NumberOfViews: 0, PictureLink: ""}
		Ads = append(Ads, ad)
	}

	log.Println("SUCCESS FOR GRABBING")
	collector.Close()
	return Ads
}

func (d *divar) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup) {

	defer wg.Done()
	collector := bInstance.MustPage("https://divar.ir/v" + ad.UniqueID)

	collector.WaitStable(10)

	if ok, _, _ := collector.HasR(`a.kt-breadcrumbs__action`, `\u06a9\u0644\u0646\u06af\u06cc`); ok {
		ad.CategoryAV = ""
		return
	}

	ad.Title = collector.MustElement("h1").MustWaitVisible().MustText()

	ad.Description = collector.MustElement(`p.kt-description-row__text.kt-description-row__text--primary`).MustWaitVisible().MustText()

	ok, _, _ := collector.HasR(`td`, `^\u0627\u0646\u0628\u0627\u0631\u06cc$`) // checks for warehouse

	ad.Anbary = ok

	ok, _, _ = collector.HasR(`td`, `^\u0622\u0633\u0627\u0646\u0633\u0648\u0631$`) // checks for elevator

	ad.Elevator = ok

	ok, _, _ = collector.HasR(`td`, `^\u067e\u0627\u0631\u06a9\u06cc\u0646\u06af$`) // checks for parking

	ad.Parking = ok

	ok, _, _ = collector.HasR(`a.kt-breadcrumbs__action`, `\u0622\u067e\u0627\u0631\u062a\u0645\u0627\u0646`) //  checks for property type

	if ok {
		ad.CategoryAV = "apartment"
	} else {
		ad.CategoryAV = "villa"
	}

	if ok, _, _ = collector.HasR(`a.kt-breadcrumbs__action`, `\u0641\u0631\u0648\u0634`); ok {
		ad.CategoryPMR = `sale`
		ad.SellPrice = getSellPrice(collector)
		ad.MortgagePrice = -1
		ad.RentPrice = -1
	} else {
		ad.CategoryPMR = `rent`
		ad.SellPrice = -1
		// TODO FINISH MORTGAGE & RENT PRICE
	}

	ad.FloorNumber = getFloor(collector)

	surface, year, rooms := getSurfaceAndYearAndRooms(collector)
	ad.Meters = surface
	ad.Age = year
	ad.NumberOfRooms = rooms

	getLocation(ad, collector)

	ad.City = getCity(collector)
	ad.Mahale = getNeighbourhood(collector)

	ad.PictureLink = getPicture(collector)

	fmt.Println(ad)
}

func getPicture(collector *rod.Page) string {
	img := collector.MustElement(`img.kt-image-block__image.kt-image-block__image--fading`)
	return img.MustProperty(`src`).Str()
}

func cleanPrices(a string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(a, `٬`, ``), `تومان`, ``))
}

func getSellPrice(collector *rod.Page) int {
	if ok, section, _ := collector.HasR(`div.kt-base-row.kt-base-row--large.kt-unexpandable-row`, `\u0642\u06cc\u0645\u062a\u0020\u06a9\u0644`); ok {
		uncleanedPrice := section.MustElement(`p.kt-unexpandable-row__value`).MustText()
		return changeFarsiToEng(cleanPrices(uncleanedPrice))

	}
	return -1
}

func getFloor(collector *rod.Page) int {
	if ok, section, _ := collector.HasR(`div.kt-base-row.kt-base-row--large.kt-unexpandable-row`, `\u0637\u0628\u0642\u0647`); ok {
		floor := section.MustElement(`p.kt-unexpandable-row__value`).MustText()
		r := regexp.MustCompile(`[\\u06f1-\\u06f9]+`) // since numbers are persian
		return changeFarsiToEng(r.FindString(floor))
	}
	return -1
}

// These 2 functions should be merged later. they do the same job
func getNeighbourhood(collector *rod.Page) string {

	all, _ := getDistricts()
	lst := all.Districts

	for i := range lst {
		if ok, _, _ := collector.HasR(`div.kt-page-title__subtitle`, lst[i].Display); ok {
			return lst[i].Display
		}

	}
	return ""
}

func getCity(collector *rod.Page) string {

	all, _ := getCities()
	lst := all.Cities

	for i := range lst {
		if ok, _, _ := collector.HasR(`div.kt-page-title__subtitle.kt-page-title__subtitle--responsive-sized`, lst[i].Display); ok {
			return lst[i].Display
		}

	}
	return ""
}

func changeFarsiToEng(a string) int {
	runesOfString := []rune(a)
	var res []rune

	for i := range runesOfString {
		res = append(res, runesOfString[i]-1728)
	}

	val, err := strconv.Atoi(string(res))

	if err != nil {
		val = -1
	}

	return val
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

func getLatitudeAndLongitude(a string) (float64, float64) {
	r := regexp.MustCompile(`[0-9][0-9\.]+`)
	result := r.FindAllString(a, 2)
	lat, errlat := strconv.ParseFloat(result[0], 64)
	long, errlong := strconv.ParseFloat(result[1], 64)

	if errlat != nil || errlong != nil {
		log.Println(`FAILED TO CONVERT TO LATITUDE OR LONGTITUDE\nLINK:` + a)
		return -1, -1
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
