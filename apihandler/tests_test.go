package apihandler_test

import (
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/sciensano/mockapi"
	mockVaccines "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHandler_Tests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockVaccines.Handler))
	defer server.Close()

	apiHandler, _ := apihandler.Create(nil)
	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}
	apiHandler.Vaccines.URL = server.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Tests
	if response, err = apiHandler.Endpoints().TableQuery("tests", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafanaJson.TableQueryResponseNumberColumn:
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
