package shutdowner

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewProcess(t *testing.T) {
	cancelFunc := func(ctx context.Context) error {
		return errors.New("same error")
	}
	proc := WithProcess("test", cancelFunc)

	require.EqualValues(t, "test", proc.name)
	require.EqualError(t, proc.cancelFunc(context.Background()), "same error")
}
