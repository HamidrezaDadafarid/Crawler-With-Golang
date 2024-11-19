package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/crawler"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func ParseRanges(input string) (uint, uint, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("بازه قیمت باید به صورت دو عدد باشد که با - از هم جدا شده اند. دوباره تلاش کنید")
	}
	minValue, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("بازه قیمت باید به صورت دو عدد باشد که با - از هم جدا شده اند. دوباره تلاش کنید")
	}
	maxValue, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("بازه قیمت باید به صورت دو عدد باشد که با - از هم جدا شده اند. دوباره تلاش کنید")
	}
	if minValue < 0 || maxValue < 0 {
		return 0, 0, fmt.Errorf("لطفا اعداد منفی وارد نکنید. دوباره تلاش کنید")
	}
	if maxValue < minValue {
		return 0, 0, fmt.Errorf("عدد دوم بازه باید کوچک تر از عدد اول آن باشد! دوباره تلاش کنید")
	}
	return uint(minValue), uint(maxValue), nil
}

func ParseDateRanges(input string) (time.Time, time.Time, error) {
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("بازه تاریخ آگهی باید به صورت دو تاریخ به صورت روز-ماه-سال باشد که با یک فاصله از هم جدا شده اند")
	}

	minDate, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("تاریخ اول نامعتبر است. لطفا دوباره تلاش کنید")
	}
	maxDate, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("تاریخ دوم نا معتبر است. لطفا دوباره تلاش کنید")
	}

	return minDate, maxDate, nil
}

func GetEnvVariable(key string) (string, error) {

	err := godotenv.Load("./.env")

	if err != nil {
		return "", errors.New("error loading .env file")
	}

	return os.Getenv(key), nil
}

func SetCrawlerConfig(key string, val string) {
	var s *crawler.Settings

	s, _ = crawler.ReadConfig()

	conft, _ := strconv.Atoi(val)

	switch key {
	case "timeout":
		s.Timeout = time.Duration(conft)
	case "max_searched_items":
		s.Items = uint(conft)
	case "ticker":
		s.Ticker = time.Duration(conft)
	}

	d, _ := json.Marshal(s)

	f, _ := os.Create(`./config/config.json`)

	f.Write(d)
}
