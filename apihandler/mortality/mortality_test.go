package mortality_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	mockAPI "github.com/clambin/sciensano/apiclient/sciensano/mocks"
	"github.com/clambin/sciensano/apihandler/mortality"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter"
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
	testResponse = []measurement.Measurement{
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "85+",
			Deaths:    100,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			AgeGroup:  "25-34",
			Deaths:    150,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "25-34",
			Deaths:    120,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			AgeGroup:  "",
			Deaths:    5,
		},
		&sciensano.APIMortalityResponseEntry{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			AgeGroup:  "65-74",
			Deaths:    100,
		},
	}

	testCases = []TestCase{
		{
			Target: "mortality",
			Response: &grafanajson.TableQueryResponse{
				Columns: []grafanajson.TableQueryResponseColumn{
					{Text: "timestamp", Data: grafanajson.TableQueryResponseTimeColumn{time.Date(2021, time.October, 21, 0, 0, 0, 0, time.UTC), time.Date(2021, time.October, 22, 0, 0, 0, 0, time.UTC)}},
					{Text: "total", Data: grafanajson.TableQueryResponseNumberColumn{250.0, 125.0}},
				},
			},
		},
		{
			Target: "mortality-region",
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
			Target: "mortality-age",
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
	getter := &mockAPI.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	h := mortality.New(client)

	targets := h.Search()
	assert.Equal(t, []string{
		"mortality",
		"mortality-age",
		"mortality-region",
	}, targets)
}

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockAPI.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	h := mortality.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(testResponse, nil)

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
			bigResponse = append(bigResponse, &sciensano.APIMortalityResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Region:    region,
				Deaths:    i,
			})
		}

		timestamp = timestamp.Add(24 * time.Hour)
	}

	getter := &mockAPI.Getter{}
	client := reporter.NewCachedClient(time.Hour)
	client.Sciensano = getter
	h := mortality.New(client)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		To: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC),
	}}}

	getter.
		On("GetMortality", mock.AnythingOfType("*context.emptyCtx")).
		Return(bigResponse, nil)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), "mortality-region", args)
		require.NoError(b, err)
	}

	getter.AssertExpectations(b)
}
