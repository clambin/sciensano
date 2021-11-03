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
	timestamp = time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)

	testVaccinationsResponse = apiclient.APIVaccinationsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Dose:      "A",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "A",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "C",
			Count:     4,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "B",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
	}
)

func TestVaccinationsEntry_Total(t *testing.T) {
	entry := sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       2,
		SingleDose: 3,
		Booster:    4,
	}
	assert.Equal(t, 10, entry.Total())
}

func TestClient_GetVaccinations(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinations(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
	}, result.Timestamps)

	require.Len(t, result.Groups, 1)

	assert.Empty(t, result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    2,
		Full:       1,
		SingleDose: 4,
		Booster:    5,
	}, *result.Groups[0].Values[0].(*sciensano.VaccinationsEntry))
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       0,
		SingleDose: 0,
		Booster:    5,
	}, *result.Groups[0].Values[1].(*sciensano.VaccinationsEntry))
}

func TestClient_GetVaccinationsByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByAgeGroup(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
	}, result.Timestamps)

	require.Len(t, result.Groups, 2)

	assert.Equal(t, "25-34", result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       0,
		SingleDose: 0,
		Booster:    0,
	}, *result.Groups[0].Values[0].(*sciensano.VaccinationsEntry))
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       0,
		SingleDose: 0,
		Booster:    0,
	}, *result.Groups[0].Values[1].(*sciensano.VaccinationsEntry))

	assert.Equal(t, "35-44", result.Groups[1].Name)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       1,
		SingleDose: 4,
		Booster:    5,
	}, *result.Groups[1].Values[0].(*sciensano.VaccinationsEntry))
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       0,
		SingleDose: 0,
		Booster:    5,
	}, *result.Groups[1].Values[1].(*sciensano.VaccinationsEntry))
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByRegion(context.Background())
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
	}, result.Timestamps)

	require.Len(t, result.Groups, 2)

	assert.Equal(t, "Brussels", result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       1,
		SingleDose: 0,
		Booster:    5,
	}, *result.Groups[0].Values[0].(*sciensano.VaccinationsEntry))
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       0,
		SingleDose: 0,
		Booster:    5,
	}, *result.Groups[0].Values[1].(*sciensano.VaccinationsEntry))

	assert.Equal(t, "Flanders", result.Groups[1].Name)
	require.Len(t, result.Groups[0].Values, 2)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    2,
		Full:       0,
		SingleDose: 4,
		Booster:    0,
	}, *result.Groups[1].Values[0].(*sciensano.VaccinationsEntry))
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       0,
		SingleDose: 0,
		Booster:    0,
	}, *result.Groups[1].Values[1].(*sciensano.VaccinationsEntry))
}

func TestClient_GetVaccinations_Failure(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(nil, fmt.Errorf("API error"))

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()

	_, err := client.GetVaccinations(ctx)
	require.Error(t, err)

	_, err = client.GetVaccinationsByRegion(ctx)
	require.Error(t, err)

	_, err = client.GetVaccinationsByAgeGroup(ctx)
	require.Error(t, err)

}

func TestClient_Vaccinations_ApplyRegions(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient
	ctx := context.Background()

	apiClient.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testVaccinationsResponse, nil)

	cases, err := client.GetVaccinationsByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)

	cases.ApplyRange(time.Time{}, timestamp)
	require.Len(t, cases.Timestamps, 1)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 1)
	require.Len(t, cases.Groups[1].Values, 1)

	cases, err = client.GetVaccinationsByAgeGroup(ctx)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Len(t, cases.Groups, 2)
	require.Len(t, cases.Groups[0].Values, 2)
	require.Len(t, cases.Groups[1].Values, 2)
}

var bigVaccinationResponse apiclient.APIVaccinationsResponse

func buildBigVaccinationResponse() {
	bigVaccinationResponse = apiclient.APIVaccinationsResponse{}

	startDate := time.Now().Add(-365 * 24 * time.Hour)
	for i := 0; i < 365; i++ {
		for _, region := range []string{"Flanders", "Brussels", "Wallonia"} {
			for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"} {
				bigVaccinationResponse = append(bigVaccinationResponse,
					apiclient.APIVaccinationsResponseEntry{
						TimeStamp: apiclient.TimeStamp{Time: startDate},
						Region:    region,
						AgeGroup:  ageGroup,
						Dose:      "A",
						Count:     i * 2,
					},
					apiclient.APIVaccinationsResponseEntry{
						TimeStamp: apiclient.TimeStamp{Time: startDate},
						Region:    region,
						AgeGroup:  ageGroup,
						Dose:      "B",
						Count:     i,
					})
			}
		}
		startDate = startDate.Add(24 * time.Hour)
	}
	fmt.Printf("response had %d entries\n", len(bigVaccinationResponse))
}

func BenchmarkClient_GetVaccination(b *testing.B) {
	buildBigVaccinationResponse()
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinations(ctx)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func BenchmarkClient_GetVaccinationsByAgeGroup(b *testing.B) {
	buildBigVaccinationResponse()
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByAgeGroup(ctx)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	buildBigVaccinationResponse()
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByRegion(ctx)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func TestAccumulateVaccinations(t *testing.T) {
	type args struct {
		entries *datasets.Dataset
	}
	tests := []struct {
		name       string
		args       args
		wantTotals *datasets.Dataset
	}{
		{
			name: "one",
			args: args{entries: &datasets.Dataset{Groups: []datasets.GroupedDatasetEntry{
				{Values: []datasets.Copyable{
					&sciensano.VaccinationsEntry{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
				}},
			}}},
			wantTotals: &datasets.Dataset{Groups: []datasets.GroupedDatasetEntry{
				{Values: []datasets.Copyable{
					&sciensano.VaccinationsEntry{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
				}},
			}},
		},
		{
			name: "many",
			args: args{entries: &datasets.Dataset{Groups: []datasets.GroupedDatasetEntry{
				{Values: []datasets.Copyable{
					&sciensano.VaccinationsEntry{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
					&sciensano.VaccinationsEntry{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
					&sciensano.VaccinationsEntry{Partial: 2, Full: 1, SingleDose: 1, Booster: 0},
					&sciensano.VaccinationsEntry{Partial: 3, Full: 2, SingleDose: 0, Booster: 1},
					&sciensano.VaccinationsEntry{Partial: 4, Full: 3, SingleDose: 1, Booster: 0},
				}},
			}}},
			wantTotals: &datasets.Dataset{Groups: []datasets.GroupedDatasetEntry{
				{Values: []datasets.Copyable{
					&sciensano.VaccinationsEntry{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
					&sciensano.VaccinationsEntry{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
					&sciensano.VaccinationsEntry{Partial: 3, Full: 1, SingleDose: 1, Booster: 1},
					&sciensano.VaccinationsEntry{Partial: 6, Full: 3, SingleDose: 1, Booster: 2},
					&sciensano.VaccinationsEntry{Partial: 10, Full: 6, SingleDose: 2, Booster: 2},
				}},
			}},
		},
	}

	for _, tt := range tests {
		sciensano.AccumulateVaccinations(tt.args.entries)

		assert.Equal(t, tt.args.entries, tt.wantTotals, tt.name)
	}
}

func TestVaccinationsEntry_GetValue(t *testing.T) {
	input := sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       2,
		SingleDose: 3,
		Booster:    4,
	}

	assert.Equal(t, 1, input.GetValue(sciensano.VaccinationTypePartial))
	assert.Equal(t, 5, input.GetValue(sciensano.VaccinationTypeFull))
	assert.Equal(t, 4, input.GetValue(sciensano.VaccinationTypeBooster))
}
