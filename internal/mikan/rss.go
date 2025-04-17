package mikan

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

// RSS represents the root rss structure.
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

// Channel represents the channel information in an RSS feed.
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents a single entry in an RSS feed.
type Item struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	GUID        string    `xml:"guid"`
	PubDate     string    `xml:"pubDate"`
	Enclosure   Enclosure `xml:"enclosure"`
}

// Enclosure represents media information within an RSS item.
type Enclosure struct {
	// URL is the torrent file downloadable url.
	URL string `xml:"url,attr"`

	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func (r *RSS) String() string {
	var buf bytes.Buffer
	for _, item := range r.Channel.Items {
		buf.WriteString(fmt.Sprintf("%s\n", item.Title))
	}
	return buf.String()
}

func (r *RSS) Torrents() []string {
	var urls []string
	for _, item := range r.Channel.Items {
		urls = append(urls, item.Enclosure.URL)
	}
	return urls
}

// Filter filters the RSS items based on the provided filters.
func Filter(r RSS, filters Filters) RSS {
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
