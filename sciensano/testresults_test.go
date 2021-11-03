package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/mocks"
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResultsResponse = apiclient.APITestResultsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			Total:     100,
			Positive:  10,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)},
			Region:    "Flanders",
			Total:     100,
			Positive:  20,
		},
	}
)

func TestGetTests(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.On("GetTestResults", mock.Anything).Return(testResultsResponse, nil)

	result, err := client.GetTestResults(ctx)
	require.NoError(t, err)
	require.Len(t, result.Timestamps, 2)
	require.Len(t, result.Groups, 1)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, 200, result.Groups[0].Values[0].(*sciensano.TestResult).Total)
	assert.Equal(t, 20, result.Groups[0].Values[0].(*sciensano.TestResult).Positive)
	assert.Equal(t, 0.1, result.Groups[0].Values[0].(*sciensano.TestResult).Ratio())
	assert.Equal(t, 100, result.Groups[0].Values[1].(*sciensano.TestResult).Total)
	assert.Equal(t, 20, result.Groups[0].Values[1].(*sciensano.TestResult).Positive)
	assert.Equal(t, 0.2, result.Groups[0].Values[1].(*sciensano.TestResult).Ratio())

	mock.AssertExpectationsForObjects(t, apiClient)
}
