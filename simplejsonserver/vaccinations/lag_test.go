package vaccinations_test

import (
	"context"
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

func TestLagHandler(t *testing.T) {
	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	client := reporter.New(time.Hour)
	client.Vaccinations.APIClient = f

	h := vaccinations.LagHandler{Reporter: client}

	response, err := h.Endpoints().Query(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "time", Data: query.TimeColumn{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)}},
			{Text: "lag", Data: query.NumberColumn{0, 0}},
		},
	}, response)

	mock.AssertExpectationsForObjects(t, f)
}
