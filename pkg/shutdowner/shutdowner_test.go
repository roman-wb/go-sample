package shutdowner

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/pkg/test/helpers"

	"github.com/stretchr/testify/require"
)

func Test_NewShutdowner(t *testing.T) {
	sd := NewShutdowner(nil)
	require.IsType(t, &Shutdowner{}, sd)
}

func Test_Success_EmptyProcess(t *testing.T) {
	require.NoError(t, helpers.FuncWithTimeout(5*time.Second, func() {
		var wantAfterErr error
		var wantAfterFuncCalled bool

		sd := NewShutdowner(
			func(name string, err error) {
				wantAfterFuncCalled = true
				wantAfterErr = err
			},
		)
		sd.Wait(context.Background(), time.Second)

		require.NoError(t, wantAfterErr)
		require.False(t, wantAfterFuncCalled)
	}))
}

func Test_Success(t *testing.T) {
	require.NoError(t, helpers.FuncWithTimeout(5*time.Second, func() {
		var wantAfterErr error
		var wantAfterFuncCalled bool

		sd := NewShutdowner(
			func(name string, err error) {
				wantAfterFuncCalled = true
				wantAfterErr = err
			},
			WithProcess("test", func(ctx context.Context) error {
				time.Sleep(time.Second)
				return nil
			}),
		)
		sd.Wait(context.Background(), 2*time.Second)

		require.NoError(t, wantAfterErr)
		require.True(t, wantAfterFuncCalled)
	}))
}

func Test_ErrorTimeout(t *testing.T) {
	require.NoError(t, helpers.FuncWithTimeout(5*time.Second, func() {
		var wantAfterErr error
		var wantAfterFuncCalled bool

		sd := NewShutdowner(
			func(name string, err error) {
				wantAfterFuncCalled = true
				wantAfterErr = err
			},
			WithProcess("test", func(ctx context.Context) error {
				time.Sleep(2 * time.Second)
				return nil
			}),
		)
		sd.Wait(context.Background(), time.Second)

		require.EqualError(t, wantAfterErr, "context deadline exceeded")
		require.True(t, wantAfterFuncCalled)
	}))
}

func Test_ErrorClose(t *testing.T) {
	require.NoError(t, helpers.FuncWithTimeout(5*time.Second, func() {
		var wantAfterErr error
		var wantAfterFuncCalled bool

		sd := NewShutdowner(
			func(name string, err error) {
				wantAfterFuncCalled = true
				wantAfterErr = err
			},
			WithProcess("test", func(ctx context.Context) error {
				return errors.New("close error")
			}),
		)
		sd.Wait(context.Background(), time.Second)

		require.EqualError(t, wantAfterErr, "close error")
		require.True(t, wantAfterFuncCalled)
	}))
}
