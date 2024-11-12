package crawler

import (
	"encoding/json"
	"os"
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
