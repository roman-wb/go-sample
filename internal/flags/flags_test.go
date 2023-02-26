package flags

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/roman-wb/go-sample/internal/flags/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.flipt.io/flipt-grpc"
)

func Test_New(t *testing.T) {
	fliptClient := &mocks.FliptClient{}
	flags := New(fliptClient, time.Second)

	require.EqualValues(t, flags.client, fliptClient)
	require.EqualValues(t, flags.timeout, time.Second)
}

func Test_IsSampleFeatureEnabled_Success_True(t *testing.T) {
	fliptClient := &mocks.FliptClient{}
	flags := New(fliptClient, time.Second)

	fliptClient.On("GetFlag", mock.Anything, &flipt.GetFlagRequest{
		Key: SampleFeatureKey,
	}).Return(&flipt.Flag{
		Key:     SampleFeatureKey,
		Enabled: true,
	}, nil).Once()
	defer fliptClient.AssertExpectations(t)

	isEnabled, err := flags.IsSampleFeatureEnabled(context.Background())
	require.NoError(t, err)
	require.True(t, isEnabled)
}

func Test_IsSampleFeatureEnabled_Success_False(t *testing.T) {
	fliptClient := &mocks.FliptClient{}
	flags := New(fliptClient, time.Second)

	fliptClient.On("GetFlag", mock.Anything, &flipt.GetFlagRequest{
		Key: SampleFeatureKey,
	}).Return(&flipt.Flag{
		Key:     SampleFeatureKey,
		Enabled: false,
	}, nil).Once()

	defer fliptClient.AssertExpectations(t)

	isEnabled, err := flags.IsSampleFeatureEnabled(context.Background())
	require.NoError(t, err)
	require.False(t, isEnabled)
}

func Test_IsSampleFeatureEnabled_Error(t *testing.T) {
	fliptClient := &mocks.FliptClient{}
	flags := New(fliptClient, time.Second)

	fliptClient.On("GetFlag", mock.Anything, &flipt.GetFlagRequest{
		Key: SampleFeatureKey,
	}).Return(nil, errors.New("some error")).Once()

	defer fliptClient.AssertExpectations(t)

	isEnabled, err := flags.IsSampleFeatureEnabled(context.Background())
	require.EqualError(t, err, "get flag: some error")
	require.False(t, isEnabled)
}
