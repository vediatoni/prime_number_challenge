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
	DatabaseConnectionString string `yaml:"databaseConnectionString" env:"DATABASE_CONNECTION_STRING" required:"true"`
	Port                     string `yaml:"port" env:"PORT" required:"true"`
}

type Service struct {
	pb.UnimplementedPrimeNumberServiceServer
	Config *Config
	Logger *log.Logger
	db     *database.Service
}

func New(config *Config) (*Service, error) {
	db, err := database.New(config.DatabaseConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to create database service: %w", err)
	}
	return &Service{
		Config: config,
		Logger: log.New(),
		db:     db,
	}, nil
}

func (s *Service) Run() error {
	s.Logger.Info("Starting server")
	port := fmt.Sprintf(":%s", s.Config.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPrimeNumberServiceServer(grpcServer, s)
	s.Logger.Info("Listening on port " + s.Config.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
