package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"
)

func NewTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(sciensanoHandler))
}

func sciensanoHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		w.Header().Set("Last-Modified", time.Now().Format(time.RFC1123))
		return
	}
	var data any
	switch req.URL.Path {
	case "/Data/COVID19BE_CASES_AGESEX.json":
		data = Cases()
	case "/Data/COVID19BE_HOSP.json":
		data = Hospitalisations()
	case "/Data/COVID19BE_MORT.json":
		data = Mortalities()
	case "/Data/COVID19BE_tests.json":
		data = TestResults()
	case "/Data/COVID19BE_VACC.json":
		data = Vaccinations()
	default:
		panic(req.URL.Path)
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
