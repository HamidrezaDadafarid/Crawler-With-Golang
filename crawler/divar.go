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
	collector.Close()
	return Ads
}

func (c *divar) GetDetails(ad *Advertisement, bInstance *rod.Browser, wg *sync.WaitGroup) {
	defer wg.Done()
	collector := bInstance.MustPage("https://divar.ir/v" + ad.UniqueID)

	collector.WaitStable(10)

	ad.Title = collector.MustElement("h1").MustWaitVisible().MustText()

	ad.Desc = collector.MustElement(`p.kt-description-row__text.kt-description-row__text--primary`).MustWaitVisible().MustText()

	ok, _, _ := collector.HasR(`td`, `^\u0627\u0646\u0628\u0627\u0631\u06cc$`) // checks for warehouse

	ad.Warehouse = ok

	ok, _, _ = collector.HasR(`td`, `^\u0622\u0633\u0627\u0646\u0633\u0648\u0631$`) // checks for elevator

	ad.Elevator = ok

	ok, _, _ = collector.HasR(`a.kt-breadcrumbs__action`, `\u0622\u067e\u0627\u0631\u062a\u0645\u0627\u0646`) //  checks for property type

	if ok {
		ad.TypeOfProperty = "apartment"
	} else {
		ad.TypeOfProperty = "villa"
	}

	surface, year, rooms := getSurfaceAndYearAndRooms(collector)

	ad.Surface = surface
	ad.YearOfBuild = year
	ad.RoomsCount = rooms

	getLocation(ad, collector)

	fmt.Println(ad)
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
