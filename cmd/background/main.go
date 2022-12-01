package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/vediatoni/prime_number_challenge/internal/background"
	"github.com/vediatoni/prime_number_challenge/pkg/config"
)

const configFilePath = "config/dev.background.yaml"

func main() {
	inf, err := config.LoadConfigFromFile(configFilePath, background.Config{})
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	cfg := inf.(background.Config)
	background, err := background.New(&cfg)
	if err != nil {
		log.Fatalf("Error creating background service: %v", err)
	}
	background.Logger.Debugf("Log level: %s", background.Logger.Level)

	err = background.Run()
	if err != nil {
		log.Fatalf("Error running background service: %v", err)
	}
}
