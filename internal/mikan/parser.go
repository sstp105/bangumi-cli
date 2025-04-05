package mikan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/parser"
	htmlutil "html"
)

type BangumiBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Bangumi struct {
	BangumiBase
	BangumiID string   `json:"bangumi_id"`
	RSSLink   string   `json:"rss_link"`
	Torrents  []string `json:"torrents"`
	Filters   Filters  `json:"filters"`
}

type Filters struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

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

func ParseSubscribedRSSLink(doc *goquery.Document) (string, error) {
	// first seen group is the user subscribed fan-sub group
	subscribedGroup := doc.Find("div.subgroup-text").First()

	rssLink := subscribedGroup.Find("a.mikan-rss").AttrOr("href", "")
	if rssLink == "" {
		return "", fmt.Errorf("RSS link not found")
	}

	return rssLink, nil
}

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
