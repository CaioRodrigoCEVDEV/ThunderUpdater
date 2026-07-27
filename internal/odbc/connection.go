package odbc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	_ "github.com/alexbrainman/odbc"
)

const (
	pingTimeout  = 10 * time.Second
	queryTimeout = 30 * time.Second
)

type odbcConnection struct {
	db  *sql.DB
	dsn string
	log *slog.Logger
}

func newConnection(log *slog.Logger) *odbcConnection {
	if log == nil {
		log = slog.Default()
	}
	return &odbcConnection{log: log}
}

func (c *odbcConnection) open(dsn string) error {
	connStr := fmt.Sprintf("DSN=%s;UID=;PWD=;", dsn)

	db, err := sql.Open("odbc", connStr)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return fmt.Errorf("%w: %w", ErrConnectionFailed, err)
	}

	c.db = db
	c.dsn = dsn
	return nil
}

func (c *odbcConnection) close() {
	if c.db != nil {
		c.db.Close()
		c.db = nil
		c.dsn = ""
	}
}

func (c *odbcConnection) isConnected() bool {
	if c.db == nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	return c.db.PingContext(ctx) == nil
}

func (c *odbcConnection) queryRowScalar(query string) (string, error) {
	if c.db == nil {
		return "", ErrDatabaseUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	var result string
	err := c.db.QueryRowContext(ctx, query).Scan(&result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrColumnNotFound
		}
		return "", fmt.Errorf("%w: %w", ErrInvalidQuery, err)
	}

	return strings.TrimSpace(result), nil
}
