package vaccines_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_TableQuery_Vaccines(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.APICache = cache
	h := vaccinesHandler.OverviewHandler{Reporter: r}

	timestamp := time.Now()
	cache.
		On("Get", "Vaccines").
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.Batch{
				Date:   vaccines.Timestamp{Time: timestamp},
				Amount: 200,
			},
			&vaccines.Batch{
				Date:   vaccines.Timestamp{Time: timestamp.Add(24 * time.Hour)},
				Amount: 200,
			},
		}, true)

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 2)
	assert.Equal(t, 100.0, response.Columns[1].Data.(simplejson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 300.0, response.Columns[1].Data.(simplejson.TableQueryResponseNumberColumn)[1])

	mock.AssertExpectationsForObjects(t, cache)
}

func TestHandler_TableQuery_VaccinesByManufacturer(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.APICache = cache
	h := vaccinesHandler.ManufacturerHandler{Reporter: r}

	timestamp := time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC)
	cache.
		On("Get", "Vaccines").
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:         vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
				Manufacturer: "A",
				Amount:       100,
			},
			&vaccines.Batch{
				Date:         vaccines.Timestamp{Time: timestamp},
				Manufacturer: "B",
				Amount:       200,
			},
			&vaccines.Batch{
				Date:         vaccines.Timestamp{Time: timestamp.Add(24 * time.Hour)},
				Manufacturer: "C",
				Amount:       200,
			},
		}, true)

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: timestamp}}}

	response, err := h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	assert.Equal(t, []simplejson.TableQueryResponseColumn{
		{
			Text: "timestamp",
			Data: simplejson.TableQueryResponseTimeColumn{
				time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{Text: "A", Data: simplejson.TableQueryResponseNumberColumn{100, 100}},
		{Text: "B", Data: simplejson.TableQueryResponseNumberColumn{0, 200}},
		{Text: "C", Data: simplejson.TableQueryResponseNumberColumn{0, 0}},
	}, response.Columns)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.APICache = cache
	h := vaccinesHandler.StatsHandler{Reporter: r}

	timestamp := time.Now()

	cache.
		On("Get", "Vaccines").
		Return(
			[]measurement.Measurement{
				&vaccines.Batch{
					Date:   vaccines.Timestamp{Time: timestamp.Add(-48 * time.Hour)},
					Amount: 50,
				},
				&vaccines.Batch{
					Date:   vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
					Amount: 100,
				},
				&vaccines.Batch{
					Date:   vaccines.Timestamp{Time: timestamp},
					Amount: 200,
				},
			}, true)

	cache.
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
			&sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Dose:      "C",
				Count:     10,
			},
		}, true)

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{
		Range: simplejson.Range{From: timestamp.Add(-24 * time.Hour)},
	}}

	response, err := h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 3)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{20.0, 30.0}, response.Columns[1].Data)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{130.0, 320.0}, response.Columns[2].Data)

	args = &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{From: timestamp, To: timestamp}}}

	response, err = h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 3)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{30.0}, response.Columns[1].Data)
	assert.Equal(t, simplejson.TableQueryResponseNumberColumn{320.0}, response.Columns[2].Data)

	mock.AssertExpectationsForObjects(t, cache, cache)
}

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.APICache = cache
	h := vaccinesHandler.DelayHandler{Reporter: r}

	cache.
		On("Get", "Vaccines").
		Return([]measurement.Measurement{
			&vaccines.Batch{
				Date:   vaccines.Timestamp{Time: time.Now().Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.Batch{
				Date:   vaccines.Timestamp{Time: time.Now().Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, true)

	cache.
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

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 2)
	require.Len(t, response.Columns[0].Data, 3)
	assert.Equal(t, 4, int(response.Columns[1].Data.(simplejson.TableQueryResponseNumberColumn)[0]))
	assert.Equal(t, 5, int(response.Columns[1].Data.(simplejson.TableQueryResponseNumberColumn)[1]))
	assert.Equal(t, 1, int(response.Columns[1].Data.(simplejson.TableQueryResponseNumberColumn)[2]))

	mock.AssertExpectationsForObjects(t, cache, cache)
}

func TestHandler_Failures(t *testing.T) {
	cache := &mockCache.Holder{}
	r := reporter.New(time.Hour)
	r.APICache = cache

	ctx := context.Background()
	args := simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{
		From: time.Time{},
		To:   time.Now(),
	}}}

	cache.On("Get", "Vaccinations").Return(nil, false)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	o := vaccinesHandler.OverviewHandler{Reporter: r}
	_, err := o.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	m := vaccinesHandler.ManufacturerHandler{Reporter: r}
	_, err = m.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	s := vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	d := vaccinesHandler.DelayHandler{Reporter: r}
	_, err = d.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]measurement.Measurement{}, true).Once()
	_, err = d.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]measurement.Measurement{}, true).Once()
	s = vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)
}
