package cases_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/cases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testCasesResponse = []apiclient.APIResponse{
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
	}
)

func TestClient_Get(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(testCasesResponse, nil)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	require.Equal(t, []string{"time", "total"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{250, 120}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetByProvince(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(testCasesResponse, nil)

	r := cases.Reporter{
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
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = entries.GetFloatValues("VlaamsBrabant")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetByRegion(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(testCasesResponse, nil)

	r := cases.Reporter{
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
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = entries.GetFloatValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetCasesByAgeGroup(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(testCasesResponse, nil)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByAgeGroup()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "25-34", "85+"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("25-34")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 120}, values)

	values, ok = entries.GetFloatValues("85+")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 0}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetCases_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(nil, errors.New("fail"))

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	require.Error(t, err)

	_, err = r.GetByRegion()
	require.Error(t, err)

	_, err = r.GetByProvince()
	require.Error(t, err)

	_, err = r.GetByAgeGroup()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkClient_GetCasesByRegion(b *testing.B) {
	var bigResponse []apiclient.APIResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APICasesResponse{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				Cases:     i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(bigResponse, nil)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetByRegion()
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
