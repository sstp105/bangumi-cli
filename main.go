/*
Copyright © 2025 Yang Li
*/
package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sstp105/bangumi-cli/cmd"
	"github.com/sstp105/bangumi-cli/internal/config"
)

var log *logrus.Logger

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file:%s", err)
	}

	config.Load()

	cmd.Execute()
}
