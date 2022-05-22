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
	cache := &mocks.Holder{}
	cache.
		On("Get", "Mortality").
		Return(testMortalityResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetMortality()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"total"}, entries.GetColumns())

	values, ok := entries.GetValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{250, 120, 100}, values)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetMortalityByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	cache.
		On("Get", "Mortality").
		Return(testMortalityResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetMortalityByRegion()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"Brussels", "Flanders"}, entries.GetColumns())

	values, ok := entries.GetValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0, 0}, values)

	values, ok = entries.GetValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120, 100}, values)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetMortalityByAgeGroup(t *testing.T) {
	cache := &mocks.Holder{}
	cache.
		On("Get", "Mortality").
		Return(testMortalityResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	entries, err := r.GetMortalityByAgeGroup()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 23, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"25-34", "55-64", "85+"}, entries.GetColumns())

	values, ok := entries.GetValues("25-34")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 120, 0}, values)

	values, ok = entries.GetValues("55-64")
	require.True(t, ok)
	assert.Equal(t, []float64{0, 0, 100}, values)

	values, ok = entries.GetValues("85+")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 0, 0}, values)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetMortality_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Mortality").Return(nil, false)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	_, err := r.GetMortality()
	require.Error(t, err)

	_, err = r.GetMortalityByRegion()
	require.Error(t, err)

	_, err = r.GetMortalityByAgeGroup()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}

func BenchmarkClient_GetMortalityByRegion(b *testing.B) {
	b.StopTimer()
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
	cache := &mocks.Holder{}
	cache.
		On("Get", "Mortality").
		Return(bigResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetMortalityByRegion()
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
