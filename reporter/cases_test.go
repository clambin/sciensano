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
	testCasesResponse = []measurement.Measurement{
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
	}
)

func TestClient_GetCases(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCases()
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "total", Values: []float64{250, 120}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetCasesByProvince(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCasesByProvince()
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "Brussels", Values: []float64{150, 0}},
			{Name: "VlaamsBrabant", Values: []float64{100, 120}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetCasesByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCasesByRegion()
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "Brussels", Values: []float64{150, 0}},
			{Name: "Flanders", Values: []float64{100, 120}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetCasesByAgeGroup(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCasesByAgeGroup()
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "25-34", Values: []float64{150, 120}},
			{Name: "85+", Values: []float64{100, 0}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetCases_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Cases").Return(nil, false)

	_, err := client.GetCases()
	require.Error(t, err)

	_, err = client.GetCasesByRegion()
	require.Error(t, err)

	_, err = client.GetCasesByProvince()
	require.Error(t, err)

	_, err = client.GetCasesByAgeGroup()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_Cases_ApplyRegions(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCasesByAgeGroup()
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)

	cases.ApplyRange(time.Time{}, time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC))
	require.Len(t, cases.Timestamps, 1)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 1)
	require.Len(t, cases.Groups[1].Values, 1)

	cases, err = client.GetCasesByAgeGroup()
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)
}

func BenchmarkClient_GetCasesByRegion(b *testing.B) {
	var bigResponse sciensano.APICasesResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APICasesResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				Cases:     i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		_, err := client.GetCasesByRegion()
		require.NoError(b, err)
	}
}
