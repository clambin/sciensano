package sciensano

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

type Test struct {
	Timestamp time.Time
	Total     int
	Positive  int
}

func (client *Client) GetTests(end time.Time) (results []Test, err error) {
	var apiResult []apiTestResponse

	if apiResult, err = client.getTests(); err == nil {
		results = groupTests(apiResult, end)
	}

	return
}

type apiTestResponse struct {
	TimeStamp string `json:"DATE"`
	Province  string `json:"PROVINCE"`
	Region    string `json:"REGION"`
	Total     int    `json:"TESTS_ALL"`
	Positive  int    `json:"TESTS_ALL_POS"`
}

func (client *Client) getTests() (stats []apiTestResponse, err error) {
	req, _ := http.NewRequest("GET", baseURL+"COVID19BE_tests.json", nil)

	var resp *http.Response
	if resp, err = client.apiClient.Do(req); err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == 200 {
			var (
				body []byte
			)
			if body, err = ioutil.ReadAll(resp.Body); err == nil {
				err = json.Unmarshal(body, &stats)
			}
		} else {
			err = errors.New(resp.Status)
		}
	}
	return
}

func groupTests(apiResult []apiTestResponse, end time.Time) (results Tests) {
	// Store the totals in a map
	totals := make(map[time.Time]Test, 0)
	for _, entry := range apiResult {
		if ts, err2 := time.Parse("2006-01-02", entry.TimeStamp); err2 == nil {
			// Skip anything after the specified end date
			if ts.After(end) {
				continue
			}

			var current Test
			var ok bool
			if current, ok = totals[ts]; ok == false {
				current = Test{Timestamp: ts}
			}
			current.Total += entry.Total
			current.Positive += entry.Positive
			totals[ts] = current
		} else {
			log.WithFields(log.Fields{
				"err":       err2,
				"timestamp": entry.TimeStamp,
			}).Warning("could not parse timestamp from API. skipping entry")
		}
	}
	// For each entry in the map, create an entry in the results slice
	for _, entry := range totals {
		results = append(results, entry)
	}
	// Maps are iterated in random order. Sort the final slice
	sort.Sort(results)

	return
}

// helper functions for sort.Sort([]Test)
type Tests []Test

func (p Tests) Len() int {
	return len(p)
}

func (p Tests) Less(i, j int) bool {
	return p[i].Timestamp.Before(p[j].Timestamp)
}

func (p Tests) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
