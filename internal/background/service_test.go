package background

import "testing"

func TestNew(t *testing.T) {
	svc, err := New(nil)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
	}
}

func TestService_Run(t *testing.T) {
	// TODO
}

func TestService_IsPrimeNumber(t *testing.T) {
	// TODO
}

func TestService_SaveToDB(t *testing.T) {
	// TODO
}

