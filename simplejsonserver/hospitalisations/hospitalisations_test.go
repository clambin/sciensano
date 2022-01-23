package hospitalisations_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/hospitalisations"
	"github.com/clambin/simplejson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var (
	testResponse = []measurement.Measurement{
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     100,
			TotalInICU:  10,
			TotalInResp: 5,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:      "Brussels",
			Province:    "Brussels",
			TotalIn:     50,
			TotalInICU:  1,
			TotalInResp: 0,
			TotalInECMO: 0,
		},
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "Flanders",
			Province:    "VlaamsBrabant",
			TotalIn:     90,
			TotalInICU:  10,
			TotalInResp: 5,
			TotalInECMO: 1,
		},
		&sciensano.APIHospitalisationsResponseEntry{
			TimeStamp:   sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:      "",
			Province:    "",
			TotalIn:     1,
			TotalInICU:  1,
			TotalInResp: 0,
			TotalInECMO: 0,
		},
		&sciensano.APIHospitalisationsResponseEntry{
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
		Response *simplejson.TableQueryResponse
	}{
		{
			Scope: hospitalisations.ScopeAll,
			Response: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "in", Data: simplejson.TableQueryResponseNumberColumn{150.0, 91.0}},
					{Text: "inICU", Data: simplejson.TableQueryResponseNumberColumn{11.0, 11.0}},
					{Text: "inResp", Data: simplejson.TableQueryResponseNumberColumn{5.0, 5.0}},
					{Text: "inECMO", Data: simplejson.TableQueryResponseNumberColumn{1.0, 1.0}},
				},
			},
		},
		{
			Scope: hospitalisations.ScopeProvince,
			Response: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: simplejson.TableQueryResponseNumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: simplejson.TableQueryResponseNumberColumn{50.0, 0.0}},
					{Text: "VlaamsBrabant", Data: simplejson.TableQueryResponseNumberColumn{100.0, 90.0}},
				},
			},
		},
		{
			Scope: hospitalisations.ScopeRegion,
			Response: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: simplejson.TableQueryResponseNumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: simplejson.TableQueryResponseNumberColumn{50.0, 0.0}},
					{Text: "Flanders", Data: simplejson.TableQueryResponseNumberColumn{100.0, 90.0}},
				},
			},
		},
	}
)

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = getter

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Hospitalisations").
		Return(testResponse, true)

	for index, testCase := range testCases {
		h := hospitalisations.Handler{
			Reporter: client,
			Scope:    testCase.Scope,
		}
		response, err := h.Endpoints().TableQuery(context.Background(), args)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.Response, response, index)
	}

	getter.AssertExpectations(t)
}

func TestHandler_Failure(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = getter

	args := &simplejson.TableQueryArgs{}

	getter.
		On("Get", "Hospitalisations").
		Return(nil, false)

	h := hospitalisations.Handler{
		Reporter: client,
		Scope:    hospitalisations.ScopeAll,
	}
	_, err := h.Endpoints().TableQuery(context.Background(), args)
	assert.Error(t, err)

	getter.AssertExpectations(t)
}

func BenchmarkHandler_TableQuery(b *testing.B) {
	var bigResponse []measurement.Measurement
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIHospitalisationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Province:  region,
				Region:    region,
				TotalIn:   i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = getter
	h := hospitalisations.Handler{
		Reporter: client,
		Scope:    hospitalisations.ScopeRegion,
	}

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Hospitalisations").
		Return(bigResponse, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), args)
		if err != nil {
			b.Fatal(err)
		}
	}
}
