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
			AgeGroup:  "85+",
			Cases:     100,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "55-64",
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
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 1)
	assert.Empty(t, cases.Groups[0].Name)
	assert.Equal(t, []int{250, 120}, cases.Groups[0].Values)
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
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "Brussels", cases.Groups[0].Name)
	assert.Equal(t, []int{150, 0}, cases.Groups[0].Values)

	assert.Equal(t, "VlaamsBrabant", cases.Groups[1].Name)
	assert.Equal(t, []int{100, 120}, cases.Groups[1].Values)

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
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "Brussels", cases.Groups[0].Name)
	assert.Equal(t, []int{150, 0}, cases.Groups[0].Values)

	assert.Equal(t, "Flanders", cases.Groups[1].Name)
	assert.Equal(t, []int{100, 120}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCasesByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	endTime := time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

	cases, err := client.GetCasesByAgeGroup(ctx, endTime)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "25-34", cases.Groups[0].Name)
	assert.Equal(t, []int{150, 120}, cases.Groups[0].Values)

	assert.Equal(t, "85+", cases.Groups[1].Name)
	assert.Equal(t, []int{100, 0}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func BenchmarkClient_GetCasesByRegion(b *testing.B) {
	var bigResponse []*apiclient.APICasesResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &apiclient.APICasesResponse{
				TimeStamp: apiclient.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				Cases:     i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	endTime := time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	for i := 0; i < 100; i++ {
		_, err := client.GetCasesByRegion(ctx, endTime)
		require.NoError(b, err)
	}
}
