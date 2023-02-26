package flipt

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/pkg/env"

	"github.com/stretchr/testify/require"
	"go.flipt.io/flipt-grpc"
)

func Test_New_Success(t *testing.T) {
	err := env.Init("../../../.env")
	require.NoError(t, err)

	fliptClient, cleanUp, err := New(
		WithServer(os.Getenv("APP_FLIPT_SERVER")),
		WithInsecure(true),
	)
	require.Implements(t, (*flipt.FliptClient)(nil), fliptClient)
	require.NotNil(t, cleanUp)
	require.Nil(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = fliptClient.DeleteFlag(ctx, &flipt.DeleteFlagRequest{
		Key: "some_unkown_test_key",
	})
	require.NoError(t, err)

	flag, err := fliptClient.CreateFlag(ctx, &flipt.CreateFlagRequest{
		Name:    "some unkown test key",
		Key:     "some_unkown_test_key",
		Enabled: true,
	})
	require.NoError(t, err)
	require.EqualValues(t, "some unkown test key", flag.Name)
	require.EqualValues(t, "some_unkown_test_key", flag.Key)
	require.True(t, flag.Enabled)

	require.NoError(t, cleanUp())
}

func Test_New_ErrorOptions(t *testing.T) {
	flipt, cleanUp, err := New(WithServer(""))

	require.EqualError(t, err, "new options: apply: server can't be blank")
	require.ErrorIs(t, err, ErrServerBlank)
	require.Nil(t, flipt)
	require.Nil(t, cleanUp)
}

func Test_New_ErrorDial(t *testing.T) {
	flipt, cleanUp, err := New(WithServer("localhost:123"))

	require.EqualError(t, err, "grpc dial: grpc: no transport security set (use grpc.WithTransportCredentials(insecure.NewCredentials()) explicitly or set credentials)")
	require.Nil(t, flipt)
	require.Nil(t, cleanUp)
}
