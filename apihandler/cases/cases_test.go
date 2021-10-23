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

var (
	testResponse = []*apiclient.APICasesResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     100,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 21, 0, 0, 0, 0, time.UTC)},
			Region:    "Brussels",
			Province:  "Brussels",
			Cases:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     120,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 22, 0, 0, 0, 0, time.UTC)},
			Region:    "",
			Province:  "",
			Cases:     5,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 10, 23, 0, 0, 0, 0, time.UTC)},
			Region:    "Flanders",
			Province:  "VlaamsBrabant",
			Cases:     100,
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
		"cases-province",
		"cases-region",
	}, targets)
}

func TestHandler_TableQuery_Cases(t *testing.T) {
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

	response, err := h.Endpoints().TableQuery(context.Background(), "cases", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 2)
	assert.Equal(t, 250.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 125.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	getter.AssertExpectations(t)
}

func TestHandler_TableQuery_CasesByProvince(t *testing.T) {
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

	response, err := h.Endpoints().TableQuery(context.Background(), "cases-province", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 4)
	require.Len(t, response.Columns[0].Data, 2)

	assert.Equal(t, "(unknown)", response.Columns[1].Text)
	require.Len(t, response.Columns[1].Data, 2)
	assert.Equal(t, 0.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 5.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	assert.Equal(t, "Brussels", response.Columns[2].Text)
	require.Len(t, response.Columns[2].Data, 2)
	assert.Equal(t, 150.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 0.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	assert.Equal(t, "VlaamsBrabant", response.Columns[3].Text)
	require.Len(t, response.Columns[3].Data, 2)
	assert.Equal(t, 100.0, response.Columns[3].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 120.0, response.Columns[3].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	getter.AssertExpectations(t)
}

func TestHandler_TableQuery_CasesByRegion(t *testing.T) {
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

	response, err := h.Endpoints().TableQuery(context.Background(), "cases-region", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 4)
	require.Len(t, response.Columns[0].Data, 2)

	assert.Equal(t, "(unknown)", response.Columns[1].Text)
	require.Len(t, response.Columns[1].Data, 2)
	assert.Equal(t, 0.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 5.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	assert.Equal(t, "Brussels", response.Columns[2].Text)
	require.Len(t, response.Columns[2].Data, 2)
	assert.Equal(t, 150.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 0.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	assert.Equal(t, "Flanders", response.Columns[3].Text)
	require.Len(t, response.Columns[3].Data, 2)
	assert.Equal(t, 100.0, response.Columns[3].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 120.0, response.Columns[3].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

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
