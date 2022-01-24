package testresults_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/testresults"
	"github.com/clambin/simplejson/v2/common"
	"github.com/clambin/simplejson/v2/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_TableQuery(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = getter
	h := testresults.Handler{Reporter: client}

	getter.
		On("Get", "TestResults").
		Return([]measurement.Measurement{
			&sciensano.APITestResultsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: time.Now().Add(-24 * time.Hour)},
				Positive:  10,
				Total:     20,
			},
		}, true)

	args := query.Args{Args: common.Args{Range: common.Range{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 4)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 20.0, response.Columns[1].Data.(query.NumberColumn)[0])
	assert.Equal(t, 10.0, response.Columns[2].Data.(query.NumberColumn)[0])
	assert.Equal(t, 0.5, response.Columns[3].Data.(query.NumberColumn)[0])

	getter.AssertExpectations(t)
}

func TestHandler_Failure(t *testing.T) {
	getter := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = getter

	args := query.Args{}

	getter.
		On("Get", "TestResults").
		Return(nil, false)

	h := testresults.Handler{
		Reporter: client,
	}
	_, err := h.Endpoints().TableQuery(context.Background(), args)
	assert.Error(t, err)

	getter.AssertExpectations(t)
}
