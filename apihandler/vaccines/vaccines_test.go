package vaccines_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	sciensanoMock "github.com/clambin/sciensano/apiclient/sciensano/mocks"
	"github.com/clambin/sciensano/apiclient/vaccines"
	vaccinesMock "github.com/clambin/sciensano/apiclient/vaccines/mocks"
	vaccinesHandler "github.com/clambin/sciensano/apihandler/vaccines"
	"github.com/clambin/sciensano/measurement"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Search(t *testing.T) {
	client := &vaccinesMock.Getter{}
	r := reporter.NewCachedClient(time.Hour)
	r.Vaccines = client
	h := vaccinesHandler.New(r)

	targets := h.Search()
	assert.Equal(t, []string{"vaccines", "vaccines-stats", "vaccines-time"}, targets)
}

func TestHandler_TableQuery_Vaccines(t *testing.T) {
	client := &vaccinesMock.Getter{}
	r := reporter.NewCachedClient(time.Hour)
	r.Vaccines = client
	h := vaccinesHandler.New(r)

	timestamp := time.Now()
	client.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Time{Time: timestamp.Add(-24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.Batch{
				Date:   vaccines.Time{Time: timestamp},
				Amount: 200,
			},
			&vaccines.Batch{
				Date:   vaccines.Time{Time: timestamp.Add(24 * time.Hour)},
				Amount: 200,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 2)
	assert.Equal(t, 100.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 300.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	client.AssertExpectations(t)
}

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	vaccineClient := &vaccinesMock.Getter{}
	sciensanoClient := &sciensanoMock.Getter{}
	r := reporter.NewCachedClient(time.Hour)
	r.Vaccines = vaccineClient
	r.Sciensano = sciensanoClient
	h := vaccinesHandler.New(r)

	timestamp := time.Now()

	vaccineClient.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Time{Time: timestamp.Add(-24 * time.Hour)},
				Amount: 100,
			},
		}, nil)

	sciensanoClient.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return([]measurement.Measurement{
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-24 * time.Hour)},
				Dose:      "A",
				Count:     10,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-24 * time.Hour)},
				Dose:      "B",
				Count:     10,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines-stats", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 3)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 20.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 80.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])

	mock.AssertExpectationsForObjects(t, sciensanoClient, vaccineClient)
}

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	vaccineClient := &vaccinesMock.Getter{}
	sciensanoClient := &sciensanoMock.Getter{}
	r := reporter.NewCachedClient(time.Hour)
	r.Vaccines = vaccineClient
	r.Sciensano = sciensanoClient
	h := vaccinesHandler.New(r)

	vaccineClient.
		On("GetBatches", mock.AnythingOfType("*context.emptyCtx")).
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Time{Time: time.Now().Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.Batch{
				Date:   vaccines.Time{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, nil)

	sciensanoClient.
		On("GetVaccinations", mock.AnythingOfType("*context.emptyCtx")).
		Return([]measurement.Measurement{
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-6 * 24 * time.Hour)},
				Dose:      "A",
				Count:     50,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-5 * 24 * time.Hour)},
				Dose:      "A",
				Count:     25,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-4 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-3 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Dose:      "A",
				Count:     40,
			},
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-1 * 24 * time.Hour)},
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

	mock.AssertExpectationsForObjects(t, sciensanoClient, vaccineClient)
}
