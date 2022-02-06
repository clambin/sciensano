package vaccinations_test

import (
	"context"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestManufacturerHandler(t *testing.T) {
	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	h := vaccinations.ManufacturerHandler{Reporter: client}

	response, err := h.Endpoints().Query(ctx, req)
	require.NoError(t, err)

	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{
			Text: "timestamp",
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
}

func TestManufacturerHandler_Failure(t *testing.T) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Vaccinations").Return(nil, false)

	h := vaccinations.ManufacturerHandler{
		Reporter: client,
	}

	_, err := h.Endpoints().Query(context.Background(), query.Request{})
	require.Error(t, err)
}
