package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/mocks"
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResponse = []*apiclient.APICasesResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     100,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			Cases:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     120,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     100,
		},
	}
)

func TestClient_GetCases(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	endTime := time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

	cases, err := client.GetCases(ctx, endTime)
	require.NoError(t, err)
	require.Len(t, cases, 2)
	assert.Equal(t, 250, cases[0].Count)
	assert.Equal(t, 120, cases[1].Count)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCasesByProvince(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	endTime := time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

	cases, err := client.GetCasesByProvince(ctx, endTime)
	require.NoError(t, err)
	require.Len(t, cases, 2)
	require.Contains(t, cases, "VlaamsBrabant")
	require.Len(t, cases["VlaamsBrabant"], 2)
	assert.Equal(t, 100, cases["VlaamsBrabant"][0].Count)
	assert.Equal(t, 120, cases["VlaamsBrabant"][1].Count)
	require.Contains(t, cases, "Brussels")
	require.Len(t, cases["Brussels"], 1)
	assert.Equal(t, 150, cases["Brussels"][0].Count)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCasesByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	endTime := time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

	cases, err := client.GetCasesByRegion(ctx, endTime)
	require.NoError(t, err)
	require.Len(t, cases, 2)
	require.Contains(t, cases, "Flanders")
	require.Len(t, cases["Flanders"], 2)
	assert.Equal(t, 100, cases["Flanders"][0].Count)
	assert.Equal(t, 120, cases["Flanders"][1].Count)
	require.Contains(t, cases, "Brussels")
	require.Len(t, cases["Brussels"], 1)
	assert.Equal(t, 150, cases["Brussels"][0].Count)

	apiClient.AssertExpectations(t)
}
