package vaccinations_test

import (
	"errors"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/vaccinations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

var (
	testVaccinationsResponse = []apiclient.APIResponse{
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "Pfizer-BioNTech",
			Region:       "Flanders",
			AgeGroup:     "25-34",
			Dose:         "A",
			Count:        1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "AstraZeneca-Oxford",
			Region:       "Flanders",
			AgeGroup:     "35-44",
			Dose:         "A",
			Count:        1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "Johnson&Johnson",
			Region:       "Flanders",
			AgeGroup:     "35-44",
			Dose:         "C",
			Count:        4,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "Moderna",
			Region:       "Brussels",
			AgeGroup:     "35-44",
			Dose:         "B",
			Count:        1,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 10, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "Moderna",
			Region:       "Brussels",
			AgeGroup:     "35-44",
			Dose:         "E",
			Count:        5,
		},
		&sciensano.APIVaccinationsResponse{
			TimeStamp:    sciensano.TimeStamp{Time: time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)},
			Manufacturer: "Moderna",
			Region:       "Brussels",
			AgeGroup:     "35-44",
			Dose:         "E",
			Count:        5,
		},
	}
)

func TestReporter_Get(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(testVaccinationsResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	entries, err := r.Get()
	require.NoError(t, err)

	assert.Equal(t, []time.Time{
		time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
	}, entries.GetTimestamps())

	assert.Equal(t, []string{"time", "partial", "full", "singledose", "booster"}, entries.GetColumns())

	values, ok := entries.GetFloatValues("partial")
	require.True(t, ok)
	assert.Equal(t, []float64{2, 0}, values)

	values, ok = entries.GetFloatValues("full")
	require.True(t, ok)
	assert.Equal(t, []float64{1, 0}, values)

	values, ok = entries.GetFloatValues("singledose")
	require.True(t, ok)
	assert.Equal(t, []float64{4, 0}, values)

	values, ok = entries.GetFloatValues("booster")
	require.True(t, ok)
	assert.Equal(t, []float64{5, 5}, values)

	mock.AssertExpectationsForObjects(t, f)
}

type vaccinationsTestCase struct {
	mode       int
	timestamps []time.Time
	values     map[string][]float64
}

func TestReporter_GetByAgeGroup(t *testing.T) {
	testCases := []vaccinationsTestCase{
		{
			mode:       vaccinations.TypePartial,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"25-34": {1},
				"35-44": {1},
			},
		},
		{
			mode:       vaccinations.TypeFull,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"35-44": {5},
			},
		},
		{
			mode: vaccinations.TypeBooster,
			timestamps: []time.Time{
				time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			},
			values: map[string][]float64{
				"35-44": {5, 5},
			},
		},
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(testVaccinationsResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	for index, testCase := range testCases {
		result, err := r.GetByAgeGroup(testCase.mode)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.timestamps, result.GetTimestamps())
		assert.Len(t, result.GetColumns(), 1+len(testCase.values))
		for column, expected := range testCase.values {
			values, ok := result.GetFloatValues(column)
			require.True(t, ok)
			assert.Equal(t, expected, values)
		}
	}
	mock.AssertExpectationsForObjects(t, f)
}

func TestReporter_GetByRegion(t *testing.T) {
	testCases := []vaccinationsTestCase{
		{
			mode:       vaccinations.TypePartial,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"Flanders": {2},
			},
		},
		{
			mode:       vaccinations.TypeFull,
			timestamps: []time.Time{time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC)},
			values: map[string][]float64{
				"Brussels": {1},
				"Flanders": {4},
			},
		},
		{
			mode: vaccinations.TypeBooster,
			timestamps: []time.Time{
				time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			},
			values: map[string][]float64{
				"Brussels": {5, 5},
			},
		},
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(testVaccinationsResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	for index, testCase := range testCases {
		result, err := r.GetByRegion(testCase.mode)
		require.NoError(t, err, index)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.timestamps, result.GetTimestamps())
		assert.Len(t, result.GetColumns(), 1+len(testCase.values))
		for column, expected := range testCase.values {
			values, ok := result.GetFloatValues(column)
			require.True(t, ok)
			assert.Equal(t, expected, values)
		}
	}
	mock.AssertExpectationsForObjects(t, f)
}

func TestReporter_GetByManufacturer(t *testing.T) {
	testCases := []vaccinationsTestCase{
		{
			timestamps: []time.Time{
				time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.March, 11, 0, 0, 0, 0, time.UTC),
			},
			values: map[string][]float64{
				"AstraZeneca-Oxford": {1, 0},
				"Pfizer-BioNTech":    {1, 0},
				"Moderna":            {6, 5},
				"Johnson&Johnson":    {4, 0},
			},
		},
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(testVaccinationsResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	for index, testCase := range testCases {
		result, err := r.GetByManufacturer()
		require.NoError(t, err, index)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.timestamps, result.GetTimestamps(), index)
		assert.Len(t, result.GetColumns(), 1+len(testCase.values), index)
		for column, expected := range testCase.values {
			values, ok := result.GetFloatValues(column)
			require.True(t, ok, column+"-"+strconv.Itoa(index))
			assert.Equal(t, expected, values, column+"-"+strconv.Itoa(index))
		}
	}
	mock.AssertExpectationsForObjects(t, f)
}

func TestReporter_Get_Failures(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(nil, errors.New("fail"))

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	_, err := r.Get()
	assert.Error(t, err)

	_, err = r.GetByRegion(vaccinations.TypePartial)
	assert.Error(t, err)

	_, err = r.GetByAgeGroup(vaccinations.TypePartial)
	assert.Error(t, err)

	_, err = r.GetByManufacturer()
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
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
}

func BenchmarkVaccinations_Get(b *testing.B) {
	buildBigVaccinationResponse()

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(bigVaccinationResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.Get()
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkVaccinations_GetByAgeGroup(b *testing.B) {
	buildBigVaccinationResponse()
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(bigVaccinationResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetByAgeGroup(vaccinations.TypeFull)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}

func BenchmarkVaccinations_GetByRegion(b *testing.B) {
	buildBigVaccinationResponse()
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(bigVaccinationResponse, nil)

	r := vaccinations.Reporter{
		ReportCache: cache.NewCache(time.Hour),
		APIClient:   f,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := r.GetByRegion(vaccinations.TypeFull)
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
}
