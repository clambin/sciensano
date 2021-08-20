package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"github.com/clambin/sciensano/sciensano/apiclient/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResultsResponse = []*apiclient.APITestResultsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Brussels",
			Total:     100,
			Positive:  10,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay.Add(24 * time.Hour)},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		},
	}
)

func TestGetTests(t *testing.T) {
	apiClient := &mocks.APIClient{}
	client := sciensano.NewClient(time.Hour)
	client.APIClient = apiClient
	ctx := context.Background()

	firstDay := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	apiClient.On("GetTestResults", mock.Anything).Return(testResultsResponse, nil)

	result, err := client.GetTests(ctx, firstDay)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, firstDay, result[0].Timestamp)
	assert.Equal(t, 200, result[0].Total)
	assert.Equal(t, 20, result[0].Positive)

	mock.AssertExpectationsForObjects(t, apiClient)
}
