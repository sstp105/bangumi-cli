package mikan

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
	"net/http"
)

const (
	baseURL = "https://mikanani.me"

	identityCookieKey = ".AspNetCore.Identity.Application"

	myBangumiPath libs.Path = "/Home/MyBangumi"
	bangumiPath   libs.Path = "/Home/Bangumi/%s"
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
	IdentityCookie string
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("create mikan client err: %s", err)
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

func (c *Client) GetMyBangumi() (string, error) {
	resp, err := c.client.R().
		Get(libs.FormatAPIPath(myBangumiPath))

	if err != nil {
		return "", fmt.Errorf("failed to fetch my bangumi page: %w", err)
	}

	return resp.String(), nil
}

func (c *Client) GetBangumi(id string) (string, error) {
	resp, err := c.client.R().
		Get(libs.FormatAPIPath(bangumiPath, id))

	if err != nil {
		return "", fmt.Errorf("failed to fetch bangumi page: %w", err)
	}

	return resp.String(), nil
}

func (c *ClientConfig) validate() error {
	if c.IdentityCookie == "" {
		return fmt.Errorf("%s is empty", identityCookieKey)
	}

	return nil
}
