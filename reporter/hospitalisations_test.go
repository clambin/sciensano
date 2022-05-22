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
	testHospitalisationsResponse = []apiclient.APIResponse{
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     100,
			TotalInICU:  50,
			TotalInResp: 25,
			TotalInECMO: 10,
		},
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Brussels",
			Province:    "Brussels",
			TotalIn:     10,
			TotalInICU:  5,
			TotalInResp: 3,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponse{
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
	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetHospitalisations()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"in", "inECMO", "inICU", "inResp"}, entries.GetColumns())

	for _, testCase := range []struct {
		column   string
		expected []float64
	}{
		{column: "in", expected: []float64{110, 50}},
		{column: "inICU", expected: []float64{55, 25}},
		{column: "inResp", expected: []float64{28, 12}},
		{column: "inECMO", expected: []float64{11, 5}},
	} {
		values, ok := entries.GetValues(testCase.column)
		require.True(t, ok, testCase.column)
		assert.Equal(t, testCase.expected, values, testCase.column)
	}

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetHospitalisationsByProvince(t *testing.T) {
	cache := &mocks.Holder{}
	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetHospitalisationsByProvince()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"Brussels", "VlaamsBrabant"}, entries.GetColumns())

	values, ok := entries.GetValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 0}, values)

	values, ok = entries.GetValues("VlaamsBrabant")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 50}, values)

	cache.AssertExpectations(t)
}

func TestClient_GetHospitalisationsByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	cache.
		On("Get", "Hospitalisations").
		Return(testHospitalisationsResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetHospitalisationsByRegion()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"Brussels", "Flanders"}, entries.GetColumns())

	values, ok := entries.GetValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 0}, values)

	values, ok = entries.GetValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 50}, values)

	cache.AssertExpectations(t)
}

func TestClient_GetHospitalisations_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Hospitalisations").Return(nil, false)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	_, err := r.GetHospitalisations()
	require.Error(t, err)

	_, err = r.GetHospitalisationsByRegion()
	require.Error(t, err)

	_, err = r.GetHospitalisationsByProvince()
	require.Error(t, err)

	cache.AssertExpectations(t)
}

func BenchmarkClient_GetHospitalisationsByRegion(b *testing.B) {
	b.StopTimer()
	var bigResponse []apiclient.APIResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APIHospitalisationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				TotalIn:   i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	cache := &mocks.Holder{}
	cache.
		On("Get", "Hospitalisations").
		Return(bigResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetHospitalisationsByRegion()
		if err != nil {
			b.Fatal(err)
		}
	}
}
