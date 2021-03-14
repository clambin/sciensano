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
	close(testsResponse)

	vaccinationResponse := make(chan []sciensano.Vaccination)
	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "",
		Value:    "",
		Response: vaccinationResponse,
	}
	vaccinations := <-vaccinationResponse
	assert.Len(t, vaccinations, 6)
	assert.Equal(t, testDate, tests[5].Timestamp)
	assert.Equal(t, 5, tests[5].Positive)
	assert.Equal(t, 10, tests[5].Total)

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "AgeGroup",
		Value:    "45-54",
		Response: vaccinationResponse,
	}
	vaccinations = <-vaccinationResponse
	assert.Len(t, vaccinations, 6)
	assert.Equal(t, testDate, tests[5].Timestamp)
	assert.Equal(t, 5, tests[5].Positive)
	assert.Equal(t, 10, tests[5].Total)

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "Region",
		Value:    "Flanders",
		Response: vaccinationResponse,
	}
	vaccinations = <-vaccinationResponse
	assert.Len(t, vaccinations, 6)
	assert.Equal(t, testDate, tests[5].Timestamp)
	assert.Equal(t, 5, tests[5].Positive)
	assert.Equal(t, 10, tests[5].Total)

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "Region",
		Value:    "Brussels",
		Response: vaccinationResponse,
	}
	vaccinations = <-vaccinationResponse
	assert.Len(t, vaccinations, 0)

	c.Vaccinations <- cache.VaccinationsRequest{
		EndTime:  testDate,
		Filter:   "Invalid",
		Value:    "",
		Response: vaccinationResponse,
	}
	vaccinations = <-vaccinationResponse
	assert.Len(t, vaccinations, 0)

	close(vaccinationResponse)

	c.Stop()

	time.Sleep(500 * time.Millisecond)
}
