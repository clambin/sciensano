package sciensano_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const testResponse = `[ 
	{"DATE": "2021-03-09", "REGION": "Flanders", "TESTS_ALL": 10, "TESTS_ALL_POS": 5},
	{"DATE": "2021-03-10", "REGION": "Flanders", "TESTS_ALL": 11, "TESTS_ALL_POS": 5},
	{"DATE": "2021-03-11", "REGION": "Flanders", "TESTS_ALL": 15, "TESTS_ALL_POS": 10}
]`

const vaccResponse = `[
	{"DATE": "2021-03-09", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 50 },
	{"DATE": "2021-03-09", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 100 },
	{"DATE": "2021-03-10", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 100 },
	{"DATE": "2021-03-10", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 150 },
	{"DATE": "2021-03-11", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 150 },
	{"DATE": "2021-03-11", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 200 },
	{"DATE": "2021-03-11", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "B", "Count": 50 }
]`

func server(req *http.Request) (resp *http.Response) {
	switch req.URL.Path {
	case "/Data/COVID19BE_tests.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(testResponse)),
		}
	case "/Data/COVID19BE_VACC.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(vaccResponse)),
		}
	default:
		resp = &http.Response{
			Status:     req.URL.Path + " not found",
			StatusCode: http.StatusNotFound,
		}
	}
	return
}

var bigResponse string

func initBigResponse() string {
	timestamp := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

	entries := make([]string, 0)
	for timestamp.Before(time.Now()) {
		for _, region := range []string{"Flanders", "Wallonia", "Brussels", "Ostbelgien"} {
			for _, ageGroup := range []string{"0-17", "18-34", "35-44", "45-54", "55-64", "65-74", "75-84", "84+"} {
				for _, dose := range []string{"A", "B"} {
					entries = append(entries, fmt.Sprintf(`	{"DATE": "%s", "REGION": "%s", "AGEGROUP": "%s", "DOSE": "%s", "Count": 0 }`,
						timestamp.Format("2006-01-02"), region, ageGroup, dose))
				}
			}
		}
		timestamp = timestamp.Add(24 * time.Hour)
	}
	return "[" + strings.Join(entries, ",") + "]"
}

func bigServer(req *http.Request) (resp *http.Response) {
	if bigResponse == "" {
		bigResponse = initBigResponse()
	}

	switch req.URL.Path {
	case "/Data/COVID19BE_tests.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(testResponse)),
		}
	case "/Data/COVID19BE_VACC.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(bigResponse)),
		}
	default:
		resp = &http.Response{
			Status:     req.URL.Path + " not found",
			StatusCode: http.StatusNotFound,
		}
	}
	return
}
