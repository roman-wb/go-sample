package flipt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DefaultOptions(t *testing.T) {
	opts := DefaultOptions()

	require.EqualValues(t, defaultServer, opts.Server)
	require.EqualValues(t, defaultInsecure, opts.Insecure)
}

func Test_NewOptions_Empty(t *testing.T) {
	opts, err := NewOptions()

	require.NoError(t, err)
	require.EqualValues(t, defaultServer, opts.Server)
	require.EqualValues(t, defaultInsecure, opts.Insecure)
}

func Test_WithServer_Success(t *testing.T) {
	opts, err := NewOptions(WithServer("localhost:9000"))

	require.NoError(t, err)
	require.EqualValues(t, "localhost:9000", opts.Server)
}

func Test_WithServer_Error(t *testing.T) {
	opts, err := NewOptions(WithServer(""))

	require.EqualError(t, err, "apply: server can't be blank")
	require.ErrorIs(t, err, ErrServerBlank)
	require.Nil(t, opts)
}

func Test_WithInsecure_True(t *testing.T) {
	opts, err := NewOptions(WithInsecure(true))

	require.NoError(t, err)
	require.True(t, opts.Insecure)
}

func Test_WithInsecure_False(t *testing.T) {
	opts, err := NewOptions(WithInsecure(false))

	require.NoError(t, err)
	require.False(t, opts.Insecure)
}
