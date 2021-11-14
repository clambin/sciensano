package reporter_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/mocks"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResultsResponse = []measurement.Measurement{
		&sciensano.APITestResultsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		},
		&sciensano.APITestResultsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			Total:     100,
			Positive:  10,
		},
		&sciensano.APITestResultsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)},
			Region:    "Flanders",
			Total:     100,
			Positive:  20,
		},
	}
)

func TestGetTests(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.On("GetTestResults", mock.Anything).Return(testResultsResponse, nil)

	result, err := client.GetTestResults(ctx)
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "total", Values: []float64{200, 100}},
			{Name: "positive", Values: []float64{20, 20}},
		},
	}, result)

	mock.AssertExpectationsForObjects(t, apiClient)
}
