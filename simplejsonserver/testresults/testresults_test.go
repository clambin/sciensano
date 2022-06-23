package testresults_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/testresults"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_TableQuery(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeTestResults).Return([]apiclient.APIResponse{
		&sciensano.APITestResultsResponse{
			TimeStamp: sciensano.TimeStamp{Time: time.Date(2022, 1, 26, 0, 0, 0, 0, time.UTC)},
			Positive:  10,
			Total:     20,
		},
	}, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.TestResults.APIClient = f
	h := testresults.Handler{Reporter: r}

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		From: time.Time{},
		To:   time.Now(),
	}}}}

	response, err := h.Endpoints().Query(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "time", Data: query.TimeColumn{time.Date(2022, 1, 26, 0, 0, 0, 0, time.UTC)}},
			{Text: "total", Data: query.NumberColumn{20.0}},
			{Text: "positive", Data: query.NumberColumn{10.0}},
			{Text: "rate", Data: query.NumberColumn{0.5}},
		},
	}, response)

	mock.AssertExpectationsForObjects(t, f)
}

func TestHandler_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeTestResults).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.TestResults.APIClient = f
	h := testresults.Handler{
		Reporter: r,
	}

	req := query.Request{}

	_, err := h.Endpoints().Query(context.Background(), req)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}
