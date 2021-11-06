package sciensano_test

import (
	"context"
	"fmt"
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
	testCasesResponse = []apiclient.Measurement{
		&apiclient.APICasesResponseEntry{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		&apiclient.APICasesResponseEntry{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		&apiclient.APICasesResponseEntry{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
	}
)

func TestClient_GetCases(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testCasesResponse, nil)

	cases, err := client.GetCases(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 1)
	assert.Empty(t, cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 250},
		&sciensano.CasesEntry{Count: 120},
	}, cases.Groups[0].Values)
}

func TestClient_GetCasesByProvince(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testCasesResponse, nil)

	cases, err := client.GetCasesByProvince(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "Brussels", cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 150},
		&sciensano.CasesEntry{Count: 0},
	}, cases.Groups[0].Values)

	assert.Equal(t, "VlaamsBrabant", cases.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 100},
		&sciensano.CasesEntry{Count: 120},
	}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCasesByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testCasesResponse, nil)

	cases, err := client.GetCasesByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "Brussels", cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 150},
		&sciensano.CasesEntry{Count: 0},
	}, cases.Groups[0].Values)

	assert.Equal(t, "Flanders", cases.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 100},
		&sciensano.CasesEntry{Count: 120},
	}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCasesByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testCasesResponse, nil)

	cases, err := client.GetCasesByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)

	assert.Equal(t, "25-34", cases.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 150},
		&sciensano.CasesEntry{Count: 120},
	}, cases.Groups[0].Values)

	assert.Equal(t, "85+", cases.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.CasesEntry{Count: 100},
		&sciensano.CasesEntry{Count: 0},
	}, cases.Groups[1].Values)

	apiClient.AssertExpectations(t)
}

func TestClient_GetCases_Failure(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetCases", mock.Anything).Return(nil, fmt.Errorf("API error"))

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()

	_, err := client.GetCases(ctx)
	require.Error(t, err)

	_, err = client.GetCasesByRegion(ctx)
	require.Error(t, err)

	_, err = client.GetCasesByProvince(ctx)
	require.Error(t, err)

	_, err = client.GetCasesByAgeGroup(ctx)
	require.Error(t, err)

}

func TestClient_Cases_ApplyRegions(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testCasesResponse, nil)

	cases, err := client.GetCasesByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)

	cases.ApplyRange(time.Time{}, time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC))
	require.Len(t, cases.Timestamps, 1)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 1)
	require.Len(t, cases.Groups[1].Values, 1)

	cases, err = client.GetCasesByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)
}

func BenchmarkClient_GetCasesByRegion(b *testing.B) {
	var bigResponse apiclient.APICasesResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &apiclient.APICasesResponseEntry{
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

	apiClient.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	for i := 0; i < 100; i++ {
		_, err := client.GetCasesByRegion(ctx)
		require.NoError(b, err)
	}
}
