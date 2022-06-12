package cases_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache/mocks"
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
	h := &mocks.Holder{}
	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
	}

	h.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	c, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, c.GetTimestamps())

	require.Equal(t, []string{"time", "total"}, c.GetColumns())

	values, ok := c.GetFloatValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{250, 120}, values)

	mock.AssertExpectationsForObjects(t, h)
}

func TestClient_GetByProvince(t *testing.T) {
	h := &mocks.Holder{}
	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
	}

	h.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	c, err := r.GetByProvince()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, c.GetTimestamps())

	assert.Equal(t, []string{"time", "Brussels", "VlaamsBrabant"}, c.GetColumns())

	values, ok := c.GetFloatValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = c.GetFloatValues("VlaamsBrabant")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

	mock.AssertExpectationsForObjects(t, h)
}

func TestClient_GetByRegion(t *testing.T) {
	h := &mocks.Holder{}

	h.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
	}

	c, err := r.GetByRegion()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, c.GetTimestamps())

	assert.Equal(t, []string{"time", "Brussels", "Flanders"}, c.GetColumns())

	values, ok := c.GetFloatValues("Brussels")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 0}, values)

	values, ok = c.GetFloatValues("Flanders")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 120}, values)

	mock.AssertExpectationsForObjects(t, h)
}

func TestClient_GetCasesByAgeGroup(t *testing.T) {
	h := &mocks.Holder{}

	h.
		On("Get", "Cases").
		Return(testCasesResponse, true)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
	}

	c, err := r.GetByAgeGroup()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, c.GetTimestamps())

	assert.Equal(t, []string{"time", "25-34", "85+"}, c.GetColumns())

	values, ok := c.GetFloatValues("25-34")
	require.True(t, ok)
	assert.Equal(t, []float64{150, 120}, values)

	values, ok = c.GetFloatValues("85+")
	require.True(t, ok)
	assert.Equal(t, []float64{100, 0}, values)

	mock.AssertExpectationsForObjects(t, h)
}

func TestClient_GetCases_Failure(t *testing.T) {
	h := &mocks.Holder{}
	h.On("Get", "Cases").Return(nil, false)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
	}

	_, err := r.Get()
	require.Error(t, err)

	_, err = r.GetByRegion()
	require.Error(t, err)

	_, err = r.GetByProvince()
	require.Error(t, err)

	_, err = r.GetByAgeGroup()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, h)
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
	h := &mocks.Holder{}
	h.
		On("Get", "Cases").
		Return(bigResponse, true)

	r := cases.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APICache:    h,
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
