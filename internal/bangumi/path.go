package bangumi

import "fmt"

// Path represents a bangumi API route in string.
type Path string

const (
	// getUserCollectionsPath is path for fetching a user's collections.
	getUserCollectionsPath Path = "/v0/users/%s/collections"
)

// getUserCollectionsURL generates the API URL for fetching user collections.
func getUserCollectionsURL(username string) string {
	return apiURL(getUserCollectionsPath, username)
}

// apiURL returns the relative API URL given the path and provided arguments (e.g. query params).
func apiURL(p Path, args ...interface{}) string {
	return fmt.Sprintf(string(p), args...)
}
