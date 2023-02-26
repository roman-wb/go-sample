package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const prefix = "APP"

type Config struct {
	Debug           bool          `envconfig:"APP_DEBUG" default:"false"`
	ShutdownTimeout time.Duration `envconfig:"APP_SHUTDOWN_TIMEOUT" default:"10s"`

	PgURL    string `envconfig:"APP_PG_URL"`
	RedisURL string `envconfig:"APP_REDIS_URL"`

	FliptServer   string        `envconfig:"APP_FLIPT_SERVER" default:"localhost:9000"`
	FliptInsecure bool          `envconfig:"APP_FLIPT_INSECURE" default:"false"`
	FliptTimeout  time.Duration `envconfig:"APP_FLIPT_TIMEOUT" default:"5s"`

	Server       string        `envconfig:"APP_SERVER" default:"localhost:3000"`
	ReadTimeout  time.Duration `envconfig:"APP_READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"APP_WRITE_TIMEOUT" default:"10s"`
}

func New[C any]() (*C, error) {
	var cfg C
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("envconfig process: %w", err)
	}

	return &cfg, nil
}
