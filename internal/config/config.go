package config

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/torrent"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port string

	bangumiClientID     string
	bangumiClientSecret string

	qbittorrentConfig torrent.QbittorrentClientConfig

	mikanAuthenticationCookie string
}

var cfg config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file:%s", err)
	}

	cfg.port = os.Getenv(PortKey)
	cfg.bangumiClientID = os.Getenv(BangumiClientIDKey)
	cfg.bangumiClientSecret = os.Getenv(BangumiClientSecretKey)

	cfg.qbittorrentConfig = torrent.QbittorrentClientConfig{
		Port:     os.Getenv(QBittorrentPortKey),
		Username: os.Getenv(QBittorrentUserNameKey),
		Password: os.Getenv(QBittorrentPasswordKey),
	}

	cfg.mikanAuthenticationCookie = os.Getenv(MikanAuthenticationCookieKey)

	if err := cfg.validate(); err != nil {
		log.Fatalf("config is invalid:%s", err)
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

func QBittorrentConfig() torrent.QbittorrentClientConfig {
	return cfg.qbittorrentConfig
}

func MikanAuthenticationCookie() string {
	return cfg.mikanAuthenticationCookie
}

func LocalServerAddress() string {
	return fmt.Sprintf("http://localhost:%s", cfg.port)
}

func (c config) validate() error {
	if c.port == "" {
		return fmt.Errorf("%s is empty", PortKey)
	}

	if c.bangumiClientID == "" {
		return fmt.Errorf("%s is empty", BangumiClientIDKey)
	}

	if c.bangumiClientSecret == "" {
		return fmt.Errorf("%s is empty", BangumiClientSecretKey)
	}

	if c.qbittorrentConfig.Port == "" {
		return fmt.Errorf("%s is empty", QBittorrentPortKey)
	}

	if c.qbittorrentConfig.Username == "" {
		return fmt.Errorf("%s is empty", QBittorrentUserNameKey)
	}

	if c.qbittorrentConfig.Password == "" {
		return fmt.Errorf("%s is empty", QBittorrentPasswordKey)
	}

	if c.mikanAuthenticationCookie == "" {
		return fmt.Errorf("%s is empty", MikanAuthenticationCookieKey)
	}

	return nil
}
