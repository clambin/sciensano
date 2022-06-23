package vaccines_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	vaccines2 "github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/vaccines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testVaccinesResponse = []apiclient.APIResponse{
		&vaccines2.APIBatchResponse{
			Date:         vaccines2.Timestamp{Time: time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Amount:       10,
		},
		&vaccines2.APIBatchResponse{
			Date:         vaccines2.Timestamp{Time: time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "B",
			Amount:       20,
		},
		&vaccines2.APIBatchResponse{
			Date:         vaccines2.Timestamp{Time: time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Amount:       40,
		},
	}
)

func TestClient_GetVaccines(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines2.TypeBatches).Return(testVaccinesResponse, nil)

	r := vaccines.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "total"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{30, 40}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetVaccinesByManufacturer(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines2.TypeBatches).Return(testVaccinesResponse, nil)

	r := vaccines.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.GetByManufacturer()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "A", "B"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 40}, values)

	values, ok = entries.GetFloatValues("B")
	require.True(t, ok)
	assert.Equal(t, []float64{20, 0}, values)

	mock.AssertExpectationsForObjects(t, f)
}

func TestClient_GetVaccines_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines2.TypeBatches).Return(nil, errors.New("fail"))

	r := vaccines.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	require.Error(t, err)

	_, err = r.GetByManufacturer()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}
