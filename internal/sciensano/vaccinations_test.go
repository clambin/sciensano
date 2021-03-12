package sciensano_test

import (
	"github.com/clambin/sciensano/internal/sciensano"
	"github.com/stretchr/testify/assert"
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
		totals sciensano.Vaccinations
		result sciensano.Vaccinations
	)
	client := sciensano.Client{VaccinationsCacheDuration: 1 * time.Hour}
	testDate := time.Now()
	totals, err = client.GetVaccinations(testDate)
	assert.Nil(t, err)

	if assert.Greater(t, len(totals), 0) {
		var firstDose, secondDose int
		for _, ageGroup := range sciensano.AgeGroups {
			result, err = client.GetVaccinationsByAge(testDate, ageGroup)
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

func TestGetAgeGroupFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
	}{
		{args: args{target: "vaccinations-18-25-first"}, wantOutput: "18-25"},
		{args: args{target: "vaccinations-85+-second"}, wantOutput: "85+"},
		{args: args{target: "vaccinations-85+-second"}, wantOutput: "85+"},
		{args: args{target: "vaccinations--second"}, wantOutput: ""},
		{args: args{target: "invalid"}, wantOutput: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := sciensano.GetAgeGroupFromTarget(tt.args.target); gotOutput != tt.wantOutput {
				t.Errorf("GetAgeGroupFromTarget() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestGetModeFromTarget(t *testing.T) {
	type args struct {
		target string
	}
	tests := []struct {
		name     string
		args     args
		wantMode string
	}{
		{name: "A", args: args{target: "vaccinations-0-16-first"}, wantMode: "A"},
		{name: "B", args: args{target: "vaccinations--second"}, wantMode: "B"},
		{name: "invalid", args: args{target: "foobar"}, wantMode: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMode := sciensano.GetModeFromTarget(tt.args.target); gotMode != tt.wantMode {
				t.Errorf("GetModeFromTarget() = %v, want %v", gotMode, tt.wantMode)
			}
		})
	}
}
