package utils

import "time"

func IsToday(checkTime time.Time) bool {
	if checkTime.IsZero() {
		return false
	}
	currentTime := time.Now()
	return checkTime.Year() == currentTime.Year() &&
		checkTime.Month() == currentTime.Month() &&
		checkTime.Day() == currentTime.Day()
}
