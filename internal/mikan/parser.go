package mikan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/parser"
	htmlutil "html"
)

// BangumiBase represents the basic information about a Mikan bangumi.
// The struct is also used to be saved as local config for subsequent fetch/update.
type BangumiBase struct {
	// ID represents Mikan bangumi id.
	ID string `json:"id"`

	// Name represents the bangumi name in Simplified Chinese language.
	Name string `json:"name"`

	// Link represents the Mikan bangumi detail page url.
	Link string `json:"link"`
}

// Bangumi represents detailed information about a Mikan bangumi
type Bangumi struct {
	BangumiBase

	// BangumiID represents bangumi.tv id.
	BangumiID string `json:"bangumi_id"`

	// RSSLink represents the rss feed url for the bangumi.
	RSSLink string `json:"rss_link"`

	// Torrents holds a slice of torrent link or files for the bangumi.
	Torrents []string `json:"torrents"`

	// Filters holds user configured Filters.
	Filters Filters `json:"filters"`
}

// Filters represents the filter settings for including or excluding specific content from the rss.
type Filters struct {
	// Include holds a slice of string that must contain.
	Include []string `json:"include"`
}

// ParseMyBangumiList parses all subscribed bangumi.
func ParseMyBangumiList(doc *goquery.Document) ([]BangumiBase, error) {
	var res []BangumiBase

	selector := ".sk-bangumi ul li"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		b, err := parseMyBangumi(s)
		if err != nil {
			return
		}
		res = append(res, *b)
	})

	if len(res) == 0 {
		return nil, fmt.Errorf("failed to parse bangumi list, no results")
	}

	return res, nil
}

// ParseBangumiID parses the bangumi.tv id from the mikan bangumi detail page.
func ParseBangumiID(doc *goquery.Document) (string, error) {
	var bangumiID string

	selector := "p.bangumi-info a[href*='bgm.tv'], p.bangumi-info a[href*='bangumi.tv']"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exist := s.Attr("href")
		if !exist {
			return
		}

		id, err := parser.ParseSuffixID(href)
		if err != nil {
			log.Errorf("failed to parse bangumi id %s", href)
			return
		}

		bangumiID = id
	})

	if bangumiID == "" {
		return "", fmt.Errorf("bangumi ID not found")
	}

	return bangumiID, nil
}

// ParseSubscribedRSSLink pares the user subscribed fan-sub rss feed link from mikan bangumi page.
func ParseSubscribedRSSLink(doc *goquery.Document) (string, error) {
	// first seen group is the user subscribed fan-sub group
	subscribedGroup := doc.Find("div.subgroup-text").First()

	rssLink := subscribedGroup.Find("a.mikan-rss").AttrOr("href", "")
	if rssLink == "" {
		return "", fmt.Errorf("RSS link not found")
	}

	return rssLink, nil
}

// parseMyBangumi parses the bangumi element from the user subscribed bangumi list.
func parseMyBangumi(s *goquery.Selection) (*BangumiBase, error) {
	a := s.Find("a.an-text")
	href, exist := a.Attr("href")
	if !exist {
		return nil, fmt.Errorf("failed to parse bangumi link")
	}

	name, exists := a.Attr("title")
	if !exists {
		return nil, fmt.Errorf("failed to parse bangumi title")
	}
	name = htmlutil.UnescapeString(name) // mikan bangumi title are escaped

	id, err := parser.ParseSuffixID(href)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bangumi id: %w", err)
	}

	return &BangumiBase{
		ID:   id,
		Name: name,
		Link: href,
	}, nil
}
