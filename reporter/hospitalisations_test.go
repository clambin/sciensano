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
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	entries, err := client.GetHospitalisations()
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

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetHospitalisationsByProvince(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	entries, err := client.GetHospitalisationsByProvince()
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
	cache.AssertExpectations(t)
}

func TestClient_GetHospitalisationsByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	entries, err := client.GetHospitalisationsByRegion()
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
	cache.AssertExpectations(t)
}

func TestClient_GetHospitalisations_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Hospitalisations").Return(nil, false)

	_, err := client.GetHospitalisations()
	require.Error(t, err)

	_, err = client.GetHospitalisationsByRegion()
	require.Error(t, err)

	_, err = client.GetHospitalisationsByProvince()
	require.Error(t, err)

	cache.AssertExpectations(t)
}

func TestClient_GetHospitalisations_ApplyRegions(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	entries, err := client.GetHospitalisationsByRegion()
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

	entries, err = client.GetHospitalisationsByRegion()
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 2)
	require.Len(t, entries.Groups[1].Values, 2)

	cache.AssertExpectations(t)
}

func BenchmarkClient_GetHospitalisationsByRegion(b *testing.B) {
	var bigResponse []measurement.Measurement
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
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Hospitalisations").
		Return(bigResponse, true)

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		_, err := client.GetHospitalisationsByRegion()
		require.NoError(b, err)
	}

	cache.AssertExpectations(b)
}
