package apihandler_test

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/apihandler/mockapi"
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
	apiHandler.Client = &mockapi.API{Tests: mockapi.DefaultTests, Vaccinations: mockapi.DefaultVaccinations}

	request := &apiserver.APIQueryRequest{
		Range: apiserver.APIQueryRequestRange{},
		Targets: []apiserver.APIQueryRequestTarget{
			{Target: "tests-total"},
			{Target: "tests-positive"},
			{Target: "tests-rate"},
			{Target: "vaccinations-first"},
			{Target: "vaccinations-second"},
			{Target: "vaccinations-45-54-first"},
			{Target: "vaccinations-Flanders-first"},
			{Target: "invalid"},
		}}

	endDate := time.Date(2021, 01, 06, 0, 0, 0, 0, time.UTC)
	request.Range.To = endDate

	var response []apiserver.APIQueryResponse
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, len(request.Targets)-1)

		for _, entry := range response {
			if assert.NotZero(t, len(entry.DataPoints), entry.Target) {

				lastIndex := len(entry.DataPoints) - 1

				// Last entry should be endDate
				// Grafana time is in msec. Convert to sec + set TZ to UTC in order to compare
				ts := time.Unix(entry.DataPoints[lastIndex][1]/1000, 0).In(time.UTC)
				assert.Equal(t, endDate, ts)

				switch entry.Target {
				case "tests-total":
					assert.Equal(t, int64(10), entry.DataPoints[lastIndex][0])
				case "tests-positive":
					assert.Equal(t, int64(5), entry.DataPoints[lastIndex][0])
				case "tests-rate":
					assert.Equal(t, int64(100*5/10), entry.DataPoints[lastIndex][0])
				case "vaccinations-first":
					assert.Equal(t, int64(15), entry.DataPoints[lastIndex][0])
				case "vaccinations-second":
					assert.Equal(t, int64(10), entry.DataPoints[lastIndex][0])
				case "vaccinations-45-54-first":
					assert.Equal(t, int64(15), entry.DataPoints[lastIndex][0])
				case "vaccinations-Flanders-first":
					assert.Equal(t, int64(15), entry.DataPoints[lastIndex][0])
				}
			}
		}
	}
}
