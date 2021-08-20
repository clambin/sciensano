package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano/apiclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAPIHandler_Tests(t *testing.T) {
	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.sciensanoClient.
		On("GetTestResults", mock.Anything).
		Return([]*apiclient.APITestResultsResponse{
			{
				TimeStamp: apiclient.TimeStamp{Time: endDate.Add(-48 * time.Hour)},
				Total:     100,
				Positive:  10,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: endDate.Add(-24 * time.Hour)},
				Total:     100,
				Positive:  10,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: endDate},
				Total:     100,
				Positive:  10,
			},
			{
				TimeStamp: apiclient.TimeStamp{Time: endDate.Add(24 * time.Hour)},
				Total:     100,
				Positive:  10,
			},
		}, nil,
		)

	// Tests
	response, err := stack.apiHandler.Endpoints().TableQuery(ctx, "tests", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		require.Len(t, column.Data, 3)
		switch data := column.Data.(type) {
		case grafanaJson.TableQueryResponseTimeColumn:
			assert.Equal(t, "timestamp", column.Text)
			assert.Equal(t, endDate, data[len(data)-1])
		case grafanaJson.TableQueryResponseNumberColumn:
			switch column.Text {
			case "total":
				assert.Equal(t, 100.0, data[len(data)-1])
			case "positive":
				assert.Equal(t, 10.0, data[len(data)-1])
			case "rate":
				assert.Equal(t, 0.1, data[len(data)-1])
			default:
				assert.Fail(t, "unexpected column", column.Text)
			}
		}
	}

	stack.sciensanoClient.AssertExpectations(t)
}
