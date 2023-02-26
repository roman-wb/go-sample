package helpers

import (
	"bytes"

	"github.com/rs/zerolog"
)

func NewLogger() (*zerolog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := zerolog.New(&buf)
	return &logger, &buf
}
