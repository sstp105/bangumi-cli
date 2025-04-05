package torrent

// Client defines a common interface for torrent clients.
type Client interface {
	// Add adds a torrent to the client using the provided link.
	Add(url string) error

	// Name returns the name of the underlying torrent client.
	Name() string
}
