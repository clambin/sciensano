package hospitalisations_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/hospitalisations"
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

func TestClient_Get(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeHospitalisations).Return(testHospitalisationsResponse, nil)

	r := hospitalisations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "in", "inICU", "inResp", "inECMO"}, entries.GetColumns())

	for _, testCase := range []struct {
		column   string
		expected []float64
	}{
		{column: "in", expected: []float64{110, 50}},
		{column: "inICU", expected: []float64{55, 25}},
		{column: "inResp", expected: []float64{28, 12}},
		{column: "inECMO", expected: []float64{11, 5}},
	} {
		values, ok := entries.GetFloatValues(testCase.column)
		require.True(t, ok, testCase.column)
		assert.Equal(t, testCase.expected, values, testCase.column)
	}

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetHospitalisationsByProvince(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeHospitalisations).Return(testHospitalisationsResponse, nil)

	r := hospitalisations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByProvince()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "Brussels", "VlaamsBrabant"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 0}, values)

	values, ok = entries.GetFloatValues("VlaamsBrabant")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 50}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetHospitalisationsByRegion(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeHospitalisations).Return(testHospitalisationsResponse, nil)

	r := hospitalisations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByRegion()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "Brussels", "Flanders"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 0}, values)

	values, ok = entries.GetFloatValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 50}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetHospitalisations_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeHospitalisations).Return(nil, errors.New("fail"))

	r := hospitalisations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	require.Error(t, err)

	_, err = r.GetByRegion()
	require.Error(t, err)

	_, err = r.GetByProvince()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkClient_GetHospitalisationsByRegion(b *testing.B) {
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

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeHospitalisations).Return(bigResponse, nil)

	r := hospitalisations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetByRegion()
		if err != nil {
			b.Fatal(err)
		}
	}
}
