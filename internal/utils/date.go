package utils

import "time"

func ConvertDate(date time.Time) string {
	return date.Format("2006-01-02")
}
