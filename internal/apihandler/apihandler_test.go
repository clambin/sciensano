package apihandler_test

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

var realTargets = map[string]bool{
	"tests":               false,
	"vaccinations":        false,
	"vacc-age-partial":    false,
	"vacc-age-full":       false,
	"vacc-region-partial": false,
	"vacc-region-full":    false,
	"vaccination-lag":     false,
}

func TestAPIHandler_Search(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)

	targets := apiHandler.Search()

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

func TestAPIHandler_Query(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)
	apiHandler.Cache.API = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}

	request := &grafana_json.QueryRequest{
		Range: grafana_json.QueryRequestRange{
			To: time.Now(),
		},
	}

	_, err = apiHandler.Query("tests", request)
	assert.NotNil(t, err)
}

func TestAPIHandler_QueryTable(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)
	apiHandler.Cache.API = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.QueryRequest{
		Range: grafana_json.QueryRequestRange{
			To: endDate,
		},
	}

	var response *grafana_json.QueryTableResponse

	// Tests
	if response, err = apiHandler.QueryTable("tests", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.QueryTableResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafana_json.QueryTableResponseNumberColumn:
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
	if response, err = apiHandler.QueryTable("vaccinations", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.QueryTableResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				assert.Equal(t, endDate, data[len(data)-1])
			case grafana_json.QueryTableResponseNumberColumn:
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
	if response, err = apiHandler.QueryTable("vacc-age-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.QueryTableResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.QueryTableResponseNumberColumn:
				switch column.Text {
				case "45-54":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 10.0, data[len(data)-1])
					}
				}
			}
		}
	}

	// Vaccinations grouped by Region
	if response, err = apiHandler.QueryTable("vacc-region-full", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.QueryTableResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.QueryTableResponseNumberColumn:
				switch column.Text {
				case "Flanders":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 10.0, data[len(data)-1])
					}
				}
			}
		}
	}

	// Lag
	if response, err = apiHandler.QueryTable("vaccination-lag", request); assert.Nil(t, err) {
		for _, column := range response.Columns {
			switch data := column.Data.(type) {
			case grafana_json.QueryTableResponseTimeColumn:
				assert.Equal(t, "timestamp", column.Text)
				if assert.NotZero(t, len(data)) {
					assert.Equal(t, endDate, data[len(data)-1])
				}
			case grafana_json.QueryTableResponseNumberColumn:
				switch column.Text {
				case "lag":
					if assert.NotZero(t, len(data)) {
						assert.Equal(t, 1.0, data[len(data)-1])
					}
				}
			}
		}
	}

	// Unknown target should return an error
	_, err = apiHandler.QueryTable("invalid", request)
	assert.NotNil(t, err)

}

func BenchmarkHandler_QueryTable(b *testing.B) {
	handler := apihandler.Handler{Cache: cache.New(0 * time.Minute)}
	handler.Cache.API = &mockapi.API{Tests: buildTestTable(720), Vaccinations: buildVaccinationTable(720)}
	go handler.Cache.Run()

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.QueryRequest{
		Range: grafana_json.QueryRequestRange{
			To: endDate,
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for j := 0; j < 100; j++ {
			for target := range realTargets {
				_, _ = handler.QueryTable(target, request)
			}
		}
		wg.Done()
	}()
	wg.Wait()
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
