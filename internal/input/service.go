package input

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	pb "prime_number_challenge/pkg/prime_number"
)

type Config struct {
	BackgroundServiceAddress string `yaml:"backgroundServiceAddress" env:"BACKGROUND_SERVICE_ADDRESS"`
	Port                     string `yaml:"port" env:"PORT"`
	LoadBalancingModel       string `yaml:"loadBalancingModel" env:"LOAD_BALANCING_MODEL"`
	LogLevel                 string `yaml:"logLevel" env:"LOG_LEVEL"`
}
type Service struct {
	Config *Config
	Logger *log.Logger
	c      pb.PrimeNumberServiceClient
}

func New(config *Config, c pb.PrimeNumberServiceClient) (*Service, error) {
	loglvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s, the correct values are: panic, fatal, error, warn, info, debug, trace", config.LogLevel)
	}

	if config == nil {
		config = DefaultConfig()
	}

	logger := log.New()
	logger.SetLevel(loglvl)
	log.Debugf("Connected to background service at %s", config.BackgroundServiceAddress)

	return &Service{
		Config: config,
		Logger: logger,
		c:      c,
	}, nil
}

func (s *Service) Run() error {
	s.Logger.Info("Starting server")
	http.HandleFunc("/", s.NumHandler)
	http.HandleFunc("/healtz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	portString := fmt.Sprintf(":%s", s.Config.Port)
	s.Logger.Infof("Listening on port %s", portString)
	err := http.ListenAndServe(portString, nil)
	if err != nil {
		return err
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Port:                     "8080",
		BackgroundServiceAddress: "localhost:50051",
	}
}
