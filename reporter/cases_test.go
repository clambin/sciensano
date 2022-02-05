package reporter_test

import (
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

func TestClient_GetCases(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	cases, err := client.GetCases()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, cases.GetTimestamps())

	require.Equal(t, []string{"total"}, cases.GetColumns())

	values, ok := cases.GetValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{250, 120}, values)

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

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, cases.GetTimestamps())

	assert.Equal(t, []string{"Brussels", "VlaamsBrabant"}, cases.GetColumns())

	values, ok := cases.GetValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = cases.GetValues("VlaamsBrabant")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

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

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, cases.GetTimestamps())

	assert.Equal(t, []string{"Brussels", "Flanders"}, cases.GetColumns())

	values, ok := cases.GetValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = cases.GetValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

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

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, cases.GetTimestamps())

	assert.Equal(t, []string{"25-34", "85+"}, cases.GetColumns())

	values, ok := cases.GetValues("25-34")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 120}, values)

	values, ok = cases.GetValues("85+")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 0}, values)

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
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.
		On("Get", "Cases").
		Return(bigResponse, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetCasesByRegion()
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
