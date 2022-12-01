package input

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
	pb "github.com/vediatoni/prime_number_challenge/pkg/prime_number"
)

type Config struct {
	BackgroundServiceAddress string `yaml:"backgroundServiceAddress" env:"BACKGROUND_SERVICE_ADDRESS"`
	Port                     string `yaml:"port" env:"PORT"`
	LoadBalancingModel       string `yaml:"loadBalancingModel" env:"LOAD_BALANCING_MODEL"`
	LogLevel                 string `yaml:"logLevel" env:"LOG_LEVEL"`
}
type Service struct {
	Config     *Config
	Logger     *log.Logger
	c          pb.PrimeNumberServiceClient
	httpServer *http.Server
}

func New(config *Config, c pb.PrimeNumberServiceClient) (*Service, error) {
	if config == nil {
		config = DefaultConfig()
	}

	loglvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s, the correct values are: panic, fatal, error, warn, info, debug, trace", config.LogLevel)
	}

	logger := log.New()
	logger.SetLevel(loglvl)

	return &Service{
		Config: config,
		Logger: logger,
		c:      c,
		httpServer: &http.Server{
			Addr: fmt.Sprintf(":%s", config.Port),
		},
	}, nil
}

func (s *Service) handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.NumHandler)
	mux.HandleFunc("/healtz", healthCheck)

	return mux
}

func (s *Service) Run() error {
	s.Logger.Info("Starting server")
	s.httpServer.Handler = s.handler()
	portString := fmt.Sprintf(":%s", s.Config.Port)
	s.Logger.Infof("Listening on port %s", portString)

	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Port:                     "8080",
		BackgroundServiceAddress: "localhost:50051",
		LoadBalancingModel:       "round_robin",
		LogLevel:                 "info",
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
