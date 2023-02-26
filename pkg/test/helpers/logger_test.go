package helpers

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func Test_NewLogger(t *testing.T) {
	logger, buf := NewLogger()

	logger.Trace().Msg("same log 1")
	logger.Info().Msg("same log 2")

	require.EqualValues(t, logger.GetLevel(), zerolog.TraceLevel)

	require.Contains(t, buf.String(), "same log 1")
	require.Contains(t, buf.String(), "same log 2")
}
