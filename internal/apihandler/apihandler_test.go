package apihandler_test

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/clambin/sciensano/internal/sciensano"
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
	realTargets = append(realTargets, sciensano.TestTargets...)
	realTargets = append(realTargets, sciensano.GetVaccinationsTargets()...)
	realTargets = append(realTargets, sciensano.GetVaccinationsByAgeTargets()...)
	realTargets = append(realTargets, sciensano.GetVaccinationsByRegionTargets()...)
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

	request := &apiserver.APIQueryRequest{
		Range: apiserver.APIQueryRequestRange{},
		Targets: []apiserver.APIQueryRequestTarget{
			{Target: "tests-total"},
			{Target: "tests-positive"},
			{Target: "tests-rate"},
			{Target: "vaccinations-first"},
			{Target: "vaccinations-second"},
			{Target: "vaccinations-45-54-first"},
			{Target: "vaccinations-Brussels-first"},
			{Target: "invalid"},
		}}

	request.Range.To = time.Date(2020, 03, 01, 0, 0, 0, 0, time.UTC)
	var response []apiserver.APIQueryResponse
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, 6)

		for _, entry := range response {
			switch entry.Target {
			case "tests-total":
				if assert.Len(t, entry.DataPoints, 1) {
					assert.Equal(t, int64(82), entry.DataPoints[0][0])
				}
			case "tests-positive":
				if assert.Len(t, entry.DataPoints, 1) {
					assert.Equal(t, int64(0), entry.DataPoints[0][0])
				}
			case "tests-rate":
				if assert.Len(t, entry.DataPoints, 1) {
					assert.Equal(t, int64(0), entry.DataPoints[0][0])
				}
			case "vaccinations-first":
				assert.Len(t, entry.DataPoints, 0)
			case "vaccinations-second":
				assert.Len(t, entry.DataPoints, 0)
			case "vaccinations-45-54-first":
				assert.Len(t, entry.DataPoints, 0)
			case "vaccinations-Flanders-first":
				assert.Len(t, entry.DataPoints, 0)
			}
		}
	}

	request.Range.To = time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC)
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, 6)

		for _, entry := range response {
			switch entry.Target {
			case "tests-total":
				assert.Equal(t, int64(29048), entry.DataPoints[len(entry.DataPoints)-1][0])
			case "tests-positive":
				assert.Equal(t, int64(1898), entry.DataPoints[len(entry.DataPoints)-1][0])
			case "tests-rate":
				assert.Equal(t, int64(6), entry.DataPoints[len(entry.DataPoints)-1][0])
			case "vaccinations-first":
				assert.Equal(t, int64(298), entry.DataPoints[len(entry.DataPoints)-1][0])
			case "vaccinations-second":
				assert.Equal(t, int64(0), entry.DataPoints[len(entry.DataPoints)-1][0])
			}
		}
	}
}
