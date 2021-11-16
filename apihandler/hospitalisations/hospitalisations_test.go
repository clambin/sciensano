package hospitalisations_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	hospitalisationsHandler "github.com/clambin/sciensano/apihandler/hospitalisations"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
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
		Target   string
		Response *grafanajson.TableQueryResponse
	}{
		{
			Target: "hospitalisations",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "in", Data: grafanajson.TableQueryResponseNumberColumn{150.0, 91.0}},
					{Text: "inICU", Data: grafanajson.TableQueryResponseNumberColumn{11.0, 11.0}},
					{Text: "inResp", Data: grafanajson.TableQueryResponseNumberColumn{5.0, 5.0}},
					{Text: "inECMO", Data: grafanajson.TableQueryResponseNumberColumn{1.0, 1.0}},
				},
			},
		},
		{
			Target: "hospitalisations-province",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: grafanajson.TableQueryResponseNumberColumn{50.0, 0.0}},
					{Text: "VlaamsBrabant", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 90.0}},
				},
			},
		},
		{
			Target: "hospitalisations-region",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "(unknown)", Data: grafanajson.TableQueryResponseNumberColumn{0.0, 1.0}},
					{Text: "Brussels", Data: grafanajson.TableQueryResponseNumberColumn{50.0, 0.0}},
					{Text: "Flanders", Data: grafanajson.TableQueryResponseNumberColumn{100.0, 90.0}},
				},
			},
		},
	}
)

func TestHandler_Search(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = getter
	h := hospitalisationsHandler.New(client)

	targets := h.Search()
	assert.Equal(t, []string{
		"hospitalisations",
		"hospitalisations-province",
		"hospitalisations-region",
	}, targets)
}

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.Sciensano = getter
	h := hospitalisationsHandler.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Hospitalisations").
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
	client.Sciensano = getter
	h := hospitalisationsHandler.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("Get", "Hospitalisations").
		Return(bigResponse, true)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "hospitalisations-region", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
