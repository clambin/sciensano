package apiclient_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache_GetVaccinations(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()

	// Cache should only call the client once.
	client.
		On("GetVaccinations", mock.Anything).
		Return([]*apiclient.APIVaccinationsResponse{{
			TimeStamp: apiclient.TimeStamp{Time: time.Now()},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Gender:    "M",
			Dose:      "A",
			Count:     100,
		}}, nil).
		Once()

	results, err := cache.GetVaccinations(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, 100, results[0].Count)

	results, err = cache.GetVaccinations(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, 100, results[0].Count)

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_GetTestResults(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()

	// Cache should only call the client once.
	client.
		On("GetTestResults", mock.Anything).
		Return([]*apiclient.APITestResultsResponse{{
			TimeStamp: apiclient.TimeStamp{Time: time.Now()},
			Region:    "Flanders",
			Total:     100,
			Positive:  10,
		}}, nil).
		Once()

	results, err := cache.GetTestResults(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, 100, results[0].Total)
	assert.Equal(t, 10, results[0].Positive)

	results, err = cache.GetTestResults(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, 100, results[0].Total)
	assert.Equal(t, 10, results[0].Positive)

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_GetCases(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()

	// Cache should only call the client once.
	client.
		On("GetCases", mock.Anything).
		Return([]*apiclient.APICasesResponse{{
			TimeStamp: apiclient.TimeStamp{Time: time.Now()},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     10,
		}}, nil).
		Once()

	results, err := cache.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)

	assert.NotZero(t, results[0].TimeStamp)
	assert.Equal(t, "VlaamsBrabant", results[0].Province)
	assert.Equal(t, 10, results[0].Cases)

	results, err = cache.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)

	assert.NotZero(t, results[0].TimeStamp)
	assert.Equal(t, "VlaamsBrabant", results[0].Province)
	assert.Equal(t, 10, results[0].Cases)

	mock.AssertExpectationsForObjects(t, client)

}
