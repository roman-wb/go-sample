package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Init_Success(t *testing.T) {
	err := Init("testdata/env")

	require.NoError(t, err)
	require.EqualValues(t, "same value", os.Getenv("TEST_VAR_NAME"))
}

func Test_Init_Error(t *testing.T) {
	err := Init("none")

	require.EqualError(t, err, "godotenv load none file: open none: no such file or directory")
}
