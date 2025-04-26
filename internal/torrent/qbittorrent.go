package torrent

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
)

const (
	QBittorrentClientName = "qBittorrent"

	QBittorrentAuthCookieKey = "SID"

	QBittorrentAPIAddPath   = "/api/v2/torrents/add"
	QBittorrentAPILoginPath = "/api/v2/auth/login"
)

var (
	ErrQBittorrentServerNotFound   = errors.New("qBittorrent port not found")
	ErrQBittorrentUserNameNotFound = errors.New("qBittorrent username not found")
	ErrQBittorrentPasswordNotFound = errors.New("qBittorrent password not found")
)

type QBittorrentClient struct {
	client *resty.Client
	cookie string
	config QBittorrentClientConfig
}

type Option func(*QBittorrentClient)

func WithHTTPClient(client *http.Client) Option {
	return func(c *QBittorrentClient) {
		c.client = resty.NewWithClient(client)
	}
}

type QBittorrentClientConfig struct {
	Server   string
	Username string
	Password string
}

func NewQBittorrentClient(cfg QBittorrentClientConfig, opts ...Option) (*QBittorrentClient, error) {
	q := &QBittorrentClient{}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid qBittorrent client config: %w", err)
	}
	q.config = cfg

	client := resty.New()
	q.client = client

	for _, opt := range opts {
		opt(q)
	}

	q.client.SetBaseURL(cfg.Server)

	authCookie, err := q.authenticate()
	if err != nil {
		return nil, fmt.Errorf("qBittorrent authentication failed:%w", err)
	}
	q.cookie = authCookie

	return q, nil
}

func (q *QBittorrentClient) Add(urls string, dest string) error {
	resp, err := q.client.R().
		SetFormData(map[string]string{
			"urls":     urls,
			"savepath": dest,
			"paused":   "true", // do not start the task immediately
		}).
		SetCookie(&http.Cookie{
			Name:  QBittorrentAuthCookieKey,
			Value: q.cookie,
		}).
		Post(QBittorrentAPIAddPath)

	if err != nil {
		return fmt.Errorf("add torrent link error:%s", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to add torrent, status code:%s", resp.Status())
	}

	return nil
}

func (q *QBittorrentClient) Name() string {
	return QBittorrentClientName
}

// authenticate logs in to the qBittorrent webUI API and returns the authentication cookie on success.
// Reference: https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#authentication
func (q *QBittorrentClient) authenticate() (string, error) {
	resp, err := q.client.R().
		SetHeader("Referer", q.client.BaseURL).
		SetFormData(map[string]string{
			"username": q.config.Username,
			"password": q.config.Password,
		}).
		Post(QBittorrentAPILoginPath)

	if err != nil {
		return "", fmt.Errorf("authentication request failed:%w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to authenticate, status code: %d", resp.StatusCode())
	}

	authCookie, err := libs.GetCookie(resp.Cookies(), QBittorrentAuthCookieKey)
	if err != nil {
		return "", fmt.Errorf("%s not found from response header:%s", QBittorrentAuthCookieKey, err)
	}

	return authCookie, nil
}

func (q *QBittorrentClientConfig) validate() error {
	if q.Server == "" {
		return ErrQBittorrentServerNotFound
	}

	if q.Username == "" {
		return ErrQBittorrentUserNameNotFound
	}

	if q.Password == "" {
		return ErrQBittorrentPasswordNotFound
	}

	return nil
}
