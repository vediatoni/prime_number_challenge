package input

import (
	"encoding/json"
	"fmt"
	"net/http"
	pb "prime_number_challenge/pkg/prime_number"
)

type Payload struct {
	Number int32 `json:"n"`
}

func (s *Service) NumHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload Payload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		s.Logger.Errorf("NumHandler: error decoding payload: %v", err)
		http.Error(w, "Oops, something went wrong :(", http.StatusBadRequest)
		return
	}
	s.Logger.Debugf("NumHandler: payload: %v", payload)

	var response *pb.IsPrimeNumberResponse
	response, err = s.c.IsPrimeNumber(r.Context(), &pb.IsPrimeNumberRequest{Number: payload.Number})
	if err != nil {
		s.Logger.Errorf("NumHandler: could not check if the number is prime: %v", err)
		http.Error(w, "Oops, something went wrong :(", http.StatusBadRequest)
		return
	}
	s.Logger.Debugf("NumHandler: response: %v", response.IsPrime)


	_, err = w.Write([]byte(fmt.Sprintf("%v", response.IsPrime)))
	if err != nil {
		s.Logger.Errorf("NumHandler: error writing response: %v", err)
		http.Error(w, "Oops, something went wrong :(", http.StatusBadRequest)
		return
	}
}
