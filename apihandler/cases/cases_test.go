package cases_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	mockAPI "github.com/clambin/sciensano/apiclient/mocks"
	casesHandler "github.com/clambin/sciensano/apihandler/cases"
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestCase struct {
	Target   string
	Response *grafanajson.TableQueryResponse
}

var (
	testResponse = []*apiclient.APICasesResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			Province:  "",
			AgeGroup:  "",
			Cases:     5,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
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
					{Text: "Timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "cases", Data: grafanajson.TableQueryResponseNumberColumn{250.0, 125.0}},
				},
			},
		},
		{
			Target: "cases-province",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "Timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
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
					{Text: "Timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
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
					{Text: "Timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 5.0}},
					{Text: "25-34", Data: grafanajson.TableQueryResponseNumberColumn{150.0, 120.0}},
					{Text: "85+", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 0.0}},
				},
			},
		},
	}
)

func TestHandler_Search(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := casesHandler.New(getter, client)

	targets := h.Search()
	assert.Equal(t, []string{
		"cases",
		"cases-age",
		"cases-province",
		"cases-region",
	}, targets)
}

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := casesHandler.New(getter, client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

	for _, testCase := range testCases {
		response, err := h.Endpoints().TableQuery(context.Background(), testCase.Target, args)
		require.NoError(t, err, testCase.Target)
		assert.Equal(t, testCase.Response, response)
	}

	getter.AssertExpectations(t)
}

func BenchmarkHandler_TableQuery(b *testing.B) {
	var bigResponse []*apiclient.APICasesResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &apiclient.APICasesResponse{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				Cases:     i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := casesHandler.New(getter, client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "cases-region", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}

func BenchmarkHandler_TableQuery_Alt(b *testing.B) {
	var bigResponse []*apiclient.APICasesResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &apiclient.APICasesResponse{
				TimeStamp: apiclient.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				Cases:     i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockAPI.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := casesHandler.New(getter, client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetCases", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "cases-region-alt", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
