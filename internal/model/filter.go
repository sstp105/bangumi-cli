package model

// Filters represents the filter settings for including or excluding specific content from the rss.
type Filters struct {
	// Include holds a slice of string that must contain.
	Include []string `json:"include"`
}
