package mikan

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/sstp105/bangumi-cli/internal/model"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	GUID        string    `xml:"guid"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	// URL is the torrent file downloadable url.
	URL string `xml:"url,attr"`

	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func (r RSS) String() string {
	var buf bytes.Buffer
	for _, item := range r.Channel.Items {
		buf.WriteString(fmt.Sprintf("%s\n", item.Title))
	}
	return buf.String()
}

func (r RSS) Torrents() []model.Torrent {
	var torrents []model.Torrent
	for _, item := range r.Channel.Items {
		torrents = append(torrents, model.Torrent{
			Link: item.Enclosure.URL,
			Title: item.Title,
		})
	}
	return torrents
}

func (r RSS) TorrentURLs() []string {
	var urls []string
	for _, item := range r.Channel.Items {
		urls = append(urls, item.Enclosure.URL)
	}
	return urls
}

func (r RSS) Filter(filters model.Filters) RSS {
	var items []Item

	for _, item := range r.Channel.Items {
		match := true
		for _, f := range filters.Include {
			if !strings.Contains(strings.ToLower(item.Title), strings.ToLower(f)) {
				match = false
				break
			}
		}

		if match {
			items = append(items, item)
		}
	}

	r.Channel.Items = items

	return r
}
