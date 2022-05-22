package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_TableQuery_Vaccines(t *testing.T) {
	cache := &mockCache.Holder{}
	timestamp := time.Now().UTC()
	cache.
		On("Get", "Vaccines").
		Return([]apiclient.APIResponse{
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp},
				Amount: 200,
			},
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(24 * time.Hour)},
				Amount: 200,
			},
		}, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache
	h := vaccinesHandler.OverviewHandler{Reporter: r}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp}}}}

	response, err := h.Endpoints().Query(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "timestamp", Data: query.TimeColumn{timestamp.Add(-24 * time.Hour), timestamp}},
		{Text: "total", Data: query.NumberColumn{100, 300}},
	}}, response)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestHandler_TableQuery_VaccinesByManufacturer(t *testing.T) {
	cache := &mockCache.Holder{}
	timestamp := time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC)
	cache.
		On("Get", "Vaccines").
		Return([]apiclient.APIResponse{
			&vaccines.APIBatchResponse{
				Date:         vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
				Manufacturer: "A",
				Amount:       100,
			},
			&vaccines.APIBatchResponse{
				Date:         vaccines.Timestamp{Time: timestamp},
				Manufacturer: "B",
				Amount:       200,
			},
			&vaccines.APIBatchResponse{
				Date:         vaccines.Timestamp{Time: timestamp.Add(24 * time.Hour)},
				Manufacturer: "C",
				Amount:       200,
			},
		}, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache
	h := vaccinesHandler.ManufacturerHandler{Reporter: r}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp}}}}

	response, err := h.Endpoints().Query(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{
			Text: "timestamp",
			Data: query.TimeColumn{
				time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{Text: "A", Data: query.NumberColumn{100, 100}},
		{Text: "B", Data: query.NumberColumn{0, 200}},
		{Text: "C", Data: query.NumberColumn{0, 0}},
	}}, response)

	mock.AssertExpectationsForObjects(t, cache)
}

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	c := &mockCache.Holder{}
	timestamp := time.Now().UTC()
	c.
		On("Get", "Vaccines").
		Return(
			[]apiclient.APIResponse{
				&vaccines.APIBatchResponse{
					Date:   vaccines.Timestamp{Time: timestamp.Add(-48 * time.Hour)},
					Amount: 50,
				},
				&vaccines.APIBatchResponse{
					Date:   vaccines.Timestamp{Time: timestamp.Add(-24 * time.Hour)},
					Amount: 100,
				},
				&vaccines.APIBatchResponse{
					Date:   vaccines.Timestamp{Time: timestamp},
					Amount: 200,
				},
			}, true)
	c.
		On("Get", "Vaccinations").
		Return([]apiclient.APIResponse{
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-72 * time.Hour)},
				Dose:      "A",
				Count:     0,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-48 * time.Hour)},
				Dose:      "A",
				Count:     0,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-24 * time.Hour)},
				Dose:      "A",
				Count:     10,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-24 * time.Hour)},
				Dose:      "B",
				Count:     10,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp},
				Dose:      "C",
				Count:     10,
			},
		}, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = c
	h := vaccinesHandler.StatsHandler{Reporter: r}

	request := query.Request{Args: query.Args{Args: common.Args{
		Range: common.Range{From: timestamp.Add(-24 * time.Hour)},
	}}}

	response, err := h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "timestamp", Data: query.TimeColumn{timestamp.Add(-24 * time.Hour), timestamp}},
		{Text: "vaccinations", Data: query.NumberColumn{20.0, 30.0}},
		{Text: "reserve", Data: query.NumberColumn{130.0, 320.0}},
	}}, response)

	request = query.Request{Args: query.Args{Args: common.Args{Range: common.Range{From: timestamp, To: timestamp}}}}

	response, err = h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "timestamp", Data: query.TimeColumn{timestamp}},
		{Text: "vaccinations", Data: query.NumberColumn{30.0}},
		{Text: "reserve", Data: query.NumberColumn{320.0}},
	}}, response)

	mock.AssertExpectationsForObjects(t, c, c)
}

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	cache := &mockCache.Holder{}
	timestamp := time.Date(2022, 1, 26, 0, 0, 0, 0, time.UTC)
	cache.
		On("Get", "Vaccines").
		Return([]apiclient.APIResponse{
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, true)

	cache.
		On("Get", "Vaccinations").
		Return([]apiclient.APIResponse{
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-6 * 24 * time.Hour)},
				Dose:      "A",
				Count:     50,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-5 * 24 * time.Hour)},
				Dose:      "A",
				Count:     25,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-4 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-3 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Dose:      "A",
				Count:     40,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-1 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
		}, true)

	r := reporter.New(time.Hour)
	r.APICache = cache
	h := vaccinesHandler.DelayHandler{Reporter: r}

	request := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: timestamp,
	}}}}

	response, err := h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "timestamp", Data: query.TimeColumn{timestamp.Add(-3 * 24 * time.Hour), timestamp.Add(-2 * 24 * time.Hour), timestamp.Add(-24 * time.Hour)}},
		{Text: "time", Data: query.NumberColumn{4, 5, 1}},
	}}, response)

	mock.AssertExpectationsForObjects(t, cache, cache)
}

func TestHandler_Failures(t *testing.T) {
	cache := &mockCache.Holder{}

	cache.On("Get", "Vaccinations").Return(nil, false)

	ctx := context.Background()
	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: time.Now(),
	}}}}

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	o := vaccinesHandler.OverviewHandler{Reporter: r}
	_, err := o.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	m := vaccinesHandler.ManufacturerHandler{Reporter: r}
	_, err = m.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	s := vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	d := vaccinesHandler.DelayHandler{Reporter: r}
	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]apiclient.APIResponse{}, true).Once()
	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]apiclient.APIResponse{}, true).Once()

	s = vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().Query(ctx, req)
	assert.Error(t, err)
}
