package model

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/bangumi"
	"strings"
)

type BangumiBase struct {
	// ID represents Mikan bangumi id.
	ID string `json:"id"`

	// Name represents the bangumi name in Simplified Chinese language.
	Name string `json:"name"`

	// Link represents the Mikan bangumi detail page url.
	Link string `json:"link"`
}

func (b BangumiBase) ConfigFileName() string {
	return fmt.Sprintf("%s.json", b.ID)
}

func (b BangumiBase) String() string {
	return fmt.Sprintf("ID:%s, Name:%s, Link:%s\n", b.ID, b.Name, b.Link)
}

func (b BangumiBase) SavePath() string {
	return fmt.Sprintf("/%s", b.Name)
}

type Bangumi struct {
	BangumiBase

	// BangumiID is a reference for  bangumi.tv id.
	BangumiID string `json:"bangumi_id"`

	// RSSLink represents the rss feed url for the subscribed fan-sub group.
	RSSLink string `json:"rss_link"`

	// Torrents holds a list of torrent urls that can be downloaded.
	Torrents []string `json:"torrents"`

	// Filters holds user configured torrent filters.
	Filters Filters `json:"filters"`

	// Episodes holds a list of episodes for the bangumi.
	Episodes []bangumi.Episode `json:"episodes"`
}

func (b Bangumi) StartEpisode() int {
	if len(b.Episodes) == 0 {
		return 1
	}
	return b.Episodes[0].Sort
}

func (b Bangumi) TorrentURLs() string {
	var builder strings.Builder
	for _, t := range b.Torrents {
		builder.WriteString(t + "\n")
	}
	return builder.String()
}
