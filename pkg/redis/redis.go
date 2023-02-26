package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CleanUpFunc func() error

func New(ctx context.Context, url string) (*redis.Client, CleanUpFunc, error) {
	options, err := redis.ParseURL(url)
	if err != nil {
		return nil, nil, fmt.Errorf("redis parse url: %w", err)
	}

	conn := redis.NewClient(options)
	if _, err := conn.Ping(ctx).Result(); err != nil {
		return nil, nil, fmt.Errorf("conn ping: %w", err)
	}

	cleanUpFunc := func() error {
		return conn.Close()
	}

	return conn, cleanUpFunc, nil
}
