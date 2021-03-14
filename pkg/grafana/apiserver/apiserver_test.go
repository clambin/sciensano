package apiserver_test

import (
	"bytes"
	"errors"
	"github.com/clambin/sciensano/pkg/grafana/apiserver"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestAPIServer_Full(t *testing.T) {
	server := apiserver.Create(newAPIHandler(), 8080)

	go func() {
		err := server.Run()

		assert.Nil(t, err)
	}()

	time.Sleep(1 * time.Second)

	body, err := call("http://localhost:8080/", "GET", "")
	if assert.Nil(t, err) {
		assert.Equal(t, "Hello", body)
	}

	body, err = call("http://localhost:8080/metrics", "GET", "")
	if assert.Nil(t, err) {
		assert.Contains(t, body, "grafana_api_duration_seconds")
		assert.Contains(t, body, "grafana_api_duration_seconds_sum")
		assert.Contains(t, body, "grafana_api_duration_seconds_count")
	}

	body, err = call("http://localhost:8080/search", "POST", "")
	if assert.Nil(t, err) {
		assert.Equal(t, `["A","B","Crash"]`, body)
	}

	req := `{
	"maxDataPoints": 100,
	"interval": "1y",
	"range": {
		"from": "2020-01-01T00:00:00.000Z",
		"to": "2020-12-31T00:00:00.000Z"
	},
	"targets": [
		{ "target": "A", "type": "foo" },
		{ "target": "B", "type": "foo" }
	]
}`
	body, err = call("http://localhost:8080/query", "POST", req)

	if assert.Nil(t, err) {
		assert.Equal(t, `[{"target":"A","datapoints":[[100,1577836800000],[101,1577836860000],[103,1577836920000]]},{"target":"B","datapoints":[[100,1577836800000],[99,1577836860000],[98,1577836920000]]}]`, body)
	}
}

func BenchmarkAPIServer(b *testing.B) {
	server := apiserver.Create(newAPIHandler(), 8080)

	go func() {
		err := server.Run()

		assert.Nil(b, err)
	}()

	time.Sleep(1 * time.Second)

	req := `{
	"maxDataPoints": 100,
	"interval": "1y",
	"range": {
		"from": "2020-01-01T00:00:00.000Z",
		"to": "2020-12-31T00:00:00.000Z"
	},
	"targets": [
		{ "target": "A", "type": "foo" },
		{ "target": "B", "type": "foo" }
	]
}`
	var body string
	var err error

	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		body, err = call("http://localhost:8080/query", "POST", req)
	}

	if assert.Nil(b, err) {
		assert.Equal(b, `[{"target":"A","datapoints":[[100,1577836800000],[101,1577836860000],[103,1577836920000]]},{"target":"B","datapoints":[[100,1577836800000],[99,1577836860000],[98,1577836920000]]}]`, body)
	}

}

func call(url, method, body string) (string, error) {
	client := &http.Client{}
	reqBody := bytes.NewBuffer([]byte(body))
	req, _ := http.NewRequest(method, url, reqBody)
	resp, err := client.Do(req)

	if err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			return string(body), nil
		}
	}

	return "", err
}

//
//
// Test APIHandler
//

type testAPIHandler struct {
}

func newAPIHandler() *testAPIHandler {
	return &testAPIHandler{}
}

func (apiHandler *testAPIHandler) Search() []string {
	return []string{"A", "B", "Crash"}
}

func (apiHandler *testAPIHandler) Query(request *apiserver.QueryRequest) ([]apiserver.QueryResponse, error) {
	var response = make([]apiserver.QueryResponse, 0)

	for _, target := range request.Targets {
		switch target {
		case "A":
			response = append(response, apiserver.QueryResponse{
				Target: "A",
				Data: []apiserver.QueryResponseData{
					{Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
					{Timestamp: time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), Value: 101},
					{Timestamp: time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), Value: 103},
				},
			})
		case "B":
			response = append(response, apiserver.QueryResponse{
				Target: "B",
				Data: []apiserver.QueryResponseData{
					{Timestamp: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), Value: 100},
					{Timestamp: time.Date(2020, 1, 1, 0, 1, 0, 0, time.UTC), Value: 99},
					{Timestamp: time.Date(2020, 1, 1, 0, 2, 0, 0, time.UTC), Value: 98},
				},
			})
		case "Crash":
			return response, errors.New("server crash")
		}
	}

	return response, nil
}
