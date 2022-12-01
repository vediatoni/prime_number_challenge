package database

import (
	"fmt"
	"testing"

	"github.com/vediatoni/prime_number_challenge/pkg/prime_number"
)

func TestNew(t *testing.T) {
	// TODO
}

func TestService_Insert(t *testing.T) {
	// TODO
}

// TestBuildInsertQuery - Tests the function that builds INSERT query from structs
func TestBuildInsertQuery(t *testing.T) {
	const tableName = "test"
	type Test struct {
		ID      int    `json:"id,omitempty"`
		Name    string `json:"name,omitempty"`
		Surname string `json:"surname,omitempty"`
		Alive   bool   `json:"alive,omitempty"`
	}

	testCases := []struct {
		name     string
		testCase interface{}
		want     string
	}{
		{
			name: "Test case 1",
			testCase: &Test{
				ID:      1,
				Name:    "John",
				Surname: "Doe",
				Alive:   true,
			},
			want: fmt.Sprintf("INSERT INTO %s (id,name,surname,alive) VALUES (1,'John','Doe',TRUE)", tableName),
		},
		{
			name: "Test case 2",
			testCase: &Test{
				ID:    2,
				Name:  "Jane",
				Alive: false,
			},
			want: fmt.Sprintf("INSERT INTO %s (id,name,alive) VALUES (2,'Jane',FALSE)", tableName),
		},
		{
			name: "Test case 3",
			testCase: &prime_number.IsPrimeNumberResponse{
				NumberTested:                 2,
				IsPrime:                      true,
				ValidationTime:               2,
				TimeNeededToValidateMicrosec: 500,
			},
			want: fmt.Sprintf("INSERT INTO %s (number_tested,is_prime,validation_time,"+
				"time_needed_to_validate_microsec) VALUES (2,TRUE,2,500)", tableName),
		},
		{
			name: "Test case 4",
			testCase: &prime_number.IsPrimeNumberResponse{
				NumberTested: 2,
				IsPrime:      true,
			},
			want: fmt.Sprintf("INSERT INTO %s (number_tested,is_prime,validation_time,"+
				"time_needed_to_validate_microsec) VALUES (2,TRUE,0,0)", tableName),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got := BuildInsertQuery(tableName, testCase.testCase)
			if got != testCase.want {
				t.Errorf("BuildInsertQuery() = %v, want %v", got, testCase.want)
			}
		})
	}
}
