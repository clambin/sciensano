package mock

import (
	log "github.com/sirupsen/logrus"
	"html"
	"net/http"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	log.Debug("apiHandler: " + html.EscapeString(req.URL.Path))

	if req.URL.Path == "/api/v1/delivered.json" {
		_, _ = w.Write([]byte(vaccinesResponse))
	} else {
		http.Error(w, "endpoint not implemented: "+req.URL.Path, http.StatusForbidden)
	}
}

const vaccinesResponse = `{
	"result":{
		"delivered":[
			{"date":"2021-01-03","amount":100,"manufacturer":"C"},
			{"date":"2021-01-02","amount":200,"manufacturer":"B"},
			{"date":"2021-01-01","amount":300,"manufacturer":"A"}
		]
	}
}`
