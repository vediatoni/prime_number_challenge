package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"prime_number_challenge/internal/input"
	"prime_number_challenge/pkg/config"
	pb "prime_number_challenge/pkg/prime_number"
)

const configFilePath = "config/dev.input.yaml"

func main() {
	inf, err := config.LoadConfigFromFile(configFilePath, input.Config{})
	if err != nil {
		log.Fatalf("Error loading config file: %v", err)
	}
	cfg := inf.(input.Config)

	// Could wrap this in a function, and maybe provide more options, but this is fine for now and for this task.
	var dialOpts []grpc.DialOption
	dialOpts = append(dialOpts, grpc.WithInsecure())
	if cfg.LoadBalancingModel != "" {
		log.Info("GRPC load balancing on background service enabled")
		tmp := fmt.Sprintf(`{"loadBalancingPolicy":"%v"}`, cfg.LoadBalancingModel)
		dialOpts = append(dialOpts, grpc.WithDefaultServiceConfig(tmp))
	}

	grpcConn, err := grpc.Dial(cfg.BackgroundServiceAddress, dialOpts...)
	if err != nil {
		log.Fatalf("Error connecting to background service: %v", err)
	}

	s, err := input.New(&cfg, pb.NewPrimeNumberServiceClient(grpcConn))
	if err != nil {
		log.Fatalf("Error creating input service: %v", err)
	}

	err = s.Run()
	if err != nil {
		log.Fatalf("Error running input service: %v", err)
	}
}
