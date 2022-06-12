package hospitalisations_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/hospitalisations"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResponse = []apiclient.APIResponse{
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     100,
			TotalInICU:  10,
			TotalInResp: 5,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Brussels",
			Province:    "Brussels",
			TotalIn:     50,
			TotalInICU:  1,
			TotalInResp: 0,
			TotalInECMO: 0,
		},
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     90,
			TotalInICU:  10,
			TotalInResp: 5,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "",
			Province:    "",
			TotalIn:     1,
			TotalInICU:  1,
			TotalInResp: 0,
			TotalInECMO: 0,
		},
		&sciensano.APIHospitalisationsResponse{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     90,
			TotalInICU:  10,
			TotalInResp: 5,
			TotalInECMO: 1,
		},
	}

	testCases = []struct {
		Scope    hospitalisations.Scope
		Response *query.TableResponse
	}{
		{
			Scope: hospitalisations.ScopeAll,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "in", Data: query.NumberColumn{150.0, 91.0}},
					{Text: "inICU", Data: query.NumberColumn{11.0, 11.0}},
					{Text: "inResp", Data: query.NumberColumn{5.0, 5.0}},
					{Text: "inECMO", Data: query.NumberColumn{1.0, 1.0}},
				},
			},
		},
		{
			Scope: hospitalisations.ScopeProvince,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: query.NumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: query.NumberColumn{50.0, 0.0}},
					{Text: "VlaamsBrabant", Data: query.NumberColumn{100.0, 90.0}},
				},
			},
		},
		{
			Scope: hospitalisations.ScopeRegion,
			Response: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: query.NumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: query.NumberColumn{50.0, 0.0}},
					{Text: "Flanders", Data: query.NumberColumn{100.0, 90.0}},
				},
			},
		},
	}
)

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockCache.Holder{}
	getter.
		On("Get", "Hospitalisations").
		Return(testResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Hospitalisations.APICache = getter

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}}

	for index, testCase := range testCases {
		h := hospitalisations.Handler{
			Reporter: r,
			Scope:    testCase.Scope,
		}
		response, err := h.Endpoints().Query(context.Background(), req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.Response, response, index)
	}

	getter.AssertExpectations(t)
}

func TestHandler_Failure(t *testing.T) {
	getter := &mockCache.Holder{}
	getter.
		On("Get", "Hospitalisations").
		Return(nil, false)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Hospitalisations.APICache = getter

	req := query.Request{}

	h := hospitalisations.Handler{
		Reporter: r,
		Scope:    hospitalisations.ScopeAll,
	}
	_, err := h.Endpoints().Query(context.Background(), req)
	assert.Error(t, err)

	getter.AssertExpectations(t)
}

func BenchmarkHospitalisationHandler(b *testing.B) {
	b.StopTimer()
	var bigResponse []apiclient.APIResponse
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIHospitalisationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				TotalIn:   i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockCache.Holder{}
	getter.
		On("Get", "Hospitalisations").
		Return(bigResponse, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Hospitalisations.APICache = getter
	h := hospitalisations.Handler{
		Reporter: r,
		Scope:    hospitalisations.ScopeRegion,
	}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
