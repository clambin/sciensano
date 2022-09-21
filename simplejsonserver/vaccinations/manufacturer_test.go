package vaccinations_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManufacturerHandler(t *testing.T) {
	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f
	h := vaccinations.ManufacturerHandler{Reporter: r}

	response, err := h.Endpoints().Query(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{
			Text: "time",
			Data: query.TimeColumn{
				time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			Text: "(unknown)",
			Data: query.NumberColumn{7, 24},
		},
	}}, response)

	mock.AssertExpectationsForObjects(t, f)
}

func TestManufacturerHandler_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f
	h := vaccinations.ManufacturerHandler{
		Reporter: r,
	}

	_, err := h.Endpoints().Query(context.Background(), query.Request{})
	require.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}
