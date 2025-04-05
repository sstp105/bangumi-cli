package mikan

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	htmlutil "html"
	"regexp"
	"strings"
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

func ParseMyBangumiList(html string) ([]BangumiBase, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse bangumi document: %w", err)
	}

	var res []BangumiBase
	doc.Find(".sk-bangumi ul li").Each(func(i int, s *goquery.Selection) {
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

func ParseBangumiID(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	var bangumiID string
	doc.Find("p.bangumi-info a[href*='bgm.tv']").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			parts := strings.Split(href, "/")
			if len(parts) > 0 {
				bangumiID = parts[len(parts)-1]
			}
		}
	})

	if bangumiID == "" {
		return "", fmt.Errorf("bangumi ID not found")
	}

	return bangumiID, nil
}

func ParseSubscribedRSSLink(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	firstSubgroupText := doc.Find("div.subgroup-text").First()

	rssLink := firstSubgroupText.Find("a.mikan-rss").AttrOr("href", "")
	if rssLink == "" {
		return "", fmt.Errorf("RSS link not found")
	}

	return rssLink, nil
}

func parseMyBangumi(s *goquery.Selection) (*BangumiBase, error) {
	a := s.Find("a.an-text")
	link, exists := a.Attr("href")
	if !exists {
		return nil, fmt.Errorf("failed to parse mikan bangumi link")
	}

	name, exists := a.Attr("title")
	if !exists {
		return nil, fmt.Errorf("failed to parse mikan bangumi title")
	}
	name = htmlutil.UnescapeString(name)

	id, err := parseBangumiID(link)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mikan bangumi id: %w", err)
	}

	return &BangumiBase{
		ID:   id,
		Name: name,
		Link: link,
	}, nil
}

func parseBangumiID(s string) (string, error) {
	re := regexp.MustCompile(`/Home/Bangumi/(\d+)`)

	match := re.FindStringSubmatch(s)
	if len(match) < 2 {
		return "", fmt.Errorf("could not parse Bangumi ID from the input string")
	}

	return match[1], nil
}
