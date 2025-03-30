package config

import (
	"fmt"
	"os"
)

type config struct {
	port                string
	bangumiClientID     string
	bangumiClientSecret string
}

var cfg config

func Load() {
	cfg.port = os.Getenv("LOCAL_SERVER_PORT")
	cfg.bangumiClientID = os.Getenv("BANGUMI_CLIENT_ID")
	cfg.bangumiClientSecret = os.Getenv("BANGUMI_CLIENT_SECRET")
}

func BangumiClientID() string {
	return cfg.bangumiClientID
}

func BangumiClientSecret() string {
	return cfg.bangumiClientSecret
}

func Port() string {
	return cfg.port
}

func LocalServerAddress() string {
	return fmt.Sprintf("http://localhost:%s", cfg.port)
}