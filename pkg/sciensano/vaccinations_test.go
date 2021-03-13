package sciensano_test

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestAPIClient_GetVaccinations(t *testing.T) {
	client := sciensano.Client{}
	firstDay := time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC)
	result, err := client.GetVaccinations(firstDay)

	if assert.Nil(t, err) {
		if assert.Len(t, result, 1) {
			assert.Equal(t, firstDay, result[0].Timestamp)
			assert.Equal(t, 298, result[0].FirstDose)
			assert.Equal(t, 0, result[0].SecondDose)
		}
	}
}

func TestAPIClient_GetVaccinationsByAge(t *testing.T) {
	var (
		err    error
		totals []sciensano.Vaccination
		result []sciensano.Vaccination
	)
	client := sciensano.Client{VaccinationsCacheDuration: 1 * time.Hour}
	testDate := time.Now()
	totals, err = client.GetVaccinations(testDate)
	assert.Nil(t, err)

	if assert.Greater(t, len(totals), 0) {
		vaccinationsByAge := make(map[string][]sciensano.Vaccination, 0)

		for _, ageGroup := range sciensano.AgeGroups {
			result, err = client.GetVaccinationsByAge(testDate, ageGroup)
			if assert.Nil(t, err) {
				vaccinationsByAge[ageGroup] = result
			}
		}

		var firstDose, secondDose int
		for _, vaccinations := range vaccinationsByAge {
			if len(vaccinations) > 0 {
				firstDose += vaccinations[len(vaccinations)-1].FirstDose
				secondDose += vaccinations[len(vaccinations)-1].SecondDose
			}
		}

		// the sum of final dose count for each age group should be the same as the final overall dose count
		assert.Equal(t, totals[len(totals)-1].FirstDose, firstDose)
		assert.Equal(t, totals[len(totals)-1].SecondDose, secondDose)
	}
}

func TestAPIClient_GetVaccinationsByRegion(t *testing.T) {
	var (
		err    error
		totals sciensano.Vaccinations
		result sciensano.Vaccinations
	)
	client := sciensano.Client{VaccinationsCacheDuration: 1 * time.Hour}
	testDate := time.Now()
	totals, err = client.GetVaccinations(testDate)
	assert.Nil(t, err)

	if assert.Greater(t, len(totals), 0) {
		var firstDose, secondDose int
		for _, region := range sciensano.Regions {
			result, err = client.GetVaccinationsByRegion(testDate, region)
			if assert.Nil(t, err) && len(result) > 0 {
				firstDose += result[len(result)-1].FirstDose
				secondDose += result[len(result)-1].SecondDose
			}
		}
		// the sum of final dose count for each age group should be the same as the final overall dose count
		assert.Equal(t, totals[len(totals)-1].FirstDose, firstDose)
		assert.Equal(t, totals[len(totals)-1].SecondDose, secondDose)
	}
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	var (
		err    error
		totals sciensano.Vaccinations
		result sciensano.Vaccinations
	)
	client := sciensano.Client{VaccinationsCacheDuration: 1 * time.Hour}
	testDate := time.Now()
	totals, err = client.GetVaccinations(testDate)
	assert.Nil(b, err)

	if assert.Greater(b, len(totals), 0) {
		var firstDose, secondDose int
		for _, region := range sciensano.Regions {
			result, err = client.GetVaccinationsByRegion(testDate, region)
			if assert.Nil(b, err) && len(result) > 0 {
				firstDose += result[len(result)-1].FirstDose
				secondDose += result[len(result)-1].SecondDose
			}
		}
		// the sum of final dose count for each age group should be the same as the final overall dose count
		assert.Equal(b, totals[len(totals)-1].FirstDose, firstDose)
		assert.Equal(b, totals[len(totals)-1].SecondDose, secondDose)
	}
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
				{FirstDose: 4, SecondDose: 3},
			}},
			wantTotals: []sciensano.Vaccination{
				{FirstDose: 4, SecondDose: 3},
			},
		},
		{
			name: "many",
			args: args{entries: []sciensano.Vaccination{
				{FirstDose: 0, SecondDose: 0},
				{FirstDose: 1, SecondDose: 0},
				{FirstDose: 2, SecondDose: 1},
				{FirstDose: 3, SecondDose: 2},
				{FirstDose: 4, SecondDose: 3},
			}},
			wantTotals: []sciensano.Vaccination{
				{FirstDose: 0, SecondDose: 0},
				{FirstDose: 1, SecondDose: 0},
				{FirstDose: 3, SecondDose: 1},
				{FirstDose: 6, SecondDose: 3},
				{FirstDose: 10, SecondDose: 6},
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