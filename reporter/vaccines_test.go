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
			Date:   vaccines.Time{Time: timestamp},
			Amount: 10,
		},
		&vaccines.Batch{
			Date:   vaccines.Time{Time: timestamp},
			Amount: 20,
		},
		&vaccines.Batch{
			Date:   vaccines.Time{Time: timestamp.Add(24 * time.Hour)},
			Amount: 40,
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

func TestClient_GetVaccines_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccines").Return(nil, false)

	client := reporter.New(time.Hour)
	client.Vaccines = cache

	_, err := client.GetVaccines()
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}
