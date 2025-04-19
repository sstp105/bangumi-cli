package mikan

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"github.com/sstp105/bangumi-cli/internal/season"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL = "https://mikanani.me"

	identityCookieKey = ".AspNetCore.Identity.Application"

	myBangumiPath libs.APIPath = "/Home/BangumiCoverFlow?year=%d&seasonStr=%s"
	bangumiPath   libs.APIPath = "/Home/Bangumi/%s"
)

var headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Referer":         baseURL,
	"Accept-Language": "en-US,en;q=0.9",
}

type Client struct {
	client *resty.Client
	config ClientConfig
}

type ClientConfig struct {
	// IdentityCookie represents the authentication cookie in Mikan
	IdentityCookie string
}

type ClientOption func(*clientOption)

type clientOption struct {
	year   int
	season season.Season
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid mikan client config: %s", err)
	}

	c := resty.New()
	c.SetBaseURL(baseURL)
	c.SetHeaders(headers)

	cookies := []*http.Cookie{
		{
			Name:  identityCookieKey,
			Value: cfg.IdentityCookie,
		},
	}
	c.SetCookies(cookies)

	return &Client{
		client: c,
		config: cfg,
	}, nil
}

// GetMyBangumi fetches the user's subscribed bangumi list from Mikan.
// By default, it returns the latest season. Use options to specify a year and season.
func (c *Client) GetMyBangumi(opts ...ClientOption) (string, error) {
	defaultYear := time.Now().Year()
	defaultSeason := season.Now()
	opt := &clientOption{
		year:   defaultYear,
		season: defaultSeason,
	}

	for _, o := range opts {
		o(opt)
	}

	resp, err := c.client.R().
		Get(libs.FormatAPIPath(myBangumiPath, opt.year, url.QueryEscape(string(opt.season))))

	if err != nil {
		return "", fmt.Errorf("failed to fetch mikan my bangumi page: %w", err)
	}

	return resp.String(), nil
}

// GetBangumi returns the html response from Mikan bangumi detail page.
func (c *Client) GetBangumi(id string) (string, error) {
	resp, err := c.client.R().Get(libs.FormatAPIPath(bangumiPath, id))

	if err != nil {
		return "", fmt.Errorf("failed to fetch mikan bangumi%s: %w", id, err)
	}

	return resp.String(), nil
}

func (c *Client) LoadRSS(url string) (*RSS, error) {
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rss: %w", err)
	}

	var rss RSS
	decoder := xml.NewDecoder(bytes.NewReader(resp.Body()))
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rss: %w", err)
	}

	return &rss, nil
}

func WithYearAndSeason(year int, season season.Season) ClientOption {
	return func(opt *clientOption) {
		opt.year = year
		opt.season = season
	}
}

func (c *ClientConfig) validate() error {
	if c.IdentityCookie == "" {
		return fmt.Errorf("%s is empty", identityCookieKey)
	}

	return nil
}
