package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAPIHandler_Vaccinations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Dose:      "A",
					Count:     100,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Dose:      "B",
					Count:     25,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Dose:      "A",
					Count:     74,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Dose:      "B",
					Count:     50,
				},
			}, nil)

	// Vaccinations
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vaccinations", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			assert.Equal(t, "timestamp", column.Text)
			assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "partial":
				require.Len(t, data, 1)
				assert.Equal(t, 100.0, data[len(data)-1])
			case "full":
				require.Len(t, data, 1)
				assert.Equal(t, 25.0, data[len(data)-1])
			default:
				assert.Fail(t, "unexpected column", column.Text)
			}
		}
	}
}

func TestAPIHandler_VaccinationsByAge(t *testing.T) {
	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					AgeGroup:  "45-54",
					Dose:      "A",
					Count:     100,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					AgeGroup:  "45-54",
					Dose:      "B",
					Count:     25,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					AgeGroup:  "45-54",
					Dose:      "A",
					Count:     74,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					AgeGroup:  "45-54",
					Dose:      "B",
					Count:     50,
				},
			}, nil)

	// Vaccinations grouped by Age
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vacc-age-full", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			require.Equal(t, "timestamp", column.Text)
			require.NotZero(t, len(data))
			assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "45-54":
				require.NotZero(t, len(data))
				assert.Equal(t, 25.0, data[len(data)-1])
			}
		}
	}

	mock.AssertExpectationsForObjects(t, &stack.demoClient)
}

func TestAPIHandler_VaccinationByAge_Rate(t *testing.T) {
	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.demoClient.
		On("GetAgeGroupFigures").
		Return(map[string]int{
			"45-54": 1000,
		})

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					AgeGroup:  "45-54",
					Dose:      "A",
					Count:     100,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					AgeGroup:  "45-54",
					Dose:      "B",
					Count:     25,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					AgeGroup:  "45-54",
					Dose:      "A",
					Count:     74,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					AgeGroup:  "45-54",
					Dose:      "B",
					Count:     50,
				},
			}, nil)

	// Vaccination rate grouped by Age
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vacc-age-rate-partial", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			assert.Equal(t, "timestamp", column.Text)
			if assert.NotZero(t, len(data)) {
				assert.Equal(t, endDate, data[len(data)-1])
			}
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "45-54":
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, 10, int(100*data[len(data)-1]))
				}
			}
		}
	}

	mock.AssertExpectationsForObjects(t, &stack.demoClient, &stack.sciensanoClient)

}

func TestAPIHandler_VaccinationByRegion(t *testing.T) {
	endDate := time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Region:    "Flanders",
					Dose:      "A",
					Count:     100,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Region:    "Flanders",
					Dose:      "B",
					Count:     25,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Region:    "Flanders",
					Dose:      "A",
					Count:     74,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Region:    "Flanders",
					Dose:      "B",
					Count:     50,
				},
			}, nil)

	// Vaccinations grouped by Region
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vacc-region-full", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			require.Equal(t, "timestamp", column.Text)
			require.NotZero(t, len(data))
			assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "Flanders":
				require.NotZero(t, len(data))
				assert.Equal(t, 25.0, data[len(data)-1])
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion_Rate(t *testing.T) {
	endDate := time.Date(2021, 03, 11, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.demoClient.
		On("GetRegionFigures").
		Return(map[string]int{
			"Flanders": 6000,
		})

	stack.sciensanoClient.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Region:    "Flanders",
					Dose:      "A",
					Count:     100,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Region:    "Flanders",
					Dose:      "B",
					Count:     25,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Region:    "Flanders",
					Dose:      "A",
					Count:     74,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
					Region:    "Flanders",
					Dose:      "B",
					Count:     50,
				},
			}, nil)

	// Vaccination rate grouped by Region
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vacc-region-rate-full", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			require.Equal(t, "timestamp", column.Text)
			require.NotZero(t, len(data))
			assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "Flanders":
				require.NotZero(t, len(data))
				assert.Equal(t, 4166, int(1000000*data[len(data)-1]))
			}
		}
	}
}

func TestAPIHandler_Vaccination_Lag(t *testing.T) {
	endDate := time.Date(2021, 03, 31, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.sciensanoClient.On("GetVaccinations", mock.Anything).Return([]*apiclient.APIVaccinationsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     200,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     50,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     300,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     320,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     155,
		},
	}, nil)

	// Lag
	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vaccination-lag", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		// case grafanaJson.TableQueryResponseTimeColumn:
		//	require.Equal(t, "timestamp", column.Text)
		//	require.NotZero(t, len(data))
		//	assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "lag":
				require.NotZero(t, len(data))
				assert.Equal(t, 15.0, data[len(data)-1])
			}
		}
	}
}
