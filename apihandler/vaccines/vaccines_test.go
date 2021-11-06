package vaccines_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	sciensanoMock "github.com/clambin/sciensano/apiclient/mocks"
	vaccinesHandler "github.com/clambin/sciensano/apihandler/vaccines"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	"github.com/clambin/sciensano/vaccines/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Search(t *testing.T) {
	client := &mocks.APIClient{}
	h := vaccinesHandler.New(nil, client)

	targets := h.Search()
	assert.Equal(t, []string{"vaccines", "vaccines-stats", "vaccines-time"}, targets)
}

func TestHandler_TableQuery_Vaccines(t *testing.T) {
	client := &mocks.APIClient{}
	h := vaccinesHandler.New(nil, client)

	client.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]*vaccines.Batch{
			{
				Date:   vaccines.Time{Time: time.Now().Add(-24 * time.Hour)},
				Amount: 100,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 100.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])

	client.AssertExpectations(t)
}

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	getter := &sciensanoMock.Getter{}
	sciensanoClient := sciensano.NewCachedClient(time.Hour)
	sciensanoClient.Getter = getter
	client := &mocks.APIClient{}
	h := vaccinesHandler.New(sciensanoClient, client)

	client.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]*vaccines.Batch{
			{
				Date:   vaccines.Time{Time: time.Now().Add(-24 * time.Hour)},
				Amount: 100,
			},
		}, nil)

	getter.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return([]apiclient.Measurement{
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-24 * time.Hour)},
				Dose:      "A",
				Count:     20,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines-stats", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 3)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 20.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 80.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])

	mock.AssertExpectationsForObjects(t, client, getter)
}

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	getter := &sciensanoMock.Getter{}
	sciensanoClient := sciensano.NewCachedClient(time.Hour)
	sciensanoClient.Getter = getter
	client := &mocks.APIClient{}
	h := vaccinesHandler.New(sciensanoClient, client)

	client.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]*vaccines.Batch{
			{
				Date:   vaccines.Time{Time: time.Now().Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			{
				Date:   vaccines.Time{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, nil)

	getter.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return([]apiclient.Measurement{
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-6 * 24 * time.Hour)},
				Dose:      "A",
				Count:     50,
			},
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-5 * 24 * time.Hour)},
				Dose:      "A",
				Count:     25,
			},
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-4 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-3 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Dose:      "A",
				Count:     40,
			},
			&apiclient.APIVaccinationsResponseEntry{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-1 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines-time", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 3)
	assert.Equal(t, 4, int(response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0]))
	assert.Equal(t, 5, int(response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1]))
	assert.Equal(t, 1, int(response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[2]))

	mock.AssertExpectationsForObjects(t, client, getter)
}
