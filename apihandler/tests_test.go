package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/stretchr/testify/assert"
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

	var response *grafanaJson.TableQueryResponse
	var err error

	// Tests
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "tests", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "total":
					assert.Equal(t, 11.0, data[len(data)-1])
				case "positive":
					assert.Equal(t, 5.0, data[len(data)-1])
				case "rate":
					assert.Equal(t, 0.45454545454545453, data[len(data)-1])
				default:
					assert.Fail(t, "unexpected column", column.Text)
				}
			}
		}
	}
}
