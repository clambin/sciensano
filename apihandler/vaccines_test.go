package apihandler_test

import (
	"fmt"
	grafana_json "github.com/clambin/grafana-json"
	apihandler2 "github.com/clambin/sciensano/apihandler"
	sciensano2 "github.com/clambin/sciensano/sciensano"
	mockapi2 "github.com/clambin/sciensano/sciensano/mockapi"
	vaccines2 "github.com/clambin/sciensano/vaccines"
	mock2 "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHandler_Vaccines(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mock2.Handler))
	defer server.Close()

	apiHandler, _ := apihandler2.Create()

	apiHandler.Sciensano = &mockapi2.API{Tests: mockapi2.DefaultTests, Vaccinations: mockapi2.DefaultVaccinations}
	apiHandler.Vaccines.URL = server.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Vaccines
	if response, err = apiHandler.Endpoints().TableQuery("vaccines", request); assert.Nil(t, err) {
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
	server := httptest.NewServer(http.HandlerFunc(mock2.Handler))
	defer server.Close()

	apiHandler, _ := apihandler2.Create()

	apiHandler.Sciensano = &mockapi2.API{Tests: mockapi2.DefaultTests, Vaccinations: mockapi2.DefaultVaccinations}
	apiHandler.Vaccines.URL = server.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Reserve
	if response, err = apiHandler.Endpoints().TableQuery("vaccines-stats", request); assert.Nil(t, err) {
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
	server := httptest.NewServer(http.HandlerFunc(mock2.Handler))
	defer server.Close()

	apiHandler, _ := apihandler2.Create()

	apiHandler.Sciensano = &mockapi2.API{Tests: mockapi2.DefaultTests, Vaccinations: mockapi2.AltVaccinations}
	apiHandler.Vaccines.URL = server.URL

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var response *grafana_json.TableQueryResponse
	var err error

	// Reserve
	if response, err = apiHandler.Endpoints().TableQuery("vaccines-time", request); assert.Nil(t, err) {
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
	vaccinations := []sciensano2.Vaccination{{
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

	batches := []vaccines2.Batch{{
		Date:   vaccines2.Time(time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 20,
	}, {
		Date:   vaccines2.Time(time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)),
		Amount: 40,
	}, {
		Date:   vaccines2.Time(time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)),
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

	timestamps, delays := apihandler2.CalculateVaccineDelay(vaccinations, batches)

	if assert.Equal(t, len(expected), len(timestamps)) && assert.Equal(t, len(expected), len(delays)) {
		for index, entry := range expected {
			assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
			assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
		}

	}
}
