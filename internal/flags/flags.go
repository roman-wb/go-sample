package flags

import (
	"context"
	"fmt"
	"time"

	"go.flipt.io/flipt-grpc"
)

//go:generate mockery --name FliptClient --srcpkg go.flipt.io/flipt-grpc --case underscore

type Flags struct {
	client  flipt.FliptClient
	timeout time.Duration
}

func New(client flipt.FliptClient, timeout time.Duration) *Flags {
	return &Flags{
		client:  client,
		timeout: timeout,
	}
}

func (f *Flags) IsSampleFeatureEnabled(ctx context.Context) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	flag, err := f.client.GetFlag(ctx, &flipt.GetFlagRequest{
		Key: SampleFeatureKey,
	})
	if err != nil {
		return false, fmt.Errorf("get flag: %w", err)
	}

	return flag.Enabled, nil
}
