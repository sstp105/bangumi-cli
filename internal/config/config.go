package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/sstp105/bangumi-cli/internal/log"
	"github.com/sstp105/bangumi-cli/internal/mikan"
	"github.com/sstp105/bangumi-cli/internal/torrent"
)

type config struct {
	port                string
	bangumiClientID     string
	bangumiClientSecret string
	qbittorrentConfig   torrent.QBittorrentClientConfig
	mikanClientConfig   mikan.ClientConfig
}

var cfg config

func init() {
	cfg.port = os.Getenv(PortKey)
	cfg.bangumiClientID = os.Getenv(BangumiClientIDKey)
	cfg.bangumiClientSecret = os.Getenv(BangumiClientSecretKey)
	cfg.qbittorrentConfig = torrent.QBittorrentClientConfig{
		Server:   os.Getenv(QBittorrentServerKey),
		Username: os.Getenv(QBittorrentUserNameKey),
		Password: os.Getenv(QBittorrentPasswordKey),
	}
	cfg.mikanClientConfig = mikan.ClientConfig{
		IdentityCookie: os.Getenv(MikanIdentityCookieKey),
	}

	if err := cfg.validate(); err != nil {
		log.Warnf("config is invalid:%s", err)
	}
}

func Port() string {
	return cfg.port
}

func BangumiClientID() string {
	return cfg.bangumiClientID
}

func BangumiClientSecret() string {
	return cfg.bangumiClientSecret
}

func QBittorrentConfig() torrent.QBittorrentClientConfig {
	return cfg.qbittorrentConfig
}

func MikanClientConfig() mikan.ClientConfig {
	return cfg.mikanClientConfig
}

func LocalServerAddress() string {
	return fmt.Sprintf("http://localhost:%s", cfg.port)
}

func (c config) validate() error {
	var errs []string

	check := func(name, val string) {
		if val == "" {
			errs = append(errs, fmt.Sprintf("%s is empty", name))
		}
	}

	check(PortKey, c.port)
	check(BangumiClientIDKey, c.bangumiClientID)
	check(BangumiClientSecretKey, c.bangumiClientSecret)
	check(QBittorrentServerKey, c.qbittorrentConfig.Server)
	check(QBittorrentUserNameKey, c.qbittorrentConfig.Username)
	check(QBittorrentPasswordKey, c.qbittorrentConfig.Password)
	check(MikanIdentityCookieKey, c.mikanClientConfig.IdentityCookie)

	if len(errs) > 0 {
		return fmt.Errorf("validation failed: %s", strings.Join(errs, ", "))
	}

	return nil
}
