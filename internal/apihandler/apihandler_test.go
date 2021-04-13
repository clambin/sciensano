package apihandler_test

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/vaccines/mock"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var realTargets = map[string]bool{
	"tests":                    false,
	"vaccinations":             false,
	"vacc-age-partial":         false,
	"vacc-age-full":            false,
	"vacc-age-rate-partial":    false,
	"vacc-age-rate-full":       false,
	"vacc-region-partial":      false,
	"vacc-region-full":         false,
	"vacc-region-rate-partial": false,
	"vacc-region-rate-full":    false,
	"vaccination-lag":          false,
	"vaccines":                 false,
}

func TestAPIHandler_Search(t *testing.T) {
	apiHandler, _ := apihandler.Create()
	targets := apiHandler.Endpoints().Search()

	for _, target := range targets {
		_, ok := realTargets[target]
		if assert.True(t, ok, "unexpected target: "+target) {
			realTargets[target] = true
		}
	}

	for target, found := range realTargets {
		assert.True(t, found, "missing target:"+target)
	}

}

func TestAPIHandler_TableQuery(t *testing.T) {
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
						assert.Equal(t, 1, int(1000000*data[len(data)-1]))
					}
				}
			}
		}
	}

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

	// Vaccines
	if response, err = apiHandler.Endpoints().TableQuery("vaccines", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.TableQueryResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					lastDate := data[len(data)-1]
					assert.Equal(t, 2021, lastDate.Year())
					assert.Equal(t, time.Month(3), lastDate.Month())
					assert.Equal(t, 18, lastDate.Day())
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

	// Unknown target should return an error
	_, err = apiHandler.TableQuery("invalid", request)
	assert.NotNil(t, err)

}

func BenchmarkHandler_QueryTable(b *testing.B) {
	handler, err := apihandler.Create()

	if assert.Nil(b, err) {
		handler.Sciensano = &mockapi.API{Tests: buildTestTable(720), Vaccinations: buildVaccinationTable(720)}

		endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
		request := &grafana_json.TableQueryArgs{
			CommonQueryArgs: grafana_json.CommonQueryArgs{
				Range: grafana_json.QueryRequestRange{
					To: endDate,
				},
			},
		}

		b.ResetTimer()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for target := range realTargets {
				if target == "vaccines" {
					continue
				}
				_, _ = handler.Endpoints().TableQuery(target, request)
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

func buildTestTable(size int) (table []sciensano.Test) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < size; i++ {
		table = append(table, sciensano.Test{
			Timestamp: testDate,
			Total:     i + 1,
			Positive:  i,
		})
		testDate = testDate.Add(24 * time.Hour)
	}
	return
}
func buildVaccinationTable(size int) (table []sciensano.Vaccination) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < size; i++ {
		table = append(table, sciensano.Vaccination{
			Timestamp:  testDate,
			FirstDose:  100 + i,
			SecondDose: i,
		})
		testDate = testDate.Add(24 * time.Hour)
	}
	return
}

func TestHandler_Annotations(t *testing.T) {
	handler, _ := apihandler.Create()
	// TODO: stub the API

	args := grafana_json.AnnotationRequestArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{
				To: time.Now(),
			},
		},
	}

	annotations, err := handler.Endpoints().Annotations("foo", "bar", &args)
	assert.Nil(t, err)
	assert.Greater(t, len(annotations), 0)
}
