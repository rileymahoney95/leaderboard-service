package utils

import (
	"time"
)

func ValidateDates(startDate *string, endDate *string) (*time.Time, *time.Time) {
	if startDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *startDate)
		if err != nil {
			return nil, nil
		}
		return &parsedDate, nil
	}

	if endDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *endDate)
		if err != nil {
			return nil, nil
		}
		return &parsedDate, nil
	}

	return nil, nil
}

