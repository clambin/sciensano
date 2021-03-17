package cache_test

import (
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCache_Run(t *testing.T) {
	c := cache.New(5 * time.Minute)
	c.API = &mockapi.API{
		Tests:        mockapi.DefaultTests,
		Vaccinations: mockapi.DefaultVaccinations,
	}

	go c.Run()

	testDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	testsResponse := make(chan []sciensano.Test)
	c.Tests <- cache.TestsRequest{
		EndTime:  testDate,
		Response: testsResponse,
	}
	tests := <-testsResponse
	assert.Len(t, tests, 6)
	assert.Equal(t, testDate, tests[5].Timestamp)
	assert.Equal(t, 5, tests[5].Positive)
	assert.Equal(t, 10, tests[5].Total)
	close(testsResponse)

	vaccinationResponse := make(chan []sciensano.Vaccination)
	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "",
		Response: vaccinationResponse,
	}
	vaccinations := <-vaccinationResponse
	if assert.Len(t, vaccinations, 6) {
		assert.Equal(t, testDate, vaccinations[5].Timestamp)
		assert.Equal(t, 5, vaccinations[5].FirstDose)
		assert.Equal(t, 4, vaccinations[5].SecondDose)
	}

	groupedVaccinationResponse := make(chan map[string][]sciensano.Vaccination)
	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:         testDate,
		Filter:          "AgeGroup",
		GroupedResponse: groupedVaccinationResponse,
	}
	groupedVaccinations := <-groupedVaccinationResponse
	if assert.Len(t, groupedVaccinations, 1) && assert.Len(t, groupedVaccinations["45-54"], 6) {
		assert.Equal(t, testDate, groupedVaccinations["45-54"][5].Timestamp)
		assert.Equal(t, 5, groupedVaccinations["45-54"][5].FirstDose)
		assert.Equal(t, 4, groupedVaccinations["45-54"][5].SecondDose)
	}

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:         testDate,
		Filter:          "Region",
		GroupedResponse: groupedVaccinationResponse,
	}
	groupedVaccinations = <-groupedVaccinationResponse
	if assert.Len(t, groupedVaccinations, 1) && assert.Len(t, groupedVaccinations["Flanders"], 6) {
		assert.Equal(t, testDate, groupedVaccinations["Flanders"][5].Timestamp)
		assert.Equal(t, 5, groupedVaccinations["Flanders"][5].FirstDose)
		assert.Equal(t, 4, groupedVaccinations["Flanders"][5].SecondDose)
	}

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:         testDate,
		Filter:          "invalid",
		GroupedResponse: groupedVaccinationResponse,
	}
	groupedVaccinations = <-groupedVaccinationResponse
	assert.Len(t, groupedVaccinations, 0)

	c.Stop()

	time.Sleep(500 * time.Millisecond)
}
