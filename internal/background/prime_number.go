package background

import (
	"context"
	"fmt"
	"prime_number_challenge/pkg/database"
	pb "prime_number_challenge/pkg/prime_number"
	"time"
)

func (s *Service) IsPrimeNumber(_ context.Context, in *pb.IsPrimeNumberRequest) (*pb.IsPrimeNumberResponse, error) {
	start := time.Now()
	result := isPrime(in.GetNumber())
	elapsed := time.Since(start)
	elapsedMs := int32(elapsed / time.Microsecond)
	s.Logger.Debugf("Prime check took %v ms", elapsedMs)

	res := pb.IsPrimeNumberResponse{
		NumberTested:                 in.GetNumber(),
		IsPrime:                      result,
		ValidationTime:               time.Now().Unix(), // works until 2038 :) also, 4 bytes
		TimeNeededToValidateMicrosec: elapsedMs,
	}

	s.SaveToDB(&res)

	return &res, nil
}

// isPrime checks if a number is prime with brute force
func isPrime(n int32) bool {
	if n <= 1 {
		return false
	}
	for i := int32(2); i <= n-1; i++ {
		if (n % i) == 0 {
			fmt.Println(i)
			return false
		}
	}
	return true
}

func (s *Service) SaveToDB(number *pb.IsPrimeNumberResponse) {
	query := database.BuildInsertQuery(number)
	s.Logger.Debugf("Query: %s", query)
	insert, err := s.db.Insert(query)
	if err != nil {
		fmt.Println(insert)
		s.Logger.Errorf("Error inserting into DB: %v", err)
	}
}