package mikan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/model"
	htmlutil "html"
)

// ParseMyBangumiList parses all user subscribed bangumi.
func ParseMyBangumiList(doc *goquery.Document) ([]model.BangumiBase, error) {
	var res []model.BangumiBase

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

		id, err := libs.ParseSuffixID(href)
		if err != nil {
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
func parseMyBangumi(s *goquery.Selection) (*model.BangumiBase, error) {
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

	id, err := libs.ParseSuffixID(href)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bangumi id: %w", err)
	}

	return &model.BangumiBase{
		ID:   id,
		Name: name,
		Link: href,
	}, nil
}
