package mikan

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
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
	URL    string `xml:"url,attr"`
	Type   string `xml:"type,attr"`
	Length string `xml:"length,attr"`
}

func LoadRSS(url string) (*RSS, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", baseURL, url))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var rss RSS
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}

func (r *RSS) FilterInclude(filters []string) []Item {
	var items []Item

	for _, item := range r.Channel.Items {
		itemTitle := item.Title
		match := true

		for _, f := range filters {
			if !strings.Contains(strings.ToLower(itemTitle), strings.ToLower(f)) {
				match = false
				break
			}
		}

		if match {
			items = append(items, item)
		}
	}

	return items
}
