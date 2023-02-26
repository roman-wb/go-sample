package flipt

import (
	"fmt"

	"go.flipt.io/flipt-grpc"
	"google.golang.org/grpc"
)

type CleanUpFunc func() error

func New(optionFuncs ...OptionFunc) (flipt.FliptClient, CleanUpFunc, error) {
	options, err := NewOptions(optionFuncs...)
	if err != nil {
		return nil, nil, fmt.Errorf("new options: %w", err)
	}

	conn, err := grpc.Dial(options.Server, options.dialOptions()...)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc dial: %w", err)
	}

	client := flipt.NewFliptClient(conn)

	cleanUpFunc := func() error {
		return conn.Close()
	}

	return client, cleanUpFunc, nil
}
