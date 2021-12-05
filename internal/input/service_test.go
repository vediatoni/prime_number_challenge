package input

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io/ioutil"
	"net"
	"net/http/httptest"
	pb "prime_number_challenge/pkg/prime_number"
	"strings"
	"testing"
)

type service struct {
	pb.UnimplementedPrimeNumberServiceServer
}

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterPrimeNumberServiceServer(s, &service{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func createTestSvc(cfg *Config) (*Service, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPrimeNumberServiceClient(conn)
	svc, err := New(cfg, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	return svc, nil
}

func TestNew(t *testing.T) {
	svc, err := createTestSvc(&Config{Port: "8080"})
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	if svc == nil {
		t.Fatal("RealService is nil")
	}
}

func TestNewNilConfig(t *testing.T) {
	svc, err := createTestSvc(nil)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	t.Logf("%v, %v", *svc.Config, *DefaultConfig())
	if *svc.Config != *DefaultConfig() {
		t.Fatalf("Config is not default: %v, it should be %v", svc.Config, DefaultConfig())
	}
}

func TestService_Run(t *testing.T) {
	svc, err := createTestSvc(nil)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	if err := svc.Run(); err != nil {
		t.Fatalf("failed to run service: %v", err)
	}

	// TODO: send http request to the running service (health check)
	httptest.NewRequest("GET", "/healtz", nil)


}

func TestRoutes(t *testing.T) {
}

func TestService_NumHandler(t *testing.T) {
	testCases := []struct {
		name    string
		request interface{}
		expect  string
	}{
		{
			name: "Test 1",
			request: &Payload{
				Number: 3,
			},
			expect: FailedToConnectToBackgroundService,
		},
		{
			name: "Test 2",
			request: &Payload{
				Number: -20,
			},
			expect: LessThanZero,
		},
		{
			name:    "Test 3",
			request: "something random",
			expect:  InvalidInput,
		},
		{
			name:    "Test 4",
			request: &struct{ rnd string }{rnd: "something random"},
			expect:  LessThanZero,
		},
	}

	svc, err := createTestSvc(nil)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(testCase.request)
			if err != nil {
				t.Fatalf("failed to marshal request: %v", err)
			}
			body := bytes.NewReader(jsonBody)
			r := httptest.NewRequest("POST", "/", body)
			w := httptest.NewRecorder()
			svc.NumHandler(w, r)

			resp := w.Result()
			defer resp.Body.Close()
			out, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			result := strings.TrimRight(string(out), "\n")
			if result != testCase.expect {
				t.Fatalf("expected %v, got %v", result, testCase.expect)
			}
		})
	}
}
