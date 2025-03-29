package config

import "os"

type Config struct {
	Port                string
	BangumiClientID     string
	BangumiClientSecret string
}

var AppConfig Config

func Load() {
	AppConfig.Port = os.Getenv("LOCAL_SERVER_PORT")
	AppConfig.BangumiClientID = os.Getenv("BANGUMI_CLIENT_ID")
	AppConfig.BangumiClientSecret = os.Getenv("BANGUMI_CLIENT_SECRET")
}