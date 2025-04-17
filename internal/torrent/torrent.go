package torrent

// Client defines a common interface for torrent clients.
type Client interface {
	// Add adds torrent urls to the client and download to dest folder.
	// urls can be multiple torrent URLs separated with newlines.
	Add(urls string, dest string) error

	// Name returns the name of the underlying torrent client. e.g. qBittorrent.
	Name() string
}
