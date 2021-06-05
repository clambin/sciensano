package apihandler_test

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/vaccines/mock"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPIHandler_Tests(t *testing.T) {
	apiHandler, _ := apihandler.Create()

	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}
	apiHandler.Vaccines.HTTPClient = mock.GetServer()

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Tests
	if response, err = apiHandler.Endpoints().TableQuery("tests", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "total":
					assert.Equal(t, 10.0, data[len(data)-1])
				case "positive":
					assert.Equal(t, 5.0, data[len(data)-1])
				case "rate":
					assert.Equal(t, 0.5, data[len(data)-1])
				default:
					assert.Fail(t, "unexpected column", column.Text)
				}
			}
		}
	}
}
