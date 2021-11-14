package reporter_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/apiclient/vaccines/mocks"
	"github.com/clambin/sciensano/measurement"
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
	apiClient := &mocks.Getter{}
	apiClient.On("GetBatches", mock.Anything).Return(testVaccinesResponse, nil)

	client := reporter.NewCachedClient(time.Hour)
	client.Vaccines = apiClient

	result, err := client.GetVaccines(context.Background())
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
