package apiserver_test

import (
	"github.com/clambin/sciensano/pkg/grafana/apiserver"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSaveSearchResponse(t *testing.T) {
	response := []string{"foo", "bar"}

	bytes, err := apiserver.SaveSearchResponse(response)
	assert.Nil(t, err)
	assert.Equal(t, []byte(`["foo","bar"]`), bytes)
}

func TestLoadQueryRequest(t *testing.T) {
	input := []byte(`{
			"range": { 
				"from": "2020-01-01T00:00:00.000Z", 
				"to": "2020-12-31T23:59:59.0Z"
			},
			"targets": [
				{ "target": "A" },
				{ "target": "B" },
				{ "target": "C" }
			]}`)

	request, err := apiserver.LoadQueryRequest(input)
	assert.Nil(t, err)
	assert.Len(t, request.Targets, 3)
	assert.Equal(t, "A", request.Targets[0])
	assert.Equal(t, "B", request.Targets[1])
	assert.Equal(t, "C", request.Targets[2])
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), request.From)
	assert.Equal(t, time.Date(2020, 12, 31, 23, 59, 59, 0, time.UTC), request.To)
}

func TestSaveQueryResponse(t *testing.T) {
	expected := []byte(
		`[{"target":"A","datapoints":[[100,1577836800000],[101,1577836860000],[103,1577836920000]]},` +
			`{"target":"B","datapoints":[[100,1577836800000],[99,1577836860000],[98,1577836920000]]}]`)
	responses := []apiserver.QueryResponse{{
		Target: "A",
		Data: []apiserver.QueryResponseData{
			{Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
			{Timestamp: time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), Value: 101},
			{Timestamp: time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), Value: 103},
		},
	}, {
		Target: "B",
		Data: []apiserver.QueryResponseData{
			{Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
			{Timestamp: time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), Value: 99},
			{Timestamp: time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), Value: 98},
		},
	}}

	bytes, err := apiserver.SaveQueryResponse(responses)
	assert.Nil(t, err)
	assert.Equal(t, expected, bytes)
}
