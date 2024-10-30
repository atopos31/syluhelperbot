package util

import "time"

func IsToday(dateStr string) (bool, error) {
	date, err := time.Parse(time.DateTime, dateStr)
	if err != nil {
		return false, err
	}
	today := time.Now()
	return date.Year() == today.Year() && date.Month() == today.Month() && date.Day() == today.Day(), nil
}
