package season

import "time"

// Season represents a season of the year.
type Season string

const (
	Spring Season = "春"
	Summer Season = "夏"
	Autumn Season = "秋"
	Winter Season = "冬"
)

// Now returns the current season.
func Now() Season {
	month := time.Now().Month()

	switch month {
	case time.January, time.February, time.March:
		return Winter
	case time.April, time.May, time.June:
		return Spring
	case time.July, time.August, time.September:
		return Summer
	case time.October, time.November, time.December:
		return Autumn
	}
	return Spring // unreachable
}
