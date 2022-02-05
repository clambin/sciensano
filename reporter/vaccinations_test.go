package reporter_test

import (
	"github.com/clambin/sciensano/apiclient"
	"time"
)

import (
	"github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	testVaccinationsResponse = []apiclient.APIResponse{
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "E",
			Count:     5,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Dose:      "A",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "A",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "35-44",
			Dose:      "C",
			Count:     4,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "35-44",
			Dose:      "B",
			Count:     1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
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
	r.APICache = cache

	entries, err := r.GetVaccinations()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"booster", "full", "partial", "singledose"}, entries.GetColumns())

	values, ok := entries.GetValues("partial")
	require.True(t, ok)
	assert.Equal(t, []float64{2, 0}, values)

	values, ok = entries.GetValues("full")
	require.True(t, ok)
	assert.Equal(t, []float64{1, 0}, values)

	values, ok = entries.GetValues("singledose")
	require.True(t, ok)
	assert.Equal(t, []float64{4, 0}, values)

	values, ok = entries.GetValues("booster")
	require.True(t, ok)
	assert.Equal(t, []float64{5, 5}, values)

	mock.AssertExpectationsForObjects(t, cache)
}

type vaccinationsTestCase struct {
	mode       reporter.VaccinationType
	timestamps []time.Time
	values     map[string][]float64
}

func TestClient_GetVaccinationsByAgeGroup(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(testVaccinationsResponse, true)

	client := reporter.New(time.Hour)
	client.APICache = cache

	testCases := []vaccinationsTestCase{
		{
			mode:       reporter.VaccinationTypePartial,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"25-34": {1},
				"35-44": {1},
			},
		},
		{
			mode:       reporter.VaccinationTypeFull,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"35-44": {5},
			},
		},
		{
			mode: reporter.VaccinationTypeBooster,
			timestamps: []time.Time{
				time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			},
			values: map[string][]float64{
				"35-44": {5, 5},
			},
		},
	}

	for index, testCase := range testCases {
		result, err := client.GetVaccinationsByAgeGroup(testCase.mode)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.timestamps, result.GetTimestamps())
		assert.Len(t, result.GetColumns(), len(testCase.values))
		for column, expected := range testCase.values {
			values, ok := result.GetValues(column)
			require.True(t, ok)
			assert.Equal(t, expected, values)
		}
	}
	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(testVaccinationsResponse, true)

	client := reporter.New(time.Hour)
	client.APICache = cache

	testCases := []vaccinationsTestCase{
		{
			mode:       reporter.VaccinationTypePartial,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"Flanders": {2},
			},
		},
		{
			mode:       reporter.VaccinationTypeFull,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"Brussels": {1},
				"Flanders": {4},
			},
		},
		{
			mode: reporter.VaccinationTypeBooster,
			timestamps: []time.Time{
				time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			},
			values: map[string][]float64{
				"Brussels": {5, 5},
			},
		},
	}

	for index, testCase := range testCases {
		result, err := client.GetVaccinationsByRegion(testCase.mode)
		require.NoError(t, err, index)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.timestamps, result.GetTimestamps())
		assert.Len(t, result.GetColumns(), len(testCase.values))
		for column, expected := range testCase.values {
			values, ok := result.GetValues(column)
			require.True(t, ok)
			assert.Equal(t, expected, values)
		}
	}
	mock.AssertExpectationsForObjects(t, cache)
}

func TestClient_GetVaccinations_Failure(t *testing.T) {
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(nil, false)

	client := reporter.New(time.Hour)
	client.APICache = cache

	_, err := client.GetVaccinations()
	require.Error(t, err)

	_, err = client.GetVaccinationsByRegion(reporter.VaccinationTypePartial)
	require.Error(t, err)

	_, err = client.GetVaccinationsByAgeGroup(reporter.VaccinationTypePartial)
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache)
}

var bigVaccinationResponse []apiclient.APIResponse

func buildBigVaccinationResponse() {
	bigVaccinationResponse = []apiclient.APIResponse{}

	startDate := time.Now().Add(-365 * 24 * time.Hour)
	for i := 0; i < 365; i++ {
		for _, region := range []string{"Flanders", "Brussels", "Wallonia"} {
			for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "85+"} {
				bigVaccinationResponse = append(bigVaccinationResponse,
					&sciensano.APIVaccinationsResponse{
						TimeStamp: sciensano.TimeStamp{Time: startDate},
						Region:    region,
						AgeGroup:  ageGroup,
						Dose:      "A",
						Count:     i * 2,
					},
					&sciensano.APIVaccinationsResponse{
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
	// fmt.Printf("responder had %d entries\n", len(bigVaccinationResponse))
}

func BenchmarkClient_GetVaccination(b *testing.B) {
	buildBigVaccinationResponse()
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.APICache = cache

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.GetVaccinations()
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkClient_GetVaccinationsByAgeGroup(b *testing.B) {
	buildBigVaccinationResponse()
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.APICache = cache

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.GetVaccinationsByAgeGroup(reporter.VaccinationTypeFull)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	buildBigVaccinationResponse()
	cache := &mocks.Holder{}
	cache.On("Get", "Vaccinations").Return(bigVaccinationResponse, true)

	client := reporter.New(0)
	client.APICache = cache

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.GetVaccinationsByRegion(reporter.VaccinationTypeFull)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
