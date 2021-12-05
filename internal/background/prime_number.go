package background

import (
	"context"
	"prime_number_challenge/pkg/database"
	pb "prime_number_challenge/pkg/prime_number"
	"time"
)

func (s *Service) IsPrimeNumber(_ context.Context, in *pb.IsPrimeNumberRequest) (*pb.IsPrimeNumberResponse, error) {
	s.Logger.Debugf("IsPrimeNumber: starting calculation for number: %d", in.Number)
	start := time.Now()
	result := isPrime(in.GetNumber())
	elapsed := time.Since(start)
	s.Logger.Debugf("Prime check took %v micro seconds", elapsed.Microseconds())

	res := pb.IsPrimeNumberResponse{
		NumberTested:                 in.GetNumber(),
		IsPrime:                      result,
		ValidationTime:               time.Now().Unix(), // works until 2038 :) also, 4 bytes
		TimeNeededToValidateMicrosec: elapsed.Microseconds(),
	}

	go s.SaveToDB(&res)

	return &res, nil
}

// isPrime checks if a number is prime with brute force
func isPrime(n int32) bool {
	if n <= 1 {
		return false
	}
	for i := int32(2); i <= n-1; i++ {
		if (n % i) == 0 {
			return false
		}
	}
	return true
}

func (s *Service) SaveToDB(number *pb.IsPrimeNumberResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := database.BuildInsertQuery("prime_numbers",number)
	s.Logger.Debugf("Query: %s", query)
	insert, err := s.db.Insert(ctx, query)
	if err != nil {
		s.Logger.Errorf("Error inserting into DB: %v", err)
	}

	s.Logger.Debugf("Insert results: %v", insert)
}
