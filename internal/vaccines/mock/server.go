package mock

import (
	"bytes"
	"github.com/clambin/gotools/httpstub"
	"io"
	"net/http"
)

const vaccinesResponse = `{
	"result":{
		"delivered":[
			{"date":"2021-03-18","amount":100,"manufacturer":"C"},
			{"date":"2021-03-15","amount":200,"manufacturer":"B"},
			{"date":"2021-03-12","amount":300,"manufacturer":"A"}
		]
	}
}`

func GetServer() *http.Client {
	return httpstub.NewTestClient(server)
}

func server(req *http.Request) (resp *http.Response) {
	if req.URL.Path == "/api/v1/delivered.json" {
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(vaccinesResponse)),
		}
	} else {
		resp = &http.Response{
			Status:     req.URL.Path + " not found",
			StatusCode: http.StatusNotFound,
		}
	}
	return
}
