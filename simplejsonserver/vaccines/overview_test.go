package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
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
	r.Vaccines.APICache = cache
	h := vaccinesHandler.OverviewHandler{Reporter: r}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp}}}}

	response, err := h.Endpoints().Query(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "time", Data: query.TimeColumn{timestamp.Add(-24 * time.Hour), timestamp}},
		{Text: "total", Data: query.NumberColumn{100, 300}},
	}}, response)

	mock.AssertExpectationsForObjects(t, cache)
}
