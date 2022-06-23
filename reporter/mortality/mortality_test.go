package mortality_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/mortality"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testMortalityResponse = []apiclient.APIResponse{
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Deaths:    100,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "25-34",
			Deaths:    150,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Deaths:    120,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "55-64",
			Deaths:    100,
		},
	}
)

func TestClient_GetMortality(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(testMortalityResponse, nil)

	r := mortality.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "total"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{250, 120, 100}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetMortalityByRegion(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(testMortalityResponse, nil)

	r := mortality.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByRegion()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "Brussels", "Flanders"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0, 0}, values)

	values, ok = entries.GetFloatValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120, 100}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetMortalityByAgeGroup(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(testMortalityResponse, nil)

	r := mortality.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByAgeGroup()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "25-34", "55-64", "85+"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("25-34")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 120, 0}, values)

	values, ok = entries.GetFloatValues("55-64")
	require.True(t, ok)
	assert.Equal(t, []float64{0, 0, 100}, values)

	values, ok = entries.GetFloatValues("85+")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 0, 0}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetMortality_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(nil, errors.New("fail"))

	r := mortality.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	require.Error(t, err)

	_, err = r.GetByRegion()
	require.Error(t, err)

	_, err = r.GetByAgeGroup()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkClient_GetMortalityByRegion(b *testing.B) {
	var bigResponse []apiclient.APIResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &sciensano.APIMortalityResponse{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Deaths:    i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(bigResponse, nil)

	r := mortality.Reporter{
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
