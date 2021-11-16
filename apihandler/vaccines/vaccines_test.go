package vaccines_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	vaccinesHandler "github.com/clambin/sciensano/apihandler/vaccines"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Search(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.Vaccines = cache
	h := vaccinesHandler.New(r)

	targets := h.Search()
	assert.Equal(t, []string{"vaccines", "vaccines-stats", "vaccines-time"}, targets)
}

func TestHandler_TableQuery_Vaccines(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.Vaccines = cache
	h := vaccinesHandler.New(r)

	timestamp := time.Now()
	cache.
		On("Get", "Vaccines").
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
		}, true)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 2)
	assert.Equal(t, 100.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 300.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[1])

	mock.AssertExpectationsForObjects(t, cache)
}

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	vaccineCache := &mockCache.Holder{}
	sciensanoCache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.Vaccines = vaccineCache
	r.Sciensano = sciensanoCache
	h := vaccinesHandler.New(r)

	timestamp := time.Now()

	vaccineCache.
		On("Get", "Vaccines").
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Time{Time: timestamp.Add(-24 * time.Hour)},
				Amount: 100,
			},
		}, true)

	sciensanoCache.
		On("Get", "Vaccinations").
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
		}, true)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "vaccines-stats", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 3)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 20.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 80.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])

	mock.AssertExpectationsForObjects(t, sciensanoCache, vaccineCache)
}

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	vaccineCache := &mockCache.Holder{}
	sciensanoCache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.Vaccines = vaccineCache
	r.Sciensano = sciensanoCache
	h := vaccinesHandler.New(r)

	vaccineCache.
		On("Get", "Vaccines").
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Time{Time: time.Now().Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.Batch{
				Date:   vaccines.Time{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, true)

	sciensanoCache.
		On("Get", "Vaccinations").
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
		}, true)

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

	mock.AssertExpectationsForObjects(t, sciensanoCache, vaccineCache)
}
