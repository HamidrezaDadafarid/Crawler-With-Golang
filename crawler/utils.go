package crawler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

func cleanPrices(a string) string {
	return strings.ReplaceAll(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(a, `٬`, ``), `تومان`, ``)), `.`, ``)
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

func GetLatitudeAndLongitude(a string) (float64, float64) {
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

func getNumbersFromSections(reg string, collector *rod.Page) string {
	if ok, section, _ := collector.HasR(`div.kt-base-row.kt-base-row--large.kt-unexpandable-row`, reg); ok {
		uncleaned := section.MustElement(`p.kt-unexpandable-row__value`).MustText()
		fmt.Println(uncleaned, reg)
		return uncleaned
	}
	return ""
}
