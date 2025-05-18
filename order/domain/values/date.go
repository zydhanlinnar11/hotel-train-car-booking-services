package values

import (
	"errors"
	"time"
)

type Date string

var (
	ErrInvalidDate = errors.New("invalid date")
)

// Validate date format: DD-MM-YYYY
func NewDate(date string) (Date, error) {
	if _, err := time.Parse("02-01-2006", date); err != nil {
		return "", ErrInvalidDate
	}

	return Date(date), nil
}

func (d Date) String() string {
	return string(d)
}
