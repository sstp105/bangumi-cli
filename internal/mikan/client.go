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

type ClientOption func(*Client)

func WithClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = resty.NewWithClient(client)
	}
}

func NewClient(cfg ClientConfig, opts ...ClientOption) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid mikan client config: %s", err)
	}

	c := &Client{}
	c.client = resty.New()

	for _, opt := range opts {
		opt(c)
	}

	c.client.SetBaseURL(baseURL)
	c.client.SetHeaders(headers)

	cookies := []*http.Cookie{
		{
			Name:  identityCookieKey,
			Value: cfg.IdentityCookie,
		},
	}
	c.client.SetCookies(cookies)

	c.config = cfg

	return c, nil
}

// GetMyBangumi fetches the user's subscribed bangumi list from mikan.
func (c *Client) GetMyBangumi(year int, season season.Season) (string, error) {
	resp, err := c.client.R().
		Get(libs.FormatAPIPath(myBangumiPath, year, url.QueryEscape(string(season))))

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

func (c *ClientConfig) validate() error {
	if c.IdentityCookie == "" {
		return fmt.Errorf("%s is empty", identityCookieKey)
	}

	return nil
}
