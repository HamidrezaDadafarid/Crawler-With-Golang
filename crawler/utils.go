package crawler

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
)

type City struct {
	Display string `json:"display"`
	Slug    string `json:"slug"`
}

type District struct {
	Display string `json:"display"`
	Slug    string `json:"slug"`
}

type Cities struct {
	Cities []City `json:"cities"`
}

type Districts struct {
	Districts []District `json:"districts"`
}

func cleanTexts(a string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(a, `٬`, ``), `تومان`, ``)), `.`, ``), `سال`, ``)
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
		return uncleaned
	}
	return ""
}

func getDistricts() (Districts, error) {
	file, err := os.Open("../database/districts.json")
	if err != nil {
		return Districts{}, err
	}
	defer file.Close()

	data, err := os.ReadFile("../database/districts.json")
	if err != nil {
		return Districts{}, err
	}

	var districts Districts
	err = json.Unmarshal(data, &districts)
	if err != nil {
		return Districts{}, err
	}

	return districts, nil

}

func getCities() (Cities, error) {
	file, err := os.Open("./database/city.json")
	if err != nil {
		return Cities{}, err
	}
	defer file.Close()

	data, err := os.ReadFile("./database/city.json")
	if err != nil {
		return Cities{}, err
	}

	var cities Cities
	err = json.Unmarshal(data, &cities)
	if err != nil {
		return Cities{}, err
	}

	return cities, nil

}

// returns ticker , timeout, max
func readConfig() (*Settings, error) {
	file, err := os.ReadFile(`./config/config.json`)

	if err != nil {
		return &Settings{}, err
	}

	var s Settings

	json.Unmarshal(file, &s)

	return &s, nil
}
