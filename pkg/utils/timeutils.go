package utils

import "time"

func StartTimeOfWeek(t time.Time) time.Time {
	days := time.Monday - t.Weekday()
	if days > 0 {
		days = -6
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, int(days))
}
