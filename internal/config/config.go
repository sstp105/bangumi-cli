package config

import (
	"fmt"
	"github.com/sstp105/bangumi-cli/internal/mikan"
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

	mikanClientConfig mikan.ClientConfig
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
		Server:   os.Getenv(QBittorrentServerKey),
		Username: os.Getenv(QBittorrentUserNameKey),
		Password: os.Getenv(QBittorrentPasswordKey),
	}

	cfg.mikanClientConfig = mikan.ClientConfig{
		IdentityCookie: os.Getenv(MikanIdentityCookieKey),
	}

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

func MikanClientConfig() mikan.ClientConfig {
	return cfg.mikanClientConfig
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

	if c.qbittorrentConfig.Server == "" {
		return fmt.Errorf("%s is empty", QBittorrentServerKey)
	}

	if c.qbittorrentConfig.Username == "" {
		return fmt.Errorf("%s is empty", QBittorrentUserNameKey)
	}

	if c.qbittorrentConfig.Password == "" {
		return fmt.Errorf("%s is empty", QBittorrentPasswordKey)
	}

	if c.mikanClientConfig.IdentityCookie == "" {
		return fmt.Errorf("%s is empty", MikanIdentityCookieKey)
	}

	return nil
}
