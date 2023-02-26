package helpers

import (
	"context"
	"time"
)

func FuncWithTimeout(timeout time.Duration, testFunc func()) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		testFunc()
		cancel()
	}()

	<-ctx.Done()

	switch ctx.Err() {
	case context.Canceled:
		return nil
	default:
		return ctx.Err()
	}
}
