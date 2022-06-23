package testresults_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/testresults"
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
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeTestResults).Return(testResultsResponse, nil)

	r := testresults.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 12, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "total", "positive", "rate"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{200, 100, 0}, values)

	values, ok = entries.GetFloatValues("positive")
	require.True(t, ok)
	assert.Equal(t, []float64{20, 20, 0}, values)

	values, ok = entries.GetFloatValues("rate")
	require.True(t, ok)
	assert.Equal(t, []float64{0.1, 0.2, 0}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetTestResults_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeTestResults).Return(nil, errors.New("fail"))

	r := testresults.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}
