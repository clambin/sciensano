package covidtests_test

import (
	"context"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/mocks"
	"github.com/clambin/sciensano/apihandler/covidtests"
	"github.com/clambin/sciensano/sciensano"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_Search(t *testing.T) {
	getter := &mocks.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := covidtests.New(client)

	targets := h.Search()
	assert.Equal(t, []string{"tests"}, targets)
}

func TestHandler_TableQuery(t *testing.T) {
	getter := &mocks.Getter{}
	client := &sciensano.Client{Getter: getter}
	h := covidtests.New(client)

	getter.
		On("GetTestResults", mock.AnythingOfType("*context.emptyCtx")).
		Return([]*apiclient.APITestResultsResponse{
			{
				TimeStamp: apiclient.TimeStamp{Time: time.Now().Add(-24 * time.Hour)},
				Positive:  10,
				Total:     20,
			},
		}, nil)

	args := &grafanajson.TableQueryArgs{CommonQueryArgs: grafanajson.CommonQueryArgs{Range: grafanajson.QueryRequestRange{
		From: time.Time{},
		To:   time.Now(),
	}}}

	response, err := h.Endpoints().TableQuery(context.Background(), "tests", args)
	require.NoError(t, err)
	require.Len(t, response.Columns, 4)
	require.Len(t, response.Columns[0].Data, 1)
	assert.Equal(t, 20.0, response.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 10.0, response.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[0])
	assert.Equal(t, 0.5, response.Columns[3].Data.(grafanajson.TableQueryResponseNumberColumn)[0])

	getter.AssertExpectations(t)
}
