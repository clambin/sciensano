package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
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

func TestHandler_TableQuery_VaccinesStats(t *testing.T) {
	vaccineClient := &mocks.Fetcher{}
	timestamp := time.Now().UTC()
	vaccineClient.
		On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines.TypeBatches).
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
			}, nil)
	vaccinationsClient := &mocks.Fetcher{}
	vaccinationsClient.
		On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).
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
		}, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccines.APIClient = vaccineClient
	r.Vaccinations.APIClient = vaccinationsClient
	h := vaccinesHandler.StatsHandler{Reporter: r}

	request := query.Request{Args: query.Args{Args: common.Args{
		Range: common.Range{From: timestamp.Add(-24 * time.Hour)},
	}}}

	response, err := h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "time", Data: query.TimeColumn{timestamp.Add(-24 * time.Hour), timestamp}},
		{Text: "vaccinations", Data: query.NumberColumn{20.0, 30.0}},
		{Text: "reserve", Data: query.NumberColumn{130.0, 320.0}},
	}}, response)

	request = query.Request{Args: query.Args{Args: common.Args{Range: common.Range{From: timestamp, To: timestamp}}}}

	response, err = h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "time", Data: query.TimeColumn{timestamp}},
		{Text: "vaccinations", Data: query.NumberColumn{30.0}},
		{Text: "reserve", Data: query.NumberColumn{320.0}},
	}}, response)

	mock.AssertExpectationsForObjects(t, vaccineClient, vaccinationsClient)
}
