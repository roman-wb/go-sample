package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_FuncWithTimeout_Success(t *testing.T) {
	require.NoError(t, FuncWithTimeout(100*time.Millisecond, func() {
		time.Sleep(10 * time.Millisecond)
	}))
}

func Test_FuncWithTimeout_Error(t *testing.T) {
	require.EqualError(t, FuncWithTimeout(100*time.Millisecond, func() {
		time.Sleep(200 * time.Millisecond)
	}), "context deadline exceeded")
}
