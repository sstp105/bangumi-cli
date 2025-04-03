package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port                string
	bangumiClientID     string
	bangumiClientSecret string
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

	if err := cfg.validate(); err != nil {
		log.Fatalf("config is invalid:%s", err)
	}
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

	return nil
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
