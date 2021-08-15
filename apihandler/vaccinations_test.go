package apihandler_test

import (
	"context"
	grafanaJson "github.com/clambin/grafana-json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPIHandler_Vaccinations(t *testing.T) {
	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Vaccinations
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vaccinations", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafanaJson.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafanaJson.TableQueryResponseNumberColumn:
				switch column.Text {
				case "partial":
					assert.Equal(t, 400.0, data[len(data)-1])
				case "full":
					assert.Equal(t, 0.0, data[len(data)-1])
				default:
					assert.Fail(t, "unexpected column", column.Text)
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationsByAge(t *testing.T) {
	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Vaccinations grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vacc-age-full", request); assert.Nil(t, err) {
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
						assert.Equal(t, 0.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByAge_Rate(t *testing.T) {
	assert.Eventually(t, demoAPIServer.AvailableData, 10*time.Second, 10*time.Millisecond)

	endDate := time.Date(2021, 03, 10, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Vaccination rate grouped by Age
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vacc-age-rate-partial", request); assert.Nil(t, err) {
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
						assert.Equal(t, 49019, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion(t *testing.T) {
	endDate := time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Vaccinations grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vacc-region-full", request); assert.Nil(t, err) {
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
						assert.Equal(t, 50.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func TestAPIHandler_VaccinationByRegion_Rate(t *testing.T) {
	assert.Eventually(t, demoAPIServer.AvailableData, 500*time.Millisecond, 10*time.Millisecond)

	endDate := time.Date(2021, 03, 11, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Vaccination rate grouped by Region
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vacc-region-rate-full", request); assert.Nil(t, err) {
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
						assert.Equal(t, 8242, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}
}

func TestAPIHandler_Vaccination_Lag(t *testing.T) {
	endDate := time.Date(2021, 03, 11, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	var response *grafanaJson.TableQueryResponse
	var err error

	// Lag
	if response, err = apiHandler.Endpoints().TableQuery(context.Background(), "vaccination-lag", request); assert.Nil(t, err) {
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
						assert.Equal(t, 2.0, data[len(data)-1])
					}
				}
			}
		}
	}
}

func BenchmarkHandler_TableQuery_VaccinationRate(b *testing.B) {
	sciensanoServer.BigResponse()

	assert.Eventually(b, demoAPIServer.AvailableData, 500*time.Millisecond, 10*time.Millisecond)

	endDate := time.Date(2022, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafanaJson.TableQueryArgs{
		CommonQueryArgs: grafanaJson.CommonQueryArgs{
			Range: grafanaJson.QueryRequestRange{To: endDate},
		},
	}

	for i := 0; i < 10; i++ {
		_, err := apiHandler.Endpoints().TableQuery(context.Background(), "vacc-region-rate-full", request)
		assert.NoError(b, err)
	}
}
