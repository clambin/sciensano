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
