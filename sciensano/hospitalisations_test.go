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
	testHospitalisationsResponse = []apiclient.Measurement{
		&apiclient.APIHospitalisationsResponseEntry{
			TimeStamp:   apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     100,
			TotalInICU:  50,
			TotalInResp: 25,
			TotalInECMO: 10,
		},
		&apiclient.APIHospitalisationsResponseEntry{
			TimeStamp:   apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Brussels",
			Province:    "Brussels",
			TotalIn:     10,
			TotalInICU:  5,
			TotalInResp: 3,
			TotalInECMO: 1,
		},
		&apiclient.APIHospitalisationsResponseEntry{
			TimeStamp:   apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     50,
			TotalInICU:  25,
			TotalInResp: 12,
			TotalInECMO: 5,
		},
	}
)

func TestClient_GetHospitalisations(t *testing.T) {
	getter := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisations(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 1)
	assert.Empty(t, entries.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.HospitalisationsEntry{
			In:     110,
			InICU:  55,
			InResp: 28,
			InECMO: 11,
		},
		&sciensano.HospitalisationsEntry{
			In:     50,
			InICU:  25,
			InResp: 12,
			InECMO: 5,
		},
	}, entries.Groups[0].Values)
}

func TestClient_GetHospitalisationsByProvince(t *testing.T) {
	getter := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByProvince(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)

	assert.Equal(t, "Brussels", entries.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.HospitalisationsEntry{
			In:     10,
			InICU:  5,
			InResp: 3,
			InECMO: 1,
		},
		&sciensano.HospitalisationsEntry{},
	}, entries.Groups[0].Values)

	assert.Equal(t, "VlaamsBrabant", entries.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.HospitalisationsEntry{
			In:     100,
			InICU:  50,
			InResp: 25,
			InECMO: 10,
		},
		&sciensano.HospitalisationsEntry{
			In:     50,
			InICU:  25,
			InResp: 12,
			InECMO: 5,
		},
	}, entries.Groups[1].Values)

	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisationsByRegion(t *testing.T) {
	getter := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)

	assert.Equal(t, "Brussels", entries.Groups[0].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.HospitalisationsEntry{
			In:     10,
			InICU:  5,
			InResp: 3,
			InECMO: 1,
		},
		&sciensano.HospitalisationsEntry{},
	}, entries.Groups[0].Values)

	assert.Equal(t, "Flanders", entries.Groups[1].Name)
	assert.Equal(t, []datasets.Copyable{
		&sciensano.HospitalisationsEntry{
			In:     100,
			InICU:  50,
			InResp: 25,
			InECMO: 10,
		},
		&sciensano.HospitalisationsEntry{
			In:     50,
			InICU:  25,
			InResp: 12,
			InECMO: 5,
		},
	}, entries.Groups[1].Values)

	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisations_Failure(t *testing.T) {
	getter := &mocks.Getter{}
	getter.On("GetHospitalisations", mock.Anything).Return(nil, fmt.Errorf("API error"))

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter

	ctx := context.Background()

	_, err := client.GetHospitalisations(ctx)
	require.Error(t, err)

	_, err = client.GetHospitalisationsByRegion(ctx)
	require.Error(t, err)

	_, err = client.GetHospitalisationsByProvince(ctx)
	require.Error(t, err)

	getter.AssertExpectations(t)
}

func TestClient_GetHospitalisations_ApplyRegions(t *testing.T) {
	getter := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testHospitalisationsResponse, nil)

	entries, err := client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 2)
	require.Len(t, entries.Groups[1].Values, 2)

	entries.ApplyRange(time.Time{}, time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC))
	require.Len(t, entries.Timestamps, 1)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 1)
	require.Len(t, entries.Groups[1].Values, 1)

	entries, err = client.GetHospitalisationsByRegion(ctx)
	require.NoError(t, err)
	require.Len(t, entries.Timestamps, 2)
	require.Len(t, entries.Groups, 2)
	require.Len(t, entries.Groups[0].Values, 2)
	require.Len(t, entries.Groups[1].Values, 2)

	getter.AssertExpectations(t)
}

func BenchmarkClient_GetHospitalisationsByRegion(b *testing.B) {
	var bigResponse apiclient.APIHospitalisationsResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels"} {
			bigResponse = append(bigResponse, &apiclient.APIHospitalisationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: ts},
				Region:    region,
				Province:  region,
				TotalIn:   i,
			})
		}
		ts = ts.Add(24 * time.Hour)
	}
	getter := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	ctx := context.Background()

	getter.
		On("GetHospitalisations", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil).Once()

	b.ResetTimer()
	for i := 0; i < 100; i++ {
		_, err := client.GetHospitalisationsByRegion(ctx)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
