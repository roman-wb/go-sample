package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_New_Success_Default(t *testing.T) {
	config, err := New[Config]()

	require.NoError(t, err)

	require.False(t, config.Debug)
	require.EqualValues(t, 10*time.Second, config.ShutdownTimeout)

	require.Empty(t, config.PgURL)
	require.Empty(t, config.RedisURL)

	require.EqualValues(t, "localhost:9000", config.FliptServer)
	require.False(t, config.FliptInsecure)
	require.EqualValues(t, 5*time.Second, config.FliptTimeout)

	require.EqualValues(t, "localhost:3000", config.Server)
	require.EqualValues(t, 10*time.Second, config.ReadTimeout)
	require.EqualValues(t, 10*time.Second, config.WriteTimeout)
}

func Test_New_Success_Custom(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_SHUTDOWN_TIMEOUT", "20s")
	os.Setenv("APP_PG_URL", "postgres://localhost:1111/db")
	os.Setenv("APP_REDIS_URL", "redis://localhost:2222/0")
	os.Setenv("APP_FLIPT_SERVER", "localhost:3333")
	os.Setenv("APP_FLIPT_INSECURE", "true")
	os.Setenv("APP_FLIPT_TIMEOUT", "30s")
	os.Setenv("APP_SERVER", "localhost:4444")
	os.Setenv("APP_READ_TIMEOUT", "40s")
	os.Setenv("APP_WRITE_TIMEOUT", "50s")

	config, err := New[Config]()

	require.NoError(t, err)

	require.True(t, config.Debug)
	require.EqualValues(t, 20*time.Second, config.ShutdownTimeout)

	require.EqualValues(t, "postgres://localhost:1111/db", config.PgURL)
	require.EqualValues(t, "redis://localhost:2222/0", config.RedisURL)

	require.EqualValues(t, "localhost:3333", config.FliptServer)
	require.True(t, config.FliptInsecure)
	require.EqualValues(t, 30*time.Second, config.FliptTimeout)

	require.EqualValues(t, "localhost:4444", config.Server)
	require.EqualValues(t, 40*time.Second, config.ReadTimeout)
	require.EqualValues(t, 50*time.Second, config.WriteTimeout)
}

func Test_New_Error(t *testing.T) {
	os.Setenv("APP_DEBUG", "some")

	_, err := New[Config]()

	require.EqualError(t, err, `envconfig process: envconfig.Process: assigning APP_APP_DEBUG to Debug: converting 'some' to type bool. details: strconv.ParseBool: parsing "some": invalid syntax`)
}
