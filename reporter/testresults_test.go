package reporter_test

import (
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/measurement/mocks"
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
		&sciensano.APITestResultsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)},
			Region:    "Flanders",
			Total:     0,
			Positive:  0,
		},
	}
)

func TestGetTests(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "TestResults").Return(testResultsResponse, true)

	result, err := client.GetTestResults()
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.March, 12, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "total", Values: []float64{200, 100, 0}},
			{Name: "positive", Values: []float64{20, 20, 0}},
			{Name: "rate", Values: []float64{0.1, 0.2, 0}},
		},
	}, result)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetTestResults_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "TestResults").Return(nil, false)

	client := reporter.New(time.Hour)
	client.APICache = cache

	_, err := client.GetTestResults()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}
