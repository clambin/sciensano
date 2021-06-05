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

func TestAPIHandler_Vaccinations(t *testing.T) {
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

	// Vaccinations
	if response, err = apiHandler.Endpoints().TableQuery("vaccinations", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafana_json.TableQueryResponseNumberColumn:
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

	// Vaccinations grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery("vacc-age-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.TableQueryResponseNumberColumn:
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

	// Vaccination rate grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery("vacc-age-rate-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "45-54":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 6, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion(t *testing.T) {
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

	// Vaccinations grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery("vacc-region-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.TableQueryResponseNumberColumn:
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

	// Vaccination rate grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery("vacc-region-rate-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.TableQueryResponseNumberColumn:
				switch column.Text {
				case "Flanders":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 10, int(data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_Vaccination_Lag(t *testing.T) {
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

	// Lag
	if response, err = apiHandler.Endpoints().TableQuery("vaccination-lag", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.TableQueryResponseNumberColumn:
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
