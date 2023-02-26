package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/pkg/env"
	"github.com/roman-wb/go-sample/pkg/test/helpers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func Test_New_Success(t *testing.T) {
	err := env.Init("../../.env")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, _ := helpers.NewLogger()

	postgresClient, cleanUp, err := New(ctx, logger, os.Getenv("APP_PG_URL"))
	defer cleanUp()

	// logger

	require.NoError(t, err)
	require.IsType(t, &pgxpool.Pool{}, postgresClient)
}

func Test_New_ErrorParse(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	err := env.Init("../../.env")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, _ := helpers.NewLogger()

	postgresClient, cleanUp, err := New(ctx, logger, "wrong dsn")

	require.EqualError(t, err, "pgxpool parse config: cannot parse `wrong dsn`: failed to parse as DSN (invalid dsn)")
	require.Nil(t, postgresClient)
	require.Nil(t, cleanUp)
}

func Test_New_ErrorPing(t *testing.T) {
	t.Skip("skipping test")

	err := env.Init("../../.env")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger, _ := helpers.NewLogger()

	postgresClient, cleanUp, err := New(ctx, logger, "postgres://test@localhost:9999")

	require.EqualError(t, err, "conn ping: failed to connect to `host=localhost user=test database=`: dial error (dial tcp localhost:9999: connect: no route to host)")
	require.Nil(t, postgresClient)
	require.Nil(t, cleanUp)
}
