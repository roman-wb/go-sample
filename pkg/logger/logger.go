package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level zerolog.Level) *zerolog.Logger {
	logger := zerolog.
		New(os.Stdout).
		Level(level).
		With().
		Timestamp().
		Str("podname", os.Getenv("PODNAME")).
		Logger()

	return &logger
}
