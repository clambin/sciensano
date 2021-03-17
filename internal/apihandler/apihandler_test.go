package apihandler_test

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/pkg/sciensano"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAPIHandler_Search(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)

	realTargets := map[string]bool{
		"tests":               false,
		"vaccinations":        false,
		"vacc-age-partial":    false,
		"vacc-age-full":       false,
		"vacc-region-partial": false,
		"vacc-region-full":    false,
	}

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
}

func BenchmarkHandler_QueryTable(b *testing.B) {
	apiHandler, err := apihandler.Create()
	if assert.Nil(b, err) {
		apiHandler.Cache.API = &mockapi.API{Tests: buildTestTable(365), Vaccinations: buildVaccinationTable(365)}

		endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
		request := &grafana_json.QueryRequest{
			Range: grafana_json.QueryRequestRange{
				To: endDate,
			},
		}

		for i := 0; i < 100; i++ {
			// Tests
			_, _ = apiHandler.QueryTable("tests", request)
			_, _ = apiHandler.QueryTable("vaccinations-full", request)
			_, _ = apiHandler.QueryTable("vacc-age-full", request)
			_, _ = apiHandler.QueryTable("vacc-reg-full", request)
		}
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
