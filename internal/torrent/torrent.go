package torrent

// Client defines a common interface for torrent clients.
type Client interface {
	// Add adds torrent urls to the client and download to dest folder.
	Add(urls string, dest string) error

	// Name returns the name of the underlying torrent client.
	Name() string
}
