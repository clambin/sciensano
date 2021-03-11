package apihandler_test

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/apihandler"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TODO: API should be stubbed so we're not dependent on external data for testing
func TestAPIHandler_Query(t *testing.T) {
	apiHandler, err := apihandler.Create()
	assert.Nil(t, err)

	request := &apiserver.APIQueryRequest{
		Range: apiserver.APIQueryRequestRange{},
		Targets: []apiserver.APIQueryRequestTarget{
			{Target: "tests-total"},
			{Target: "tests-positive"},
			{Target: "vaccine-first"},
			{Target: "vaccine-second"},
			{Target: "invalid"},
		}}

	request.Range.To = time.Date(2020, 03, 01, 0, 0, 0, 0, time.UTC)
	var response []apiserver.APIQueryResponse
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, 4)

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
			case "vaccine-first":
				assert.Len(t, entry.DataPoints, 0)
			case "vaccine-second":
				assert.Len(t, entry.DataPoints, 0)
			}
		}
	}

	request.Range.To = time.Date(2020, 12, 28, 0, 0, 0, 0, time.UTC)
	response, err = apiHandler.Query(request)

	if assert.Nil(t, err) {
		assert.Len(t, response, 4)

		for _, entry := range response {
			switch entry.Target {
			case "tests-total":
				if assert.Len(t, entry.DataPoints, 303) {
					assert.Equal(t, int64(82), entry.DataPoints[0][0])
				}
			case "tests-positive":
				if assert.Len(t, entry.DataPoints, 303) {
					assert.Equal(t, int64(0), entry.DataPoints[0][0])
				}
			case "vaccine-first":
				if assert.Len(t, entry.DataPoints, 1) {
					assert.Equal(t, int64(298), entry.DataPoints[0][0])
				}
			case "vaccine-second":
				if assert.Len(t, entry.DataPoints, 1) {
					assert.Equal(t, int64(0), entry.DataPoints[0][0])
				}
			}
		}
	}

}
