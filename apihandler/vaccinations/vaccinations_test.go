package vaccinations_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	vaccinationsHandler "github.com/clambin/sciensano/apihandler/vaccinations"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	timestamp = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	vaccinationTestData = []measurement.Measurement{
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "25-34", Dose: "C", Count: 1},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "35-44", Dose: "E", Count: 1},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "35-44", Dose: "B", Count: 2},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "25-34", Dose: "A", Count: 2},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "", AgeGroup: "", Dose: "A", Count: 1},

		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "B", Count: 3},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "C", Count: 4},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "A", Count: 5},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "25-34", Dose: "E", Count: 5},

		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "A", Count: 9},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "E", Count: 9},
	}

	testCases = []struct {
		Target   string
		Response *grafanaJson.TableQueryResponse
	}{
		{
			Target: "vaccinations",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "partial", Data: grafanaJson.TableQueryResponseNumberColumn{3, 8}},
					{Text: "full", Data: grafanaJson.TableQueryResponseNumberColumn{3, 10}},
					{Text: "booster", Data: grafanaJson.TableQueryResponseNumberColumn{1, 6}},
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
		{
			Target: "vacc-region-partial",
			Response: &grafanaJson.TableQueryResponse{
				Columns: []grafanaJson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanaJson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{1, 1}},
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
					{Text: "(unknown)", Data: grafanaJson.TableQueryResponseNumberColumn{1, 1}},
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
	}
)

func TestHandler_Search(t *testing.T) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = cache
	demo := &mockDemographics.Demographics{}
	h := vaccinationsHandler.New(client, demo)

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
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = cache
	demo := &mockDemographics.Demographics{}
	h := vaccinationsHandler.New(client, demo)

	endDate := timestamp.Add(24 * time.Hour)
	args := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	cache.
		On("Get", "Vaccinations").
		Return(vaccinationTestData, true)

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

	ctx := context.Background()
	for _, testCase := range testCases {
		response, err := h.Endpoints().TableQuery(ctx, testCase.Target, args)
		require.NoError(t, err, testCase.Target)
		assert.Equal(t, testCase.Response, response, testCase.Target)
	}

	mock.AssertExpectationsForObjects(t, cache, demo)
}

func BenchmarkHandler_TableQuery(b *testing.B) {
	var bigResponse []measurement.Measurement
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "A",
				Count:     i + 100,
			})
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "B",
				Count:     i,
			})
		}

		ts = ts.Add(24 * time.Hour)
	}

	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = cache
	h := vaccinationsHandler.New(client, nil)

	args := &grafanaJson.TableQueryArgs{CommonQueryArgs: grafanaJson.CommonQueryArgs{Range: grafanaJson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	cache.
		On("Get", "Vaccinations").
		Return(bigResponse, true)

	_, err := h.Endpoints().TableQuery(context.Background(), "vacc-region-full", args)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		_, _ = h.Endpoints().TableQuery(context.Background(), "vacc-region-full", args)
	}

	cache.AssertExpectations(b)
}
