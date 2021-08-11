package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/demographics"
	"github.com/clambin/sciensano/demographics/mock"
	"github.com/clambin/sciensano/sciensano/mockapi"
	mockVaccines "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHandler_Vaccinations(t *testing.T) {
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

	// Vaccinations
	if response, err = apiHandler.Endpoints().TableQuery("vaccinations", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "partial":
					assert.Equal(t, 15.0, data[len(data)-1])
				case "full":
					assert.Equal(t, 10.0, data[len(data)-1])
				default:
					assert.Fail(t, "unexpected column", column.Text)
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationsByAge(t *testing.T) {
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

	// Vaccinations grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery("vacc-age-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "45-54":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 10.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByAge_Rate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockVaccines.Handler))
	defer server.Close()

	demoServer := mock.New("")
	defer demoServer.Close()
	demo := demographics.New()
	demo.URL = demoServer.URL()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := demo.Run(ctx, 100*time.Millisecond)
		assert.NoError(t, err)
	}()

	assert.Eventually(t, demo.AvailableData, 2000*time.Millisecond, 10*time.Millisecond)

	apiHandler, _ := apihandler.Create(demo)

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

	// Vaccination rate grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery("vacc-age-rate-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "45-54":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 1960, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion(t *testing.T) {
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

	// Vaccinations grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery("vacc-region-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "Flanders":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 10.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion_Rate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(mockVaccines.Handler))
	defer server.Close()

	demoServer := mock.New("")
	defer demoServer.Close()
	demo := demographics.New()
	demo.URL = demoServer.URL()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = demo.Run(ctx, 100*time.Millisecond)
	}()

	assert.Eventually(t, demo.AvailableData, 500*time.Millisecond, 10*time.Millisecond)

	apiHandler, _ := apihandler.Create(demo)

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

	// Vaccination rate grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery("vacc-region-rate-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "Flanders":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 1648, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_Vaccination_Lag(t *testing.T) {
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

	// Lag
	if response, err = apiHandler.Endpoints().TableQuery("vaccination-lag", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "lag":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 1.0, data[len(data)-1])
					}
				}
			}
		}
	}
}
