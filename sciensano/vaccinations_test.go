package sciensano_test

import (
	"context"
	"fmt"
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
	lastDay = time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)

	vaccinationsResponse = []*apiclient.APIVaccinationsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Dose:      "A",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "A",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "C",
			Count:     4,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "B",
			Count:     1,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: lastDay.Add(24 * time.Hour)},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "B",
			Count:     1,
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
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinations(context.Background(), lastDay)
	require.NoError(t, err)

	assert.Equal(t, []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)}, result.Timestamps)

	require.Len(t, result.Groups, 1)

	assert.Empty(t, result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 1)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    2,
		Full:       1,
		SingleDose: 4,
		Booster:    5,
	}, *result.Groups[0].Values[0])
}

func TestClient_GetVaccinationsByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByAgeGroup(context.Background(), lastDay)
	require.NoError(t, err)

	assert.Equal(t, []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)}, result.Timestamps)

	require.Len(t, result.Groups, 2)

	assert.Equal(t, "25-34", result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 1)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       0,
		SingleDose: 0,
		Booster:    0,
	}, *result.Groups[0].Values[0])

	assert.Equal(t, "35-44", result.Groups[1].Name)
	require.Len(t, result.Groups[0].Values, 1)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    1,
		Full:       1,
		SingleDose: 4,
		Booster:    5,
	}, *result.Groups[1].Values[0])
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByRegion(context.Background(), lastDay)
	require.NoError(t, err)

	assert.Equal(t, []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)}, result.Timestamps)

	require.Len(t, result.Groups, 2)

	assert.Equal(t, "Brussels", result.Groups[0].Name)
	require.Len(t, result.Groups[0].Values, 1)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    0,
		Full:       1,
		SingleDose: 0,
		Booster:    5,
	}, *result.Groups[0].Values[0])

	assert.Equal(t, "Flanders", result.Groups[1].Name)
	require.Len(t, result.Groups[0].Values, 1)
	assert.Equal(t, sciensano.VaccinationsEntry{
		Partial:    2,
		Full:       0,
		SingleDose: 4,
		Booster:    0,
	}, *result.Groups[1].Values[0])
}

var bigVaccinationResponse []*apiclient.APIVaccinationsResponse

func buildBigVaccinationResponse() []*apiclient.APIVaccinationsResponse {
	if bigVaccinationResponse == nil {
		startDate := time.Now().Add(-365 * 24 * time.Hour)
		for i := 0; i < 365; i++ {
			for _, region := range []string{"Flanders", "Brussels", "Wallonia"} {
				for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"} {
					bigVaccinationResponse = append(bigVaccinationResponse,
						&apiclient.APIVaccinationsResponse{
							TimeStamp: apiclient.TimeStamp{Time: startDate},
							Region:    region,
							AgeGroup:  ageGroup,
							Dose:      "A",
							Count:     i * 2,
						},
						&apiclient.APIVaccinationsResponse{
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

	return bigVaccinationResponse
}

func BenchmarkClient_GetVaccinationsByAgeGroup(b *testing.B) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(buildBigVaccinationResponse(), nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByAgeGroup(ctx, lastDay)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(buildBigVaccinationResponse(), nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByRegion(ctx, lastDay)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func TestAccumulateVaccinations(t *testing.T) {
	type args struct {
		entries *sciensano.Vaccinations
	}
	tests := []struct {
		name       string
		args       args
		wantTotals *sciensano.Vaccinations
	}{
		{
			name: "one",
			args: args{entries: &sciensano.Vaccinations{Groups: []sciensano.GroupedVaccinationsEntry{
				{Values: []*sciensano.VaccinationsEntry{
					{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
				}},
			}}},
			wantTotals: &sciensano.Vaccinations{Groups: []sciensano.GroupedVaccinationsEntry{
				{Values: []*sciensano.VaccinationsEntry{
					{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
				}},
			}},
		},
		{
			name: "many",
			args: args{entries: &sciensano.Vaccinations{Groups: []sciensano.GroupedVaccinationsEntry{
				{Values: []*sciensano.VaccinationsEntry{
					{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
					{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
					{Partial: 2, Full: 1, SingleDose: 1, Booster: 0},
					{Partial: 3, Full: 2, SingleDose: 0, Booster: 1},
					{Partial: 4, Full: 3, SingleDose: 1, Booster: 0},
				}},
			}}},
			wantTotals: &sciensano.Vaccinations{Groups: []sciensano.GroupedVaccinationsEntry{
				{Values: []*sciensano.VaccinationsEntry{
					{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
					{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
					{Partial: 3, Full: 1, SingleDose: 1, Booster: 1},
					{Partial: 6, Full: 3, SingleDose: 1, Booster: 2},
					{Partial: 10, Full: 6, SingleDose: 2, Booster: 2},
				}},
			}},
		},
	}

	for _, tt := range tests {
		sciensano.AccumulateVaccinations(tt.args.entries)

		assert.Equal(t, tt.args.entries, tt.wantTotals, tt.name)
	}
}
