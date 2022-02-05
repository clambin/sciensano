package reporter_test

import (
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testVaccinesResponse = []apiclient.APIResponse{
		&vaccines.APIBatchResponse{
			Date:         vaccines.Timestamp{Time: time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Amount:       10,
		},
		&vaccines.APIBatchResponse{
			Date:         vaccines.Timestamp{Time: time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "B",
			Amount:       20,
		},
		&vaccines.APIBatchResponse{
			Date:         vaccines.Timestamp{Time: time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "A",
			Amount:       40,
		},
	}
)

func TestClient_GetVaccines(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(testVaccinesResponse, true)

	client := reporter.New(time.Hour)
	client.APICache = cache

	entries, err := client.GetVaccines()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"total"}, entries.GetColumns())

	values, ok := entries.GetValues("total")
	require.True(t, ok)
	assert.Equal(t, []float64{30, 40}, values)
}

func TestClient_GetVaccinesByManufacturer(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(testVaccinesResponse, true)

	client := reporter.New(time.Hour)
	client.APICache = cache

	entries, err := client.GetVaccinesByManufacturer()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"A", "B"}, entries.GetColumns())

	values, ok := entries.GetValues("A")
	require.True(t, ok)
	assert.Equal(t, []float64{10, 40}, values)

	values, ok = entries.GetValues("B")
	require.True(t, ok)
	assert.Equal(t, []float64{20, 0}, values)
}

func TestClient_GetVaccines_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(nil, false)

	client := reporter.New(time.Hour)
	client.APICache = cache

	_, err := client.GetVaccines()
	require.Error(t, err)

	_, err = client.GetVaccinesByManufacturer()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}
