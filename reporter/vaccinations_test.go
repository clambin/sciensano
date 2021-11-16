package reporter_test

import (
	"fmt"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	timestamp = time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)

	testVaccinationsResponse = []measurement.Measurement{
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Dose:      "A",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "A",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "C",
			Count:     4,
		},
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "B",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
		&sciensano.APIVaccinationsResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
	}
)

func TestClient_GetVaccinations(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(testVaccinationsResponse, true)

	r := reporter.New(time.Hour)
	r.Sciensano = cache

	result, err := r.GetVaccinations()
	require.NoError(t, err)

	assert.Equal(t, &datasets.Dataset{
		Timestamps: []time.Time{
			time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
			time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
		},
		Groups: []datasets.DatasetGroup{
			{Name: "partial", Values: []float64{2, 0}},
			{Name: "full", Values: []float64{1, 0}},
			{Name: "singledose", Values: []float64{4, 0}},
			{Name: "booster", Values: []float64{5, 5}},
		}}, result)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetVaccinationsByAgeGroup(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(testVaccinationsResponse, true)

	client := reporter.New(time.Hour)
	client.Sciensano = cache

	testCases := []struct {
		mode   int
		output *datasets.Dataset
	}{
		{
			mode: reporter.VaccinationTypePartial,
			output: &datasets.Dataset{
				Timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
				Groups: []datasets.DatasetGroup{
					{Name: "25-34", Values: []float64{1}},
					{Name: "35-44", Values: []float64{1}},
				},
			},
		},
		{
			mode: reporter.VaccinationTypeFull,
			output: &datasets.Dataset{
				Timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
				Groups: []datasets.DatasetGroup{
					{Name: "35-44", Values: []float64{5}},
				},
			},
		},
		{
			mode: reporter.VaccinationTypeBooster,
			output: &datasets.Dataset{
				Timestamps: []time.Time{
					time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
				},
				Groups: []datasets.DatasetGroup{
					{Name: "35-44", Values: []float64{5, 5}},
				},
			},
		},
	}

	for index, testCase := range testCases {
		result, err := client.GetVaccinationsByAgeGroup(testCase.mode)
		require.NoError(t, err, index)

		assert.Equal(t, testCase.output, result, index)
	}
	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(testVaccinationsResponse, true)

	client := reporter.New(time.Hour)
	client.Sciensano = cache

	testCases := []struct {
		mode   int
		output *datasets.Dataset
	}{
		{
			mode: reporter.VaccinationTypePartial,
			output: &datasets.Dataset{
				Timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
				Groups: []datasets.DatasetGroup{
					{Name: "Flanders", Values: []float64{2}},
				},
			},
		},
		{
			mode: reporter.VaccinationTypeFull,
			output: &datasets.Dataset{
				Timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
				Groups: []datasets.DatasetGroup{
					{Name: "Brussels", Values: []float64{1}},
					{Name: "Flanders", Values: []float64{4}},
				},
			},
		},
		{
			mode: reporter.VaccinationTypeBooster,
			output: &datasets.Dataset{
				Timestamps: []time.Time{
					time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
					time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
				},
				Groups: []datasets.DatasetGroup{
					{Name: "Brussels", Values: []float64{5, 5}},
				},
			},
		},
	}

	for index, testCase := range testCases {
		result, err := client.GetVaccinationsByRegion(testCase.mode)
		require.NoError(t, err, index)

		assert.Equal(t, testCase.output, result, index)
	}
	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetVaccinations_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(nil, false)

	client := reporter.New(time.Hour)
	client.Sciensano = cache

	_, err := client.GetVaccinations()
	require.Error(t, err)

	_, err = client.GetVaccinationsByRegion(reporter.VaccinationTypePartial)
	require.Error(t, err)

	_, err = client.GetVaccinationsByAgeGroup(reporter.VaccinationTypePartial)
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_Vaccinations_ApplyRegions(t *testing.T) {
	cache := &mocks.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = cache

	cache.
		On("Get", "Vaccinations").
		Return(testVaccinationsResponse, true)

	cases, err := client.GetVaccinationsByAgeGroup(reporter.VaccinationTypeBooster)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Equal(t, []datasets.DatasetGroup{{Name: "35-44", Values: []float64{5, 5}}}, cases.Groups)

	cases.ApplyRange(time.Time{}, timestamp)
	require.Len(t, cases.Timestamps, 1)
	require.Equal(t, []datasets.DatasetGroup{{Name: "35-44", Values: []float64{5}}}, cases.Groups)

	mock.AssertExpectationsForObjects(t, cache)
}

var bigVaccinationResponse []measurement.Measurement

func buildBigVaccinationResponse() {
	bigVaccinationResponse = []measurement.Measurement{}

	startDate := time.Now().Add(-365 * 24 * time.Hour)
	for i := 0; i < 365; i++ {
		for _, region := range []string{"Flanders", "Brussels", "Wallonia"} {
			for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"} {
				bigVaccinationResponse = append(bigVaccinationResponse,
					&sciensano.APIVaccinationsResponseEntry{
						TimeStamp: sciensano.TimeStamp{Time: startDate},
						Region:    region,
						AgeGroup:  ageGroup,
						Dose:      "A",
						Count:     i * 2,
					},
					&sciensano.APIVaccinationsResponseEntry{
						TimeStamp: sciensano.TimeStamp{Time: startDate},
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
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.Sciensano = cache

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := client.GetVaccinations()
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, cache)
}

func BenchmarkClient_GetVaccinationsByAgeGroup(b *testing.B) {
	buildBigVaccinationResponse()
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.Sciensano = cache

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := client.GetVaccinationsByAgeGroup(reporter.VaccinationTypeFull)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, cache)
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	buildBigVaccinationResponse()
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.Sciensano = cache

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := client.GetVaccinationsByRegion(reporter.VaccinationTypeFull)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, cache)
}
