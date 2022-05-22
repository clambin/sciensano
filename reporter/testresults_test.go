package reporter_test

import (
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResultsResponse = []apiclient.APIResponse{
		&sciensano.APITestResultsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		},
		&sciensano.APITestResultsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Total:     100,
			Positive:  10,
		},
		&sciensano.APITestResultsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Total:     100,
			Positive:  20,
		},
		&sciensano.APITestResultsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, time.March, 12, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Total:     0,
			Positive:  0,
		},
	}
)

func TestGetTests(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "TestResults").Return(testResultsResponse, true).Once()

	c := reporter.NewWithOptions(time.Hour, client.Options{})
	c.APICache = cache

	entries, err := c.GetTestResults()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 12, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"positive", "rate", "total"}, entries.GetColumns())

	values, ok := entries.GetValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{200, 100, 0}, values)

	values, ok = entries.GetValues("positive")
	require.True(t, ok)
	assert.Equal(t, []float64{20, 20, 0}, values)

	values, ok = entries.GetValues("rate")
	require.True(t, ok)
	assert.Equal(t, []float64{0.1, 0.2, 0}, values)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetTestResults_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "TestResults").Return(nil, false).Once()

	c := reporter.NewWithOptions(time.Hour, client.Options{})
	c.APICache = cache

	_, err := c.GetTestResults()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}
