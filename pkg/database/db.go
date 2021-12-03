package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	Logger *log.Logger
	conn   *sql.DB
}

// New - Create a new database connection
func New(dbConnectionString string) (*Service, error) {
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Service{
		Logger: log.New(),
		conn:   db,
	}, nil
}

// Insert - Insert something via query
func (s *Service) Insert(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	s.Logger.Debugf("Insert: %s", query)
	return s.conn.ExecContext(ctx, query, args...)
}

// docker run --name postgres -e POSTGRES_PASSWORD=Test12345 -e POSTGRES_DB=prime_numbers -p 5432:5432 -d postgres
// migrate -database postgres://postgres:Test12345@localhost:5432/prime_numbers?sslmode=disable -source file://migrations up
