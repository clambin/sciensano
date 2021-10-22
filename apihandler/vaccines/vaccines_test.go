package vaccines_test

import (
	"context"
	"fmt"
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
	sciensanoClient := &sciensano.Client{Getter: getter}
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
		Return([]*apiclient.APIVaccinationsResponse{
			{
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
	sciensanoClient := &sciensano.Client{Getter: getter}
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
		Return([]*apiclient.APIVaccinationsResponse{
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-6 * 24 * time.Hour)},
				Dose:      "A",
				Count:     50,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-5 * 24 * time.Hour)},
				Dose:      "A",
				Count:     25,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-4 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-3 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Dose:      "A",
				Count:     40,
			},
			{
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

func TestVaccineDelay(t *testing.T) {
	vaccinations := []sciensano.Vaccination{{
		Timestamp: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
		Partial:   10,
		Full:      0,
	}, {
		Timestamp: time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
		Partial:   15,
		Full:      1,
	}, {
		Timestamp: time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
		Partial:   15,
		Full:      4,
	}, {
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Partial:   25,
		Full:      5,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Partial:   35,
		Full:      10,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Partial:   35,
		Full:      15,
	}}

	batches := []*vaccines.Batch{{
		Date:   vaccines.Time{Time: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 20,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 40,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 50,
	}}

	expected := []struct {
		Timestamp time.Time
		Value     float64
	}{{
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Value:     45,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Value:     28,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Value:     42,
	}}

	timestamps, delays := vaccinesHandler.CalculateVaccineDelay(vaccinations, batches)

	if assert.Equal(t, len(expected), len(timestamps)) && assert.Equal(t, len(expected), len(delays)) {
		for index, entry := range expected {
			assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
			assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
		}

	}
}
