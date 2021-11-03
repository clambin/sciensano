package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/mocks"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testMortalityResponse = apiclient.APIMortalityResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Deaths:    100,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "25-34",
			Deaths:    150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Deaths:    120,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "55-64",
			Deaths:    100,
		},
	}
)

func TestClient_GetMortality(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortality(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 3)
	require.Len(t, cases.Groups, 1)
	assert.Empty(t, cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 250},
		&sciensano.MortalityEntry{Count: 120},
		&sciensano.MortalityEntry{Count: 100},
	}, cases.Groups[0].Values)
}

func TestClient_GetMortalityByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortalityByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 3)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "Brussels", cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 150},
		&sciensano.MortalityEntry{Count: 0},
		&sciensano.MortalityEntry{Count: 0},
	}, cases.Groups[0].Values)

	assert.Equal(t, "Flanders", cases.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 100},
		&sciensano.MortalityEntry{Count: 120},
		&sciensano.MortalityEntry{Count: 100},
	}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func TestClient_GetMortalityByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testMortalityResponse, nil)

	cases, err := client.GetMortalityByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 3)
	require.Len(t, cases.Groups, 3)

	assert.Equal(t, "25-34", cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 150},
		&sciensano.MortalityEntry{Count: 120},
		&sciensano.MortalityEntry{Count: 0},
	}, cases.Groups[0].Values)

	assert.Equal(t, "55-64", cases.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 0},
		&sciensano.MortalityEntry{Count: 0},
		&sciensano.MortalityEntry{Count: 100},
	}, cases.Groups[1].Values)

	assert.Equal(t, "85+", cases.Groups[2].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.MortalityEntry{Count: 100},
		&sciensano.MortalityEntry{Count: 0},
		&sciensano.MortalityEntry{Count: 0},
	}, cases.Groups[2].Values)

	apiClient.AssertExpectations(t)
}

func BenchmarkClient_GetMortalityByRegion(b *testing.B) {
	var bigResponse apiclient.APIMortalityResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, apiclient.APIMortalityResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: ts},
				Region:    region,
				Deaths:    i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	for i := 0; i < 100; i++ {
		_, err := client.GetMortalityByRegion(ctx)
		require.NoError(b, err)
	}
}
