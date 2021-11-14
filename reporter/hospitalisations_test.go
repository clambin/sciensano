package reporter_test

import (
	"context"
	"fmt"
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
	testHospitalisationsResponse = []measurement.Measurement{
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     100,
			TotalInICU:  50,
			TotalInResp: 25,
			TotalInECMO: 10,
		},
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Brussels",
			Province:    "Brussels",
			TotalIn:     10,
			TotalInICU:  5,
			TotalInResp: 3,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     50,
			TotalInICU:  25,
			TotalInResp: 12,
			TotalInECMO: 5,
		},
	}
)

func TestClient_GetHospitalisations(t *testing.T) {
	getter := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisations(ctx)
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "in", Values: []float64{110, 50}},
			{Name: "inICU", Values: []float64{55, 25}},
			{Name: "inResp", Values: []float64{28, 12}},
			{Name: "inECMO", Values: []float64{11, 5}},
		},
	}, entries)

	mock.AssertExpectationsForObjects(t, getter)
}

func TestClient_GetHospitalisationsByProvince(t *testing.T) {
	getter := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByProvince(ctx)
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "Brussels", Values: []float64{10, 0}},
			{Name: "VlaamsBrabant", Values: []float64{100, 50}},
		},
	}, entries)
	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisationsByRegion(t *testing.T) {
	getter := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "Brussels", Values: []float64{10, 0}},
			{Name: "Flanders", Values: []float64{100, 50}},
		},
	}, entries)
	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisations_Failure(t *testing.T) {
	getter := &mocks.Getter{}
	getter.On("GetHospitalisations", mock.Anything).Return(nil, fmt.Errorf("API error"))

	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter

	ctx := context.Background()

	_, err := client.GetHospitalisations(ctx)
	require.Error(t, err)

	_, err = client.GetHospitalisationsByRegion(ctx)
	require.Error(t, err)

	_, err = client.GetHospitalisationsByProvince(ctx)
	require.Error(t, err)

	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisations_ApplyRegions(t *testing.T) {
	getter := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 2)
	require.Len(t, entries.Groups[1].Values, 2)

	entries.ApplyRange(time.Time{}, time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC))
	require.Len(t, entries.Timestamps, 1)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 1)
	require.Len(t, entries.Groups[1].Values, 1)

	entries, err = client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 2)
	require.Len(t, entries.Groups[1].Values, 2)

	getter.AssertExpectations(t)
}

func BenchmarkClient_GetHospitalisationsByRegion(b *testing.B) {
	var bigResponse sciensano.APIHospitalisationsResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APIHospitalisationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				TotalIn:   i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	getter := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil).Once()

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		_, err := client.GetHospitalisationsByRegion(ctx)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
