package flipt

import (
	"errors"
	"fmt"

	"google.golang.org/grpc"
	credsInsecure "google.golang.org/grpc/credentials/insecure"
)

const (
	defaultServer   = "localhost:9000"
	defaultInsecure = false
)

var ErrServerBlank = errors.New("server can't be blank")

type OptionFunc func(*Options) error

type Options struct {
	Server   string
	Insecure bool
}

func DefaultOptions() *Options {
	return &Options{
		Server:   defaultServer,
		Insecure: defaultInsecure,
	}
}

func NewOptions(optionFuncs ...OptionFunc) (*Options, error) {
	options := DefaultOptions()
	for _, optionFunc := range optionFuncs {
		if err := optionFunc(options); err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}
	}

	return options, nil
}

func (o *Options) dialOptions() []grpc.DialOption {
	var dialOptions []grpc.DialOption

	if o.Insecure {
		dialOptions = append(dialOptions, grpc.WithTransportCredentials(
			credsInsecure.NewCredentials(),
		))
	}

	return dialOptions
}

func WithServer(server string) OptionFunc {
	return func(opts *Options) error {
		if server == "" {
			return ErrServerBlank
		}

		opts.Server = server

		return nil
	}
}

func WithInsecure(insecure bool) OptionFunc {
	return func(opts *Options) error {
		opts.Insecure = insecure

		return nil
	}
}
