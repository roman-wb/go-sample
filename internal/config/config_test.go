package config

import (
	"os"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/pkg/config"

	"github.com/stretchr/testify/require"
)

func Test_New_Success_Default(t *testing.T) {
	config, err := config.New[Config]()
	require.NoError(t, err)

	// Global config
	require.False(t, config.Config.Debug)
	require.EqualValues(t, 10*time.Second, config.Config.ShutdownTimeout)

	// Local config
	require.EqualValues(t, "App Name", config.Name)
}

func Test_New_Success_Custom(t *testing.T) {
	os.Setenv("APP_DEBUG", "true")
	os.Setenv("APP_NAME", "Custom App Name")

	config, err := config.New[Config]()
	require.NoError(t, err)

	// Global config
	require.True(t, config.Config.Debug)
	require.EqualValues(t, 10*time.Second, config.Config.ShutdownTimeout)

	// Local config
	require.EqualValues(t, "Custom App Name", config.Name)
}

func Test_New_Error(t *testing.T) {
	os.Setenv("APP_DEBUG", "some")

	_, err := config.New[Config]()

	require.EqualError(t, err, `envconfig process: envconfig.Process: assigning APP_APP_DEBUG to Debug: converting 'some' to type bool. details: strconv.ParseBool: parsing "some": invalid syntax`)
}
