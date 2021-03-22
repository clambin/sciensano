package vaccines_test

import (
	"bytes"
	"github.com/clambin/gotools/httpstub"
	"github.com/clambin/sciensano/internal/vaccines"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestVaccines(t *testing.T) {
	v := vaccines.Create()
	v.HTTPClient = httpstub.NewTestClient(server)

	response := make(chan []vaccines.Batch)
	v.Request <- response
	batches := <-response

	if assert.Len(t, batches, 3) {
		assert.Equal(t, "A", batches[0].Manufacturer)
		assert.Equal(t, "B", batches[1].Manufacturer)
		assert.Equal(t, "C", batches[2].Manufacturer)

		accu := vaccines.AccumulateBatches(batches)

		if assert.Len(t, accu, 3) {
			assert.Equal(t, int64(300), accu[0].Amount)
			assert.Equal(t, int64(500), accu[1].Amount)
			assert.Equal(t, int64(600), accu[2].Amount)
		}
	}

}

const vaccinesResponse = `{
	"result":{
		"delivered":[
			{"date":"2021-03-18","amount":100,"manufacturer":"C"},
			{"date":"2021-03-15","amount":200,"manufacturer":"B"},
			{"date":"2021-03-12","amount":300,"manufacturer":"A"}
		]
	}
}`

func server(req *http.Request) (resp *http.Response) {
	if req.URL.Path == "/api/v1/delivered.json" {
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(vaccinesResponse)),
		}
	} else {
		resp = &http.Response{
			Status:     req.URL.Path + " not found",
			StatusCode: http.StatusNotFound,
		}
	}
	return
}
