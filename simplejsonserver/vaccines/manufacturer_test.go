package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
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

func TestHandler_TableQuery_VaccinesByManufacturer(t *testing.T) {
	timestamp := time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC)

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines.TypeBatches).Return([]apiclient.APIResponse{
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
	}, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccines.APIClient = f
	h := vaccinesHandler.ManufacturerHandler{Reporter: r}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp}}}}

	response, err := h.Endpoints().Query(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{
			Text: "time",
			Data: query.TimeColumn{
				time.Date(2021, time.September, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2021, time.September, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{Text: "A", Data: query.NumberColumn{100, 100}},
		{Text: "B", Data: query.NumberColumn{0, 200}},
		{Text: "C", Data: query.NumberColumn{0, 0}},
	}}, response)

	mock.AssertExpectationsForObjects(t, f)
}
