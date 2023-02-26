package redis

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/pkg/env"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func Test_New_Success(t *testing.T) {
	err := env.Init("../../.env")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient, cleanUp, err := New(ctx, os.Getenv("APP_REDIS_URL"))
	defer cleanUp()

	require.NoError(t, err)
	require.IsType(t, &redis.Client{}, redisClient)
}

func Test_New_ErrorParse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient, cleanUp, err := New(ctx, "...")

	require.EqualError(t, err, "redis parse url: redis: invalid URL scheme: ")
	require.Nil(t, redisClient)
	require.Nil(t, cleanUp)
}

func Test_New_ErrorPing(t *testing.T) {
	t.Skip("skipping test")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	redisClient, cleanUp, err := New(ctx, "redis://localhost:1111")

	require.EqualError(t, err, "conn ping: dial tcp [::1]:1111: connect: no route to host")
	require.Nil(t, redisClient)
	require.Nil(t, cleanUp)
}
