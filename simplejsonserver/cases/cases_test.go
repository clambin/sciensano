package cases_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/cases"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestCase struct {
	Scope    cases.Scope
	Response *query.TableResponse
}

var (
	testResponse = []apiclient.APIResponse{
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "85+",
			Cases:     100,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			AgeGroup:  "25-34",
			Cases:     150,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "25-34",
			Cases:     120,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			Province:  "",
			AgeGroup:  "",
			Cases:     5,
		},
		&sciensano.APICasesResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			AgeGroup:  "65-74",
			Cases:     100,
		},
	}

	testCases = []TestCase{
		{
			Scope: cases.ScopeAll,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "total", Data: query.NumberColumn{250.0, 125.0}},
				},
			},
		},
		{
			Scope: cases.ScopeProvince,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: query.NumberColumn{0.0, 5.0}},
					{Text: "Brussels", Data: query.NumberColumn{150.0, 0.0}},
					{Text: "VlaamsBrabant", Data: query.NumberColumn{100.0, 120.0}},
				},
			},
		},
		{
			Scope: cases.ScopeRegion,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: query.NumberColumn{0.0, 5.0}},
					{Text: "Brussels", Data: query.NumberColumn{150.0, 0.0}},
					{Text: "Flanders", Data: query.NumberColumn{100.0, 120.0}},
				},
			},
		},
		{
			Scope: cases.ScopeAge,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: query.NumberColumn{0.0, 5.0}},
					{Text: "25-34", Data: query.NumberColumn{150.0, 120.0}},
					{Text: "65-74", Data: query.NumberColumn{0.0, 0.0}},
					{Text: "85+", Data: query.NumberColumn{100.0, 0.0}},
				},
			},
		},
	}
)

func TestHandler_TableQuery(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(testResponse, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Cases.APIClient = f

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}}

	for index, testCase := range testCases {
		h := cases.Handler{
			Reporter: r,
			Scope:    testCase.Scope,
		}
		response, err := h.Endpoints().Query(context.Background(), req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.Response, response, index)
	}

	mock.AssertExpectationsForObjects(t, f)
}

func TestHandler_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Cases.APIClient = f

	req := query.Request{}

	h := cases.Handler{
		Reporter: r,
		Scope:    cases.ScopeAll,
	}
	_, err := h.Endpoints().Query(context.Background(), req)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkCasesHandler(b *testing.B) {
	var bigResponse []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APICasesResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				Cases:     i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeCases).Return(bigResponse, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Cases.APIClient = f
	h := cases.Handler{
		Reporter: r,
		Scope:    cases.ScopeRegion,
	}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
