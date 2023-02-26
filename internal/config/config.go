package config

import "github.com/roman-wb/go-sample/pkg/config"

type Config struct {
	*config.Config

	Name string `envconfig:"APP_NAME" default:"App Name"`
}
