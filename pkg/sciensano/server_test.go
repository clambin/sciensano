package sciensano_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

const testResponse = `[ 
	{"DATE": "2021-03-09", "REGION": "Flanders", "TESTS_ALL": 10, "TESTS_ALL_POS": 5},
	{"DATE": "2021-03-10", "REGION": "Flanders", "TESTS_ALL": 11, "TESTS_ALL_POS": 5},
	{"DATE": "2021-03-11", "REGION": "Flanders", "TESTS_ALL": 15, "TESTS_ALL_POS": 10}
]`

// 	TimeStamp string `json:"DATE"`
//	Region    string `json:"REGION"`
//	AgeGroup  string `json:"AGEGROUP"`
//	Gender    string `json:"SEX"`
//	Dose      string `json:"DOSE"`
//	Count     int    `json:"Count"`

const vaccResponse = `[
	{"DATE": "2021-03-09", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 50 },
	{"DATE": "2021-03-09", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 100 },
	{"DATE": "2021-03-10", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 100 },
	{"DATE": "2021-03-10", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 150 },
	{"DATE": "2021-03-11", "REGION": "Brussels", "AGEGROUP": "35-44", "DOSE": "A", "Count": 150 },
	{"DATE": "2021-03-11", "REGION": "Flanders", "AGEGROUP": "45-54", "DOSE": "A", "Count": 200 }
]`

func server(req *http.Request) (resp *http.Response) {
	switch req.URL.Path {
	case "/Data/COVID19BE_tests.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(testResponse)),
		}
	case "/Data/COVID19BE_VACC.json":
		resp = &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString(vaccResponse)),
		}
	default:
		resp = &http.Response{
			Status:     req.URL.Path + " not found",
			StatusCode: http.StatusNotFound,
		}
	}
	return
}
