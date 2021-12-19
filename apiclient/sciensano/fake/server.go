package fake

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Handler implements a fake Sciensano server
type Handler struct {
	Fail      bool
	Slow      bool
	Responses map[string]string
	Count     int
	lock      sync.Mutex
}

// Handle processes an incoming HTTP request
func (handler *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	handler.Count++

	if handler.Slow && wait(req.Context(), 1*time.Second) == false {
		http.Error(w, "context exceeded", http.StatusRequestTimeout)
		return
	}

	if handler.Fail {
		http.Error(w, "server set to Fail", http.StatusInternalServerError)
		return
	}

	if handler.Responses == nil {
		handler.Responses = defaultResponses
	}

	response, ok := handler.Responses[req.URL.Path]

	if ok {
		_, _ = w.Write([]byte(response))
	} else {
		http.Error(w, "endpoint not implemented: "+html.EscapeString(req.URL.Path), http.StatusNotImplemented)
	}
}

// BigResponse creates a big response
func (handler *Handler) BigResponse() {
	handler.Responses = defaultResponses
	handler.Responses["/Data/COVID19BE_VACC.json"] = bigVaccinationResponse()
}

func wait(ctx context.Context, duration time.Duration) (passed bool) {
	timer := time.NewTimer(duration)
loop:
	for {
		select {
		case <-timer.C:
			break loop
		case <-ctx.Done():
			return false
		}
	}
	return true
}

var defaultResponses = map[string]string{
	"/Data/COVID19BE_tests.json": `[ 
		{"DATE": "2021-03-09", "REGION": "Flanders", "TESTS_ALL": 10, "TESTS_ALL_POS": 5},
		{"DATE": "2021-03-10", "REGION": "Flanders", "TESTS_ALL": 11, "TESTS_ALL_POS": 5},
		{"DATE": "2021-03-11", "REGION": "Flanders", "TESTS_ALL": 15, "TESTS_ALL_POS": 10}
]`,

	"/Data/COVID19BE_VACC.json": `[
		{"DATE": "2021-03-09", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "COUNT": 50 },
		{"DATE": "2021-03-09", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "COUNT": 100 },
		{"DATE": "2021-03-10", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "COUNT": 100 },
		{"DATE": "2021-03-10", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "COUNT": 150 },
		{"DATE": "2021-03-11", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "COUNT": 150 },
		{"DATE": "2021-03-11", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "COUNT": 200 },
		{"DATE": "2021-03-11", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "B", "COUNT": 50 }
]`,

	"/Data/COVID19BE_CASES_AGESEX.json": `[
		{"DATE":"2020-03-01","PROVINCE":"VlaamsBrabant","REGION":"Flanders","AGEGROUP":"40-49","SEX":"M","CASES":1},
		{"DATE":"2020-03-01","PROVINCE":"Brussels","REGION":"Brussels","AGEGROUP":"40-49","SEX":"M","CASES":2}
]`,

	"/Data/COVID19BE_MORT.json": `[
		{"DATE":"2020-03-10","REGION":"Brussels","AGEGROUP":"85+","SEX":"F","DEATHS":1},
		{"DATE":"2020-03-10","REGION":"Brussels","AGEGROUP":"85+","FEMALE":"F","DEATHS":2}
]`,

	"/Data/COVID19BE_HOSP.json": `[
		{"DATE":"2020-03-15","PROVINCE":"Brussels","REGION":"Brussels","NR_REPORTING":15,"TOTAL_IN":58,"TOTAL_IN_ICU":11,"TOTAL_IN_RESP":8,"TOTAL_IN_ECMO":0,"NEW_IN":7,"NEW_OUT":2},
		{"DATE":"2020-03-15","PROVINCE":"VlaamsBrabant","REGION":"Flanders","NR_REPORTING":6,"TOTAL_IN":13,"TOTAL_IN_ICU":2,"TOTAL_IN_RESP":0,"TOTAL_IN_ECMO":1,"NEW_IN":2,"NEW_OUT":0}
]`,
}

func bigVaccinationResponse() string {
	timestamp := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

	entries := make([]string, 0)
	count := 0
	for timestamp.Before(time.Now()) {
		count++
		for _, region := range []string{"Flanders", "Wallonia", "Brussels", "Ostbelgien"} {
			for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "84+"} {
				for _, dose := range []string{"A", "B"} {
					entries = append(entries, fmt.Sprintf(`	{"DATE": "%s", "REGION": "%s", "AGEGROUP": "%s", "DOSE": "%s", "Count": %d }`,
						timestamp.Format("2006-01-02"), region, ageGroup, dose, count))
				}
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return "[" + strings.Join(entries, ",") + "]"
}
