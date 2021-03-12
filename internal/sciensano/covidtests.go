package sciensano

import (
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

var TestTargets = []string{"tests-positive", "tests-total", "tests-rate"}

type Test struct {
	Timestamp time.Time
	Total     int
	Positive  int
}

type Tests []Test

func (client *Client) GetTests(end time.Time) (results Tests, err error) {
	var apiResult []apiTestResponse

	if apiResult, err = client.getTests(); err == nil {
		results = groupTests(apiResult, end)
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
func (p Tests) Len() int {
	return len(p)
}

func (p Tests) Less(i, j int) bool {
	return p[i].Timestamp.Before(p[j].Timestamp)
}

func (p Tests) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
