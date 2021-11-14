package reporter_test

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/sciensano/mocks"
	"github.com/clambin/sciensano/measurement"
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
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	r := reporter.NewCachedClient(time.Hour)
	r.Sciensano = apiClient

	result, err := r.GetVaccinations(context.Background())
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

	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetVaccinationsByAgeGroup(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient

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
		result, err := client.GetVaccinationsByAgeGroup(context.Background(), testCase.mode)
		require.NoError(t, err, index)

		assert.Equal(t, testCase.output, result, index)
	}
	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetVaccinationsByRegion(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(testVaccinationsResponse, nil)

	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient

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
		result, err := client.GetVaccinationsByRegion(context.Background(), testCase.mode)
		require.NoError(t, err, index)

		assert.Equal(t, testCase.output, result, index)
	}
	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_GetVaccinations_Failure(t *testing.T) {
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(nil, fmt.Errorf("API error"))

	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient

	ctx := context.Background()

	_, err := client.GetVaccinations(ctx)
	require.Error(t, err)

	_, err = client.GetVaccinationsByRegion(ctx, reporter.VaccinationTypePartial)
	require.Error(t, err)

	_, err = client.GetVaccinationsByAgeGroup(ctx, reporter.VaccinationTypePartial)
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, apiClient)
}

func TestClient_Vaccinations_ApplyRegions(t *testing.T) {
	apiClient := &mocks.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = apiClient
	ctx := context.Background()

	apiClient.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return(testVaccinationsResponse, nil)

	cases, err := client.GetVaccinationsByAgeGroup(ctx, reporter.VaccinationTypeBooster)
	require.NoError(t, err)
	require.Len(t, cases.Timestamps, 2)
	require.Equal(t, []datasets.DatasetGroup{{Name: "35-44", Values: []float64{5, 5}}}, cases.Groups)

	cases.ApplyRange(time.Time{}, timestamp)
	require.Len(t, cases.Timestamps, 1)
	require.Equal(t, []datasets.DatasetGroup{{Name: "35-44", Values: []float64{5}}}, cases.Groups)

	mock.AssertExpectationsForObjects(t, apiClient)
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
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := reporter.NewCachedClient(0)
	client.Sciensano = apiClient

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := client.GetVaccinations(ctx)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func BenchmarkClient_GetVaccinationsByAgeGroup(b *testing.B) {
	buildBigVaccinationResponse()
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := reporter.NewCachedClient(0)
	client.Sciensano = apiClient

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := client.GetVaccinationsByAgeGroup(ctx, reporter.VaccinationTypeFull)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}

func BenchmarkClient_GetVaccinationsByRegion(b *testing.B) {
	buildBigVaccinationResponse()
	apiClient := &mocks.Getter{}
	apiClient.On("GetVaccinations", mock.Anything).Return(bigVaccinationResponse, nil)

	client := reporter.NewCachedClient(0)
	client.Sciensano = apiClient

	ctx := context.Background()

	b.ResetTimer()

	for i := 0; i < 10; i++ {
		_, err := client.GetVaccinationsByRegion(ctx, reporter.VaccinationTypeFull)
		require.NoError(b, err)
	}

	mock.AssertExpectationsForObjects(b, apiClient)
}
