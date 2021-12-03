package input

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
	pb "prime_number_challenge/pkg/prime_number"
)

type Config struct {
	BackgroundServiceAddress string `yaml:"backgroundServiceAddress" env:"BACKGROUND_SERVICE_ADDRESS" required:"true"`
	Port                     string `yaml:"port" env:"PORT" required:"true"`
}

type Service struct {
	Config *Config
	Logger *log.Logger
	conn   *grpc.ClientConn
	c      pb.PrimeNumberServiceClient
}

func New(config *Config) (*Service, error) {
	logger := log.New()
	conn, err := grpc.Dial(config.BackgroundServiceAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	log.Debugf("Connected to background service at %s", config.BackgroundServiceAddress)

	return &Service{
		Config: config,
		Logger: logger,
		conn:   conn,
		c:      pb.NewPrimeNumberServiceClient(conn),
	}, nil
}

func (s *Service) Run() error {
	s.Logger.Info("Starting server")

	http.HandleFunc("/", s.NumHandler)
	http.HandleFunc("/healtz", func (w http.ResponseWriter, r *http.Request) {
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