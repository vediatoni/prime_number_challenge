package main

import (
	log "github.com/sirupsen/logrus"
	"prime_number_challenge/internal/input"
	"prime_number_challenge/pkg/config"
)

const configFilePath = "config/dev.input.yaml"

func main() {
	inf, err := config.LoadConfigFromFile(configFilePath, input.Config{})
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	cfg := inf.(input.Config)

	s, err := input.New(&cfg)
	if err != nil {
		log.Fatalf("Error creating input service: %v", err)
	}

	err = s.Run()
	if err != nil{
		log.Fatalf("Error running input service: %v", err)
	}
}
