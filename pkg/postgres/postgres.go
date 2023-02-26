package postgres

import (
	"context"
	"fmt"

	adapter "github.com/jackc/pgx-zerolog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"
)

type CleanUpFunc func()

func New(ctx context.Context, logger *zerolog.Logger, url string) (*pgxpool.Pool, CleanUpFunc, error) {
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool parse config: %w", err)
	}

	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   adapter.NewLogger(*logger, adapter.WithSubDictionary("pgx")),
		LogLevel: tracelog.LogLevelInfo,
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, nil, fmt.Errorf("pgxpool new with config: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, nil, fmt.Errorf("conn ping: %w", err)
	}

	cleanUpFunc := func() {
		conn.Close()
	}

	return conn, cleanUpFunc, nil
}
