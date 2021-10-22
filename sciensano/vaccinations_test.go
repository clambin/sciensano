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
	"reflect"
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

func TestClient_GetVaccinations(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinations(context.Background(), lastDay)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, lastDay, result[0].Timestamp)
	assert.Equal(t, 2, result[0].Partial)
	assert.Equal(t, 1, result[0].Full)
	assert.Equal(t, 4, result[0].SingleDose)
	assert.Equal(t, 5, result[0].Booster)

}

func TestClient_GetVaccinationsByAge(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByAge(context.Background(), lastDay)
	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Contains(t, result, "25-34")
	require.Contains(t, result, "35-44")
	require.Len(t, result["25-34"], 1)
	assert.Equal(t, 1, result["25-34"][0].Partial)
	assert.Equal(t, 0, result["25-34"][0].Full)
	assert.Equal(t, 0, result["25-34"][0].SingleDose)
	assert.Equal(t, 0, result["25-34"][0].Booster)
	require.Len(t, result["35-44"], 1)
	assert.Equal(t, 1, result["35-44"][0].Partial)
	assert.Equal(t, 1, result["35-44"][0].Full)
	assert.Equal(t, 4, result["35-44"][0].SingleDose)
	assert.Equal(t, 5, result["35-44"][0].Booster)

	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(vaccinationsResponse, nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	result, err := client.GetVaccinationsByRegion(context.Background(), lastDay)
	require.NoError(t, err)
	require.Len(t, result, 2)

	require.Contains(t, result, "Flanders")
	require.Len(t, result["Flanders"], 1)
	assert.Equal(t, 2, result["Flanders"][0].Partial)
	assert.Equal(t, 0, result["Flanders"][0].Full)
	assert.Equal(t, 4, result["Flanders"][0].SingleDose)
	assert.Equal(t, 0, result["Flanders"][0].Booster)

	require.Contains(t, result, "Brussels")
	require.Len(t, result["Brussels"], 1)
	assert.Equal(t, 0, result["Brussels"][0].Partial)
	assert.Equal(t, 1, result["Brussels"][0].Full)
	assert.Equal(t, 0, result["Brussels"][0].SingleDose)
	assert.Equal(t, 5, result["Brussels"][0].Booster)

	mock.AssertExpectationsForObjects(t, apiClient)
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

func BenchmarkClient_GetVaccinationsByAge(b *testing.B) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(buildBigVaccinationResponse(), nil)

	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = apiClient

	ctx := context.Background()
	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByAge(ctx, lastDay)
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
		entries []sciensano.Vaccination
	}
	tests := []struct {
		name       string
		args       args
		wantTotals []sciensano.Vaccination
	}{
		{
			name:       "empty",
			args:       args{},
			wantTotals: []sciensano.Vaccination{},
		},
		{
			name: "one",
			args: args{entries: []sciensano.Vaccination{
				{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
			}},
			wantTotals: []sciensano.Vaccination{
				{Partial: 4, Full: 3, SingleDose: 2, Booster: 1},
			},
		},
		{
			name: "many",
			args: args{entries: []sciensano.Vaccination{
				{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
				{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
				{Partial: 2, Full: 1, SingleDose: 1, Booster: 0},
				{Partial: 3, Full: 2, SingleDose: 0, Booster: 1},
				{Partial: 4, Full: 3, SingleDose: 1, Booster: 0},
			}},
			wantTotals: []sciensano.Vaccination{
				{Partial: 0, Full: 0, SingleDose: 0, Booster: 0},
				{Partial: 1, Full: 0, SingleDose: 0, Booster: 1},
				{Partial: 3, Full: 1, SingleDose: 1, Booster: 1},
				{Partial: 6, Full: 3, SingleDose: 1, Booster: 2},
				{Partial: 10, Full: 6, SingleDose: 2, Booster: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTotals := sciensano.AccumulateVaccinations(tt.args.entries); !reflect.DeepEqual(gotTotals, tt.wantTotals) {
				t.Errorf("AccumulateVaccinations() = %v, want %v", gotTotals, tt.wantTotals)
			}
		})
	}
}

func TestVaccination_Total(t *testing.T) {
	vaccination := sciensano.Vaccination{
		Partial:    1,
		Full:       2,
		SingleDose: 3,
		Booster:    4,
	}

	assert.Equal(t, 10, vaccination.Total())
}
