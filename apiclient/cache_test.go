package apiclient_test

import (
	"context"
	"fmt"
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
	timestamp := time.Now()

	// Cache should only call the client once.
	client.
		On("GetVaccinations", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
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
	assert.Equal(t, &apiclient.APIVaccinationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		AgeGroup:  "25-34",
		Gender:    "M",
		Dose:      "A",
		Count:     100,
	}, results[0])

	results, err = cache.GetVaccinations(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APIVaccinationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		AgeGroup:  "25-34",
		Gender:    "M",
		Dose:      "A",
		Count:     100,
	}, results[0])

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_GetVaccinations_Failure(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()
	timestamp := time.Now()

	// Set up a failing call
	client.
		On("GetVaccinations", mock.Anything).
		Return(nil, fmt.Errorf("not available")).
		Once()

	results, err := cache.GetVaccinations(ctx)
	require.Error(t, err)

	// Cache should only call the client once.
	client.
		On("GetVaccinations", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Region:    "Flanders",
				AgeGroup:  "25-34",
				Gender:    "M",
				Dose:      "A",
				Count:     100,
			},
		}, nil).
		Once()

	results, err = cache.GetVaccinations(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APIVaccinationsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		AgeGroup:  "25-34",
		Gender:    "M",
		Dose:      "A",
		Count:     100,
	}, results[0])

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_GetTestResults(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()
	timestamp := time.Now()

	// Cache should only call the client once.
	client.
		On("GetTestResults", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APITestResultsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Region:    "Flanders",
				Province:  "VlaamsBrabant",
				Total:     100,
				Positive:  10,
			},
		}, nil).
		Once()

	results, err := cache.GetTestResults(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APITestResultsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		Province:  "VlaamsBrabant",
		Total:     100,
		Positive:  10,
	}, results[0])

	results, err = cache.GetTestResults(ctx)
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APITestResultsResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		Province:  "VlaamsBrabant",
		Total:     100,
		Positive:  10,
	}, results[0])

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_GetCases(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()
	timestamp := time.Now()

	// Cache should only call the client once.
	client.
		On("GetCases", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APICasesResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Region:    "Flanders",
				Province:  "VlaamsBrabant",
				AgeGroup:  "85+",
				Cases:     10,
			},
		}, nil).
		Once()

	results, err := cache.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APICasesResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Province:  "VlaamsBrabant",
		Region:    "Flanders",
		AgeGroup:  "85+",
		Cases:     10,
	}, results[0])

	results, err = cache.GetCases(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APICasesResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Province:  "VlaamsBrabant",
		Region:    "Flanders",
		AgeGroup:  "85+",
		Cases:     10,
	}, results[0])

	mock.AssertExpectationsForObjects(t, client)

}

func TestCache_GetMortality(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()
	timestamp := time.Now()

	// Cache should only call the client once.
	client.
		On("GetMortality", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APIMortalityResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Region:    "Flanders",
				AgeGroup:  "85+",
				Deaths:    10,
			},
		}, nil).
		Once()

	results, err := cache.GetMortality(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APIMortalityResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		AgeGroup:  "85+",
		Deaths:    10,
	}, results[0])

	results, err = cache.GetMortality(ctx)

	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, &apiclient.APIMortalityResponseEntry{
		TimeStamp: apiclient.TimeStamp{Time: timestamp},
		Region:    "Flanders",
		AgeGroup:  "85+",
		Deaths:    10,
	}, results[0])

	mock.AssertExpectationsForObjects(t, client)
}

func TestCache_All(t *testing.T) {
	client := &mocks.Getter{}
	cache := apiclient.Cache{
		Getter:    client,
		Retention: time.Hour,
	}
	ctx := context.Background()

	client.
		On("GetTestResults", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APITestResultsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now()},
				Region:    "Flanders",
				Total:     100,
				Positive:  10,
			},
		}, nil).
		Once()
	client.
		On("GetVaccinations", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now()},
				Region:    "Flanders",
				AgeGroup:  "25-34",
				Gender:    "M",
				Dose:      "A",
				Count:     100,
			},
		}, nil).
		Once()

	client.
		On("GetCases", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APICasesResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now()},
				Region:    "Flanders",
				Province:  "VlaamsBrabant",
				Cases:     10,
			},
		}, nil).
		Once()

	client.
		On("GetMortality", mock.Anything).
		Return([]apiclient.Measurement{
			&apiclient.APIMortalityResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now()},
				Region:    "Flanders",
				Deaths:    10,
			},
		}, nil).
		Once()

	for i := 0; i < 500; i++ {
		go func() {
			_, err := cache.GetTestResults(ctx)
			require.NoError(t, err)
		}()
		go func() {
			_, err := cache.GetVaccinations(ctx)
			require.NoError(t, err)
		}()
		go func() {
			_, err := cache.GetCases(ctx)
			require.NoError(t, err)
		}()
		go func() {
			_, err := cache.GetMortality(ctx)
			require.NoError(t, err)
		}()
	}

	mock.AssertExpectationsForObjects(t, client)
}
