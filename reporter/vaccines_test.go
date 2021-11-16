package reporter_test

import (
	"github.com/clambin/sciensano/apiclient/vaccines"
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
	testVaccinesResponse = []measurement.Measurement{
		&vaccines.Batch{
			Date:         vaccines.Timestamp{Time: timestamp},
			Manufacturer: "A",
			Amount:       10,
		},
		&vaccines.Batch{
			Date:         vaccines.Timestamp{Time: timestamp},
			Manufacturer: "B",
			Amount:       20,
		},
		&vaccines.Batch{
			Date:         vaccines.Timestamp{Time: timestamp.Add(24 * time.Hour)},
			Manufacturer: "A",
			Amount:       40,
		},
	}
)

func TestClient_GetVaccines(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(testVaccinesResponse, true)

	client := reporter.New(time.Hour)
	client.Vaccines = cache

	result, err := client.GetVaccines()
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			timestamp,
			timestamp.Add(24 * time.Hour),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "total", Values: []float64{30, 40}},
		},
	}, result)
}

func TestClient_GetVaccinesByManufacturer(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(testVaccinesResponse, true)

	client := reporter.New(time.Hour)
	client.Vaccines = cache

	result, err := client.GetVaccinesByManufacturer()
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			timestamp,
			timestamp.Add(24 * time.Hour),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "A", Values: []float64{10, 40}},
			{Name: "B", Values: []float64{20, 0}},
		},
	}, result)
}

func TestClient_GetVaccines_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(nil, false)

	client := reporter.New(time.Hour)
	client.Vaccines = cache

	_, err := client.GetVaccines()
	require.Error(t, err)

	_, err = client.GetVaccinesByManufacturer()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}
