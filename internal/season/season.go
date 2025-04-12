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
	return nowAt(time.Now())
}

// ID is the integer identifier for the season.
func (s Season) ID() ID {
	switch s {
	case Winter:
		return 1
	case Spring:
		return 2
	case Summer:
		return 3
	case Autumn:
		return 4
	}
	return -1
}

func (s Season) String() string {
	return string(s)
}

// nowAt returns the Season for a given time.
func nowAt(t time.Time) Season {
	month := t.Month()
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
	return ""
}
