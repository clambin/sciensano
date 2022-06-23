package mortality_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/cases"
	"github.com/clambin/sciensano/simplejsonserver/mortality"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestCase struct {
	Scope    mortality.Scope
	Response *query.TableResponse
}

var (
	testResponse = []apiclient.APIResponse{
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Deaths:    100,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "25-34",
			Deaths:    150,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Deaths:    120,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			AgeGroup:  "",
			Deaths:    5,
		},
		&sciensano.APIMortalityResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "65-74",
			Deaths:    100,
		},
	}

	testCases = []TestCase{
		{
			Scope: mortality.ScopeAll,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "total", Data: query.NumberColumn{250.0, 125.0}},
				},
			},
		},
		{
			Scope: mortality.ScopeRegion,
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
			Scope: mortality.ScopeAge,
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
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(testResponse, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Mortality.APIClient = f

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}}

	for index, testCase := range testCases {
		h := mortality.Handler{
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
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Mortality.APIClient = f

	req := query.Request{}

	h := mortality.Handler{
		Reporter: r,
		Scope:    cases.ScopeAll,
	}
	_, err := h.Endpoints().Query(context.Background(), req)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}

func BenchmarkMortalityHandler(b *testing.B) {
	var bigResponse []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIMortalityResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Region:    region,
				Deaths:    i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeMortality).Return(bigResponse, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Mortality.APIClient = f
	h := mortality.Handler{
		Reporter: r,
		Scope:    mortality.ScopeRegion,
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
