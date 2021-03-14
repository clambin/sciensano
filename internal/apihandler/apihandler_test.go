package apihandler_test

import (
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/pkg/sciensano/mockapi"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func TestAPIHandler_Search(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)

	targets := apiHandler.Search()

	realTargets := make([]string, 0)
	realTargets = append(realTargets, apihandler.GetTestTargets()...)
	realTargets = append(realTargets, apihandler.GetVaccinationsTargets()...)
	realTargets = append(realTargets, apihandler.GetVaccinationsByAgeTargets()...)
	realTargets = append(realTargets, apihandler.GetVaccinationsByRegionTargets()...)
	sort.Strings(realTargets)

	if assert.Len(t, targets, len(realTargets)) {
		for index, target := range targets {
			assert.Equal(t, realTargets[index], target)
		}
	}
}

// TODO: API should be stubbed so we're not dependent on external data for testing
func TestAPIHandler_Query(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)
	apiHandler.Cache.API = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}

	request := &apiserver.QueryRequest{
		Targets: []string{
			"tests-total",
			"tests-positive",
			"tests-rate",
			"vaccinations-first",
			"vaccinations-second",
			"vaccinations-45-54-first",
			"vaccinations-Flanders-first",
			"invalid",
		}}

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request.To = endDate

	var response []apiserver.QueryResponse
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, len(request.Targets)-1)

		for _, entry := range response {
			if assert.NotZero(t, len(entry.Data), entry.Target) {
				lastIndex := len(entry.Data) - 1

				// Last entry should be endDate
				assert.Equal(t, endDate, entry.Data[lastIndex].Timestamp)

				switch entry.Target {
				case "tests-total":
					assert.Equal(t, int64(10), entry.Data[lastIndex].Value)
				case "tests-positive":
					assert.Equal(t, int64(5), entry.Data[lastIndex].Value)
				case "tests-rate":
					assert.Equal(t, int64(100*5/10), entry.Data[lastIndex].Value)
				case "vaccinations-first":
					assert.Equal(t, int64(15), entry.Data[lastIndex].Value)
				case "vaccinations-second":
					assert.Equal(t, int64(10), entry.Data[lastIndex].Value)
				case "vaccinations-45-54-first":
					assert.Equal(t, int64(15), entry.Data[lastIndex].Value)
				case "vaccinations-Flanders-first":
					assert.Equal(t, int64(15), entry.Data[lastIndex].Value)
				}
			}
		}
	}
}
