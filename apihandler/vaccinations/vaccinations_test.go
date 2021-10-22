package vaccinations_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	mockAPI "github.com/clambin/sciensano/apiclient/mocks"
	vaccinationsHandler "github.com/clambin/sciensano/apihandler/vaccinations"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Search(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, nil)

	targets := h.Search()
	assert.Equal(t, []string{
		"vacc-age-booster",
		"vacc-age-full",
		"vacc-age-partial",
		"vacc-age-rate-booster",
		"vacc-age-rate-full",
		"vacc-age-rate-partial",
		"vacc-region-booster",
		"vacc-region-full",
		"vacc-region-partial",
		"vacc-region-rate-booster",
		"vacc-region-rate-full",
		"vacc-region-rate-partial",
		"vaccination-lag",
		"vaccinations",
	}, targets)
}

func TestHandler_TableQuery_Vaccinations(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, nil)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	getter.
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
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Dose:      "C",
					Count:     10,
				},
				{
					TimeStamp: apiclient.TimeStamp{Time: endDate},
					Dose:      "E",
					Count:     5,
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
	response, err := h.Endpoints().TableQuery(context.Background(), "vaccinations", request)
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
				assert.Equal(t, 35.0, data[len(data)-1])
			case "booster":
				require.Len(t, data, 1)
				assert.Equal(t, 5.0, data[len(data)-1])
			default:
				assert.Fail(t, "unexpected column", column.Text)
			}
		}
	}

	getter.AssertExpectations(t)
}

func TestHandler_TableQuery_GroupedVaccination_ByAge(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, nil)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	getter.
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
	// TODO: partial, booster calls
	response, err := h.Endpoints().TableQuery(context.Background(), "vacc-age-full", request)
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

	getter.AssertExpectations(t)
}

func TestHandler_TableQuery_GroupedRatedVaccination_ByAge(t *testing.T) {
	getter := &mockAPI.Getter{}
	demographics := &mockDemographics.Demographics{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, demographics)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	demographics.
		On("GetAgeGroupFigures").
		Return(map[string]int{
			"45-54": 1000,
		}, nil)

	getter.
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
	response, err := h.Endpoints().TableQuery(context.Background(), "vacc-age-rate-partial", request)
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

	mock.AssertExpectationsForObjects(t, demographics, getter)
}

func TestHandler_TableQuery_GroupedVaccination_ByRegion(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, nil)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	getter.
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
	response, err := h.Endpoints().TableQuery(context.Background(), "vacc-region-full", request)
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

	getter.AssertExpectations(t)
}

func TestHandler_TableQuery_GroupedRatedVaccination_ByRegion(t *testing.T) {
	getter := &mockAPI.Getter{}
	demographics := &mockDemographics.Demographics{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, demographics)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	demographics.
		On("GetRegionFigures").
		Return(map[string]int{
			"Flanders": 6000,
		})

	getter.
		On("GetVaccinations", mock.Anything).
		Return(
			[]*apiclient.APIVaccinationsResponse{
				{TimeStamp: apiclient.TimeStamp{Time: endDate}, Region: "Flanders", Dose: "A", Count: 100},
				{TimeStamp: apiclient.TimeStamp{Time: endDate}, Region: "Flanders", Dose: "B", Count: 25},
				{TimeStamp: apiclient.TimeStamp{Time: endDate}, Region: "Ostbelgien", Dose: "A", Count: 25},
				{TimeStamp: apiclient.TimeStamp{Time: endDate}, Region: "Ostbelgien", Dose: "B", Count: 5},
				{TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)}, Region: "Flanders", Dose: "A", Count: 174},
				{TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)}, Region: "Flanders", Dose: "B", Count: 50},
				{TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)}, Region: "Ostbelgien", Dose: "A", Count: 25},
				{TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)}, Region: "Ostbelgien", Dose: "B", Count: 5},
			}, nil)

	// Vaccination rate grouped by Region
	response, err := h.Endpoints().TableQuery(context.Background(), "vacc-region-rate-full", request)
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
			case "Ostbelgien":
				require.NotZero(t, len(data))
				assert.Equal(t, 64, int(1000000*data[len(data)-1]))
			}
		}
	}

	mock.AssertExpectationsForObjects(t, getter, demographics)
}

func TestHandler_TableQuery_Vaccination_Lag(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := vaccinationsHandler.New(client, nil)

	endDate := time.Date(2021, 03, 31, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	getter.
		On("GetVaccinations", mock.Anything).
		Return([]*apiclient.APIVaccinationsResponse{
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
	response, err := h.Endpoints().TableQuery(context.Background(), "vaccination-lag", request)
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

	getter.AssertExpectations(t)
}
