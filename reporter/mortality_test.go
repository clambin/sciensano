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
	testMortalityResponse = []measurement.Measurement{
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Deaths:    100,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "25-34",
			Deaths:    150,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Deaths:    120,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "55-64",
			Deaths:    100,
		},
	}
)

func TestClient_GetMortality(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortality(ctx)
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "total", Values: []float64{250, 120, 100}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetMortalityByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortalityByRegion(ctx)
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "Brussels", Values: []float64{150, 0, 0}},
			{Name: "Flanders", Values: []float64{100, 120, 100}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetMortalityByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortalityByAgeGroup(ctx)
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "25-34", Values: []float64{150, 120, 0}},
			{Name: "55-64", Values: []float64{0, 0, 100}},
			{Name: "85+", Values: []float64{100, 0, 0}},
		},
	}, cases)

	mock.AssertExpectationsForObjects(t, apiClient)
}

func BenchmarkClient_GetMortalityByRegion(b *testing.B) {
	var bigResponse sciensano.APIMortalityResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APIMortalityResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Deaths:    i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	for i := 0; i < 100; i++ {
		_, err := client.GetMortalityByRegion(ctx)
		require.NoError(b, err)
	}
}
