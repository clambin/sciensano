package apihandler_test

import (
	"context"
	"fmt"
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/mockapi"
	"github.com/clambin/sciensano/vaccines"
	"github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHandler_Vaccines(t *testing.T) {
	server := mock.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))
	defer apiServer.Close()

	apiHandler, _ := apihandler.Create(nil)

	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}
	apiHandler.Vaccines.URL = apiServer.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Vaccines
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vaccines", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					lastDate := data[len(data)-1]
					assert.Equal(t, 2021, lastDate.Year())
					assert.Equal(t, time.Month(1), lastDate.Month())
					assert.Equal(t, 3, lastDate.Day())
				}
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "vaccines":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 600.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_Vaccines_Stats(t *testing.T) {
	server := mock.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))
	defer apiServer.Close()

	apiHandler, _ := apihandler.Create(nil)

	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}
	apiHandler.Vaccines.URL = apiServer.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Reserve
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vaccines-stats", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					lastDate := data[len(data)-1]
					assert.Equal(t, 2021, lastDate.Year())
					assert.Equal(t, time.Month(1), lastDate.Month())
					assert.Equal(t, 6, lastDate.Day())
				}
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "vaccinations":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 25.0, data[len(data)-1])
					}
				case "reserve":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 575.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_Vaccines_Time(t *testing.T) {
	server := mock.Server{}
	apiServer := httptest.NewServer(http.HandlerFunc(server.Handler))
	defer apiServer.Close()

	apiHandler, _ := apihandler.Create(nil)

	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.AltVaccinations}
	apiHandler.Vaccines.URL = apiServer.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Reserve
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vaccines-time", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					lastDate := data[len(data)-1]
					assert.Equal(t, 2021, lastDate.Year())
					assert.Equal(t, time.Month(1), lastDate.Month())
					assert.Equal(t, 4, lastDate.Day())
				}
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "time":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 1.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestVaccineDelay(t *testing.T) {
	vaccinations := []sciensano.Vaccination{{
		Timestamp:  time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
		FirstDose:  10,
		SecondDose: 0,
	}, {
		Timestamp:  time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  15,
		SecondDose: 1,
	}, {
		Timestamp:  time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
		FirstDose:  15,
		SecondDose: 4,
	}, {
		Timestamp:  time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  25,
		SecondDose: 5,
	}, {
		Timestamp:  time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		FirstDose:  35,
		SecondDose: 10,
	}, {
		Timestamp:  time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		FirstDose:  35,
		SecondDose: 15,
	}}

	batches := []vaccines.Batch{{
		Date:   vaccines.Time(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 20,
	}, {
		Date:   vaccines.Time(time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 40,
	}, {
		Date:   vaccines.Time(time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 50,
	}}

	expected := []struct {
		Timestamp time.Time
		Value     float64
	}{{
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Value:     45,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Value:     28,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Value:     42,
	}}

	timestamps, delays := apihandler.CalculateVaccineDelay(vaccinations, batches)

	if assert.Equal(t, len(expected), len(timestamps)) && assert.Equal(t, len(expected), len(delays)) {
		for index, entry := range expected {
			assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
			assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
		}

	}
}
