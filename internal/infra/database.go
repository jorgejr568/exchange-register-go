package infra

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

var (
	ErrNotFound = fmt.Errorf("row not found: %w", sql.ErrNoRows)
)

type DB interface {
	// Exec executes a query without returning any rows.
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// Query executes a query that returns rows, typically a SELECT.
	Query(ctx context.Context, target interface{}, query string, args ...interface{}) error

	// QueryOne executes a query that returns one row, typically a SELECT.
	QueryOne(ctx context.Context, target interface{}, query string, args ...interface{}) error

	// Close closes the database, releasing any open resources.
	Close() error
}

type ksqlPgDB struct{ db *ksql.DB }

func (k ksqlPgDB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return k.db.Exec(ctx, query, args...)
}

func (k ksqlPgDB) Query(ctx context.Context, target interface{}, query string, args ...interface{}) error {
	return k.db.Query(ctx, target, query, args...)
}

func (k ksqlPgDB) QueryOne(ctx context.Context, target interface{}, query string, args ...interface{}) error {
	err := k.db.QueryOne(ctx, target, query, args...)
	if err != nil {
		if errors.Is(err, ksql.ErrRecordNotFound) {
			return ErrNotFound
		}

		return err
	}

	return nil
}

func (k ksqlPgDB) Close() error {
	return k.db.Close()
}

func NewKsqlPgDB(ctx context.Context, connection string) (DB, error) {
	db, err := kpgx.New(ctx, connection, ksql.Config{
		MaxOpenConns: 10,
		TLSConfig:    nil,
	})
	if err != nil {
		return nil, err
	}

	return &ksqlPgDB{
		db: &db,
	}, nil
}
