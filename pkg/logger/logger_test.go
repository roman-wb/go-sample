package logger

import (
	"bytes"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	os.Setenv("PODNAME", "same pod")

	logger := New(zerolog.InfoLevel)

	var buf bytes.Buffer
	testLogger := logger.Output(&buf)

	testLogger.Debug().Msg("same log 1")
	testLogger.Info().Msg("same log 2")

	require.EqualValues(t, logger.GetLevel(), zerolog.InfoLevel)
	require.EqualValues(t, testLogger.GetLevel(), zerolog.InfoLevel)

	require.NotContains(t, buf.String(), "same log 1")
	require.Contains(t, buf.String(), "same log 2")
	require.Contains(t, buf.String(), `"podname":"same pod"`)
}
