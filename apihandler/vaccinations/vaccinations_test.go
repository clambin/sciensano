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

type TestCase struct {
	Target   string
	Response *grafanaJson.TableQueryResponse
}

var (
	timestamp = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	vaccinationTestData = apiclient.APIVaccinationsResponse{
		{TimeStamp: apiclient.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "25-34", Dose: "C", Count: 1},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "35-44", Dose: "E", Count: 1},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "35-44", Dose: "B", Count: 2},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "25-34", Dose: "A", Count: 2},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp}, Region: "", AgeGroup: "", Dose: "A", Count: 0},

		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "B", Count: 3},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "C", Count: 4},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "A", Count: 5},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "25-34", Dose: "E", Count: 5},

		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "A", Count: 9},
		{TimeStamp: apiclient.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "E", Count: 9},
	}

	testCases = []TestCase{
		{
			Target: "vaccinations",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "partial", Data: grafanaJson.TableQueryResponseNumberColumn{2, 7}},
					{Text: "full", Data: grafanaJson.TableQueryResponseNumberColumn{3, 10}},
					{Text: "booster", Data: grafanaJson.TableQueryResponseNumberColumn{1, 6}},
				},
			},
		},
		{
			Target: "vacc-region-partial",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{2, 7}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
				},
			},
		},
		{
			Target: "vacc-region-rate-partial",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{0.2, 0.7}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
				},
			},
		},
		{
			Target: "vacc-region-full",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{2, 6}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{1, 4}},
				},
			},
		},
		{
			Target: "vacc-region-rate-full",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{0.2, 0.6}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{0.01, 0.04}},
				},
			},
		},
		{
			Target: "vacc-region-booster",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{0, 5}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{1, 1}},
				},
			},
		},
		{
			Target: "vacc-region-rate-booster",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0.5}},
					{Text: "Flanders", Data: grafanaJson.TableQueryResponseNumberColumn{0.01, 0.01}},
				},
			},
		},
		{
			Target: "vacc-age-partial",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{2, 2}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{0, 5}},
				},
			},
		},
		{
			Target: "vacc-age-rate-partial",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{0.02, 0.02}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0.5}},
				},
			},
		},
		{
			Target: "vacc-age-full",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{1, 4}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{2, 6}},
				},
			},
		},
		{
			Target: "vacc-age-rate-full",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{0.01, 0.04}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{0.2, 0.6}},
				},
			},
		},
		{
			Target: "vacc-age-booster",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{0, 5}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{1, 1}},
				},
			},
		},
		{
			Target: "vacc-age-rate-booster",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0.05}},
					{Text: "35-44", Data: grafanaJson.TableQueryResponseNumberColumn{0.1, 0.1}},
				},
			},
		},
		{
			Target: "vaccination-lag",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "lag", Data: grafanaJson.TableQueryResponseNumberColumn{0, 0}},
				},
			},
		},
	}
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

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	demo := &mockDemographics.Demographics{}
	h := vaccinationsHandler.New(client, demo)

	endDate := timestamp.Add(24 * time.Hour)
	args := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	getter.
		On("GetVaccinations", mock.Anything).
		Return(vaccinationTestData, nil)

	demo.
		On("GetRegionFigures").
		Return(map[string]int{
			"Flanders": 100,
			"Brussels": 10,
		})

	demo.
		On("GetAgeGroupFigures").
		Return(map[string]int{
			"25-34": 100,
			"35-44": 10,
		})

	for _, testCase := range testCases {
		response, err := h.Endpoints().TableQuery(context.Background(), testCase.Target, args)
		require.NoError(t, err, testCase.Target)
		assert.Equal(t, testCase.Response, response, testCase.Target)
	}

	mock.AssertExpectationsForObjects(t, getter, demo)
}

func BenchmarkHandler_TableQuery(b *testing.B) {
	var bigResponse apiclient.APIVaccinationsResponse
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "A",
				Count:     i + 100,
			})
			bigResponse = append(bigResponse, apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "B",
				Count:     i,
			})
		}

		ts = ts.Add(24 * time.Hour)
	}

	getter := &mockAPI.Getter{}
	client := sciensano.NewCachedClient(time.Hour)
	client.Getter = getter
	h := vaccinationsHandler.New(client, nil)

	args := &grafanaJson.TableQueryArgs{CommonQueryArgs: grafanaJson.CommonQueryArgs{Range: grafanaJson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "vacc-region-full", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
