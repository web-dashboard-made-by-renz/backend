package utils

import (
	"fmt"
	"time"
)

func ParseTimestamp(dateStr string) (time.Time, error) {
	formats := []string{
		"1/2/2006 15:04:05",
		"1/2/2006 3:04:05 PM",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
		"01/02/2006",
		"2/1/2006 15:04:05",
		"02/01/2006 15:04:05",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
