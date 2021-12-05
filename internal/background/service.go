package background

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"prime_number_challenge/pkg/database"
	pb "prime_number_challenge/pkg/prime_number"
)

type Config struct {
	DatabaseConnectionString string `yaml:"databaseConnectionString" env:"DATABASE_CONNECTION_STRING"`
	SelfSvcAddress           string `yaml:"selfSvcAddress" env:"SELF_SVC_ADDRESS"`
	LogLevel                 string `yaml:"logLevel" env:"LOG_LEVEL"`
}

type Service struct {
	pb.UnimplementedPrimeNumberServiceServer
	Config *Config
	Logger *log.Logger
	db     *database.Service
}

func New(config *Config) (*Service, error) {
	logger := log.New()
	loglvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %s, the correct values are: panic, fatal, error, warn, info, debug, trace", config.LogLevel)
	}
	logger.SetLevel(loglvl)
	db, err := database.New(config.DatabaseConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}
	log.Debugf("Database service created and connected")

	return &Service{
		Config: config,
		Logger: logger,
		db:     db,
	}, nil
}

func (s *Service) Run() error {
	s.Logger.Info("Starting server")
	lis, err := net.Listen("tcp", s.Config.SelfSvcAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPrimeNumberServiceServer(grpcServer, s)
	s.Logger.Info("Listening on " + s.Config.SelfSvcAddress)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
