package vaccinations_test

import (
	"context"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v2/common"
	"github.com/clambin/simplejson/v2/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLagHandler(t *testing.T) {
	ctx := context.Background()
	args := query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}

	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	h := vaccinations.LagHandler{Reporter: client}

	cache.On("Get", "Vaccinations").Return(nil, false).Once()
	_, err := h.Endpoints().TableQuery(ctx, args)
	assert.Error(t, err)

	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true).Once()

	response, err := h.Endpoints().TableQuery(ctx, args)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: query.TimeColumn{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)}},
			{Text: "lag", Data: query.NumberColumn{0, 0}},
		},
	}, response)

}
