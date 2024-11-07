package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseRanges(input string) (int, int, error) {
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
	return minValue, maxValue, nil
}

// TODO: only year or date
func ParseDateRanges(input string) (string, string, error) {
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid format")
	}

	minDate, err := time.Parse("2006-01-02", parts[0])
	if err != nil {
		return "", "", fmt.Errorf("invalid min date format")
	}
	maxDate, err := time.Parse("2006-01-02", parts[1])
	if err != nil {
		return "", "", fmt.Errorf("invalid max date format")
	}

	return minDate.Format("2006-01-02"), maxDate.Format("2006-01-02"), nil
}
