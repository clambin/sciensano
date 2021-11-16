package cases_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	casesHandler "github.com/clambin/sciensano/apihandler/cases"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestCase struct {
	Target   string
	Response *grafanajson.TableQueryResponse
}

var (
	testResponse = []measurement.Measurement{
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			Province:  "",
			AgeGroup:  "",
			Cases:     5,
		},
		&sciensano.APICasesResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "65-74",
			Cases:     100,
		},
	}

	testCases = []TestCase{
		{
			Target: "cases",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "total", Data: grafanajson.TableQueryResponseNumberColumn{250.0, 125.0}},
				},
			},
		},
		{
			Target: "cases-province",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 5.0}},
					{Text: "Brussels", Data: grafanajson.TableQueryResponseNumberColumn{150.0, 0.0}},
					{Text: "VlaamsBrabant", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 120.0}},
				},
			},
		},
		{
			Target: "cases-region",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 5.0}},
					{Text: "Brussels", Data: grafanajson.TableQueryResponseNumberColumn{150.0, 0.0}},
					{Text: "Flanders", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 120.0}},
				},
			},
		},
		{
			Target: "cases-age",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 5.0}},
					{Text: "25-34", Data: grafanajson.TableQueryResponseNumberColumn{150.0, 120.0}},
					{Text: "65-74", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 0.0}},
					{Text: "85+", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 0.0}},
				},
			},
		},
	}
)

func TestHandler_Search(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = getter
	h := casesHandler.New(client)

	targets := h.Search()
	assert.Equal(t, []string{
		"cases",
		"cases-age",
		"cases-province",
		"cases-region",
	}, targets)
}

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = getter
	h := casesHandler.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Cases").
		Return(testResponse, true)

	for _, testCase := range testCases {
		response, err := h.Endpoints().TableQuery(context.Background(), testCase.Target, args)
		require.NoError(t, err, testCase.Target)
		assert.Equal(t, testCase.Response, response, testCase.Target)
	}

	getter.AssertExpectations(t)
}

func BenchmarkHandler_TableQuery(b *testing.B) {
	var bigResponse []measurement.Measurement
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APICasesResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				Cases:     i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = getter
	h := casesHandler.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Cases").
		Return(bigResponse, true)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "cases-region", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
