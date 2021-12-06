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
	"net/http"
	"net/http/httptest"
	"prime_number_challenge/internal/background"
	pb "prime_number_challenge/pkg/prime_number"
	"strings"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func initBackgroundSvc() {
	bg, err := background.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterPrimeNumberServiceServer(s, bg)
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
	initBackgroundSvc()

	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial bufnet: %v", err)
	}
	svc, err := New(cfg, pb.NewPrimeNumberServiceClient(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	return svc, nil
}

func TestNew(t *testing.T) {
	svc, err := createTestSvc(&Config{
		Port:     "8080",
		LogLevel: "debug",
	})
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}
	if svc == nil {
		t.Fatal("service is nil")
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

	serviceRunning := make(chan struct{})
	serviceDone := make(chan struct{})
	go func() {
		close(serviceRunning)
		if err := svc.Run(); err != http.ErrServerClosed && err != nil {
			t.Fatalf("failed to run service: %v", err)
		}
		defer close(serviceDone)
	}()

	<-serviceRunning

	resp, err := http.Get(fmt.Sprintf("http://localhost%v/healtz", svc.httpServer.Addr))
	if err != nil {
		t.Fatalf("failed to get healthz: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("healthz status code is not OK: %v", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = svc.httpServer.Shutdown(ctx)
	if err != nil {
		t.Fatalf("failed to shutdown server: %v", err)
	}

	<-serviceDone
}

func TestService_Healtz(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(healthCheck))
	defer svr.Close()
	resp, err := http.Get(svr.URL)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("response status code is not OK: %v", resp.StatusCode)
	}
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
			expect: "true",
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
		{
			name:    "Test 5",
			request: &Payload{Number: 1},
			expect:  "false",
		},
		{
			name:    "Test 6",
			request: &Payload{Number: 70},
			expect:  "false",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			svc, err := createTestSvc(nil)
			if err != nil {
				t.Fatalf("failed to create service: %v", err)
			}
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
				t.Fatalf("expected %v, got %v", testCase.expect, result)
			}
		})
	}
}
