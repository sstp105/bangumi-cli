package torrent

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/sstp105/bangumi-cli/internal/libs"
)

const (
	QBittorrentClientName = "qbittorrent"

	QBittorrentAuthCookieName = "SID"

	QBittorrentAPIAddPath   = "/api/v2/torrents/add"
	QBittorrentAPILoginPath = "/api/v2/auth/login"
)

var (
	ErrQBittorrentServerNotFound   = errors.New("qBittorrent port not found")
	ErrQBittorrentUserNameNotFound = errors.New("qBittorrent username not found")
	ErrQBittorrentPasswordNotFound = errors.New("qBittorrent password not found")
)

// QBittorrentClient represents a client for interacting with a qBittorrent web API.
type QBittorrentClient struct {
	// client is the Resty HTTP client used to make requests to the qBittorrent Web API.
	client *resty.Client

	// cookie stores the authentication cookie used in subsequent requests to the qBittorrent Web API.
	cookie string

	// config holds the configuration settings (port, username, and password) for the qBittorrent client.
	config QbittorrentClientConfig
}

// QbittorrentClientConfig represents the configuration required to interact with a qBittorrent client.
type QbittorrentClientConfig struct {
	// Server is the port where the qBittorrent Web API is running (e.g. localhost:8080).
	Server string

	// Username is the username used for authenticating against the qBittorrent Web API.
	Username string

	// Password is the password used for authenticating against the qBittorrent Web API.
	Password string
}

// NewQBittorrentClient https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#authentication
func NewQBittorrentClient(cfg QbittorrentClientConfig) (*QBittorrentClient, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("client config is invalid:%s", err)
	}

	c := resty.New()
	c.SetBaseURL(cfg.Server)

	authCookie, err := loginQbittorrent(c, cfg)
	if err != nil {
		return nil, fmt.Errorf("login qbittorrent failed:%s", err)
	}

	return &QBittorrentClient{
		client: c,
		cookie: authCookie,
		config: cfg,
	}, nil
}

// Add sends a request to the qBittorrent client to add a torrent using the provided magnet link or torrent URL.
func (q *QBittorrentClient) Add(urls string, dest string) error {
	resp, err := q.client.R().
		SetFormData(map[string]string{
			"urls":     urls,
			"savepath": dest,
			"paused":   "true",
		}).
		SetCookie(&http.Cookie{
			Name:  QBittorrentAuthCookieName,
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

// Name returns the name of the torrent client.
func (q *QBittorrentClient) Name() string {
	return QBittorrentClientName
}

// validate checks that the configuration fields for the qBittorrent client are set correctly.
func (q *QbittorrentClientConfig) validate() error {
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

// loginQbittorrent logs in to the qBittorrent webUI API and returns the authentication cookie on success.
// Reference: https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)#authentication
func loginQbittorrent(client *resty.Client, cfg QbittorrentClientConfig) (string, error) {
	resp, err := client.R().
		SetFormData(map[string]string{
			"username": cfg.Username,
			"password": cfg.Password,
		}).
		SetHeader("Referer", client.BaseURL).
		Post(QBittorrentAPILoginPath)

	if err != nil {
		return "", fmt.Errorf("login qbittorrent client error:%s", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return "", fmt.Errorf("failed to login to qbittorrent, status: %d", resp.StatusCode())
	}

	authCookie, err := libs.GetCookie(resp.Cookies(), QBittorrentAuthCookieName)
	if err != nil {
		return "", fmt.Errorf("failed to get %s from response header:%s", QBittorrentAuthCookieName, err)
	}

	return authCookie, nil
}
