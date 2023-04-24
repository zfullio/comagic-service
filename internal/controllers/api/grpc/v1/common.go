package v1

import (
	"time"
)

func pbDateNormalize(s string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
