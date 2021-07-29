package apihandler_test

import (
	"github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/mockapi"
	mockVaccines "github.com/clambin/sciensano/vaccines/mock"
	"github.com/stretchr/testify/assert"
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
	"vaccines-stats":           false,
	"vaccines-time":            false,
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

func TestAPIHandler_Invalid(t *testing.T) {
	apiHandler, _ := apihandler.Create()

	apiHandler.Sciensano = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}
	apiHandler.Vaccines.HTTPClient = mockVaccines.GetServer()

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	var err error

	// Unknown target should return an error
	_, err = apiHandler.TableQuery("invalid", request)
	assert.NotNil(t, err)

}

func BenchmarkHandler_QueryTable(b *testing.B) {
	handler, err := apihandler.Create()

	if assert.Nil(b, err) {
		handler.Sciensano = &mockapi.API{Tests: buildTestTable(720), Vaccinations: buildVaccinationTable(720)}
		handler.Vaccines.HTTPClient = mockVaccines.GetServer()

		endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
		request := &grafana_json.TableQueryArgs{
			CommonQueryArgs: grafana_json.CommonQueryArgs{
				Range: grafana_json.QueryRequestRange{
					To: endDate,
				},
			},
		}
		b.ResetTimer()
		for target := range realTargets {
			for i := 0; i < 100; i++ {
				_, _ = handler.Endpoints().TableQuery(target, request)
			}
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

/*
func buildVaccineTable(size int) (table []vaccines.Batch) {
	testDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < size; i++ {
		table = append(table, vaccines.Batch{
			Date:  vaccines.Time(testDate),
			Amount: 200+i,
		})
		testDate = testDate.Add(24 * time.Hour)
	}
	return
}
*/
