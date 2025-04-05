package mikan

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

const (
	baseURL = "https://mikanani.me"

	identityCookieKey = ".AspNetCore.Identity.Application"
	
	myBangumiPath = "/Home/MyBangumi"
)

var (
	headers = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Referer":         baseURL,
		"Accept-Language": "en-US,en;q=0.9",
	}
)

type Client struct {
	client *resty.Client
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
	}, nil
}

func (c *Client) GetSubscribedAnimation() (string, error) {
	resp, err := c.client.R().
		Get(myBangumiPath)

	if err != nil {
		return "", fmt.Errorf("failed to fetch subscribed animation: %w", err)
	}

	return resp.String(), nil
}

func (c *Client) GetBangumi(path string) (string, error) {
	resp, err := c.client.R().
		Get(path)

	if err != nil {
		return "", fmt.Errorf("failed to fetch bangumi page: %w", err)
	}

	return resp.String(), nil
}

func (c *ClientConfig) validate() error {
	if c.IdentityCookie == "" {
		return fmt.Errorf("mikan identity cookie is empty")
	}

	return nil
}
