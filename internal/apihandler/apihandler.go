package apihandler

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/sciensano"
	log "github.com/sirupsen/logrus"
	"sort"
)

// APIHandler implements a Grafana SimpleJson API that gets BE covid stats
type APIHandler struct {
	apiClient sciensano.APIClient
}

func Create() (*APIHandler, error) {
	return &APIHandler{apiClient: sciensano.APIClient{}}, nil
}

var (
	// TODO: should be build dynamically based on content. e.g. to provide stats by age group, we'd need:
	// tests-age-18-35-positive
	// tests-ago-18-35-total
	// etc.
	targets = map[string][]string{
		"tests":   sciensano.TestTargets,
		"vaccine": sciensano.VaccineTargets,
	}
)

func findTargetGroup(target string) string {
	for group, entries := range targets {
		for _, entry := range entries {
			if target == entry {
				return group
			}
		}
	}
	return ""
}

func allTargets() (output []string) {
	for _, entries := range targets {
		output = append(output, entries...)
	}
	sort.Strings(output)
	return
}

// Search returns all supported targets
func (apiHandler *APIHandler) Search() []string {
	return allTargets()
}

// Query the DB and return the requested targets
func (apiHandler *APIHandler) Query(request *apiserver.APIQueryRequest) (response []apiserver.APIQueryResponse, err error) {
	var (
		testStats    sciensano.Tests
		vaccineStats sciensano.Vaccines
		group        string
	)

	for _, target := range request.Targets {
		if group = findTargetGroup(target.Target); group == "" {
			log.WithField("target", target.Target).Warning("invalid target")
			continue
		}
		if group == "tests" {
			if testStats == nil {
				if testStats, err = apiHandler.apiClient.GetTests(request.Range.To); err != nil {
					log.WithField("err", err).Warning("unable to get test statistics")
					continue
				}
			}
			response = append(response, buildTestPart(testStats, target.Target))
		} else if group == "vaccine" {
			if vaccineStats == nil {
				if vaccineStats, err = apiHandler.apiClient.GetVaccines(request.Range.To); err != nil {
					log.WithField("err", err).Warning("unable to get vaccine statistics")
					continue
				}
			}
			response = append(response, buildVaccinePart(vaccineStats, target.Target))
		}
	}
	return
}

func buildTestPart(entries sciensano.Tests, target string) (response apiserver.APIQueryResponse) {
	var timestamp, value int64

	response.Target = target
	response.DataPoints = make([][2]int64, len(entries))
	for index, entry := range entries {
		timestamp = entry.Timestamp.UnixNano() / 1000000
		value = 0
		switch target {
		case "tests-total":
			value = int64(entry.Total)
		case "tests-positive":
			value = int64(entry.Positive)
		case "tests-rate":
			if entry.Total > 0 {
				value = int64(100 * entry.Positive / entry.Total)
			}
		}
		response.DataPoints[index] = [2]int64{value, timestamp}
	}
	return
}

func buildVaccinePart(entries sciensano.Vaccines, target string) (response apiserver.APIQueryResponse) {
	var timestamp, value int64

	response.Target = target
	response.DataPoints = make([][2]int64, len(entries))
	for index, entry := range entries {
		timestamp = entry.Timestamp.UnixNano() / 1000000
		value = 0
		if target == "vaccine-first" {
			value = int64(entry.FirstDose)
		} else {
			value = int64(entry.SecondDose)
		}
		response.DataPoints[index] = [2]int64{value, timestamp}
	}
	return
}
