package apihandler

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/sciensano"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

// APIHandler implements a Grafana SimpleJson API that gets BE covid stats
type APIHandler struct {
	apiClient sciensano.Client
}

func Create() (*APIHandler, error) {
	client := sciensano.Client{
		VaccinationsCacheDuration: 5 * time.Second,
	}
	return &APIHandler{apiClient: client}, nil
}

var (
	// TODO: should be build dynamically based on content. e.g. to provide stats by age group, we'd need:
	// tests-age-18-35-positive
	// tests-ago-18-35-total
	// etc.
	targets = map[string][]string{
		"tests":   sciensano.TestTargets,
		"vaccine": sciensano.GetVaccinationsTargets(),
		"vac-age": sciensano.GetVaccinationsByAgeTargets(),
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
		vaccineStats sciensano.Vaccinations
		group        string
	)

	for _, target := range request.Targets {
		if group = findTargetGroup(target.Target); group == "" {
			log.WithField("target", target.Target).Warning("invalid target")
			continue
		}
		if group == "tests" {
			if testStats, err = apiHandler.apiClient.GetTests(request.Range.To); err != nil {
				log.WithField("err", err).Warning("unable to get test statistics")
				continue
			}
			response = append(response, buildTestPart(testStats, target.Target))
		} else if group == "vaccine" {
			if vaccineStats, err = apiHandler.apiClient.GetVaccinations(request.Range.To); err != nil {
				log.WithField("err", err).Warning("unable to get vaccine statistics")
				continue
			}
			response = append(response, buildVaccinePart(vaccineStats, target.Target))
		} else if group == "vac-age" {
			if vaccineStats, err = apiHandler.apiClient.GetVaccinationsByAge(request.Range.To, sciensano.GetAgeGroupFromTarget(target.Target)); err != nil {
				log.WithField("err", err).Warning("unable to get vaccine statistics by age")
				continue
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

func buildVaccinePart(entries sciensano.Vaccinations, target string) (response apiserver.APIQueryResponse) {
	var timestamp, value int64

	response.Target = target
	response.DataPoints = make([][2]int64, len(entries))

	mode := sciensano.GetModeFromTarget(target)

	for index, entry := range entries {
		timestamp = entry.Timestamp.UnixNano() / 1000000
		value = 0
		if mode == "A" {
			value = int64(entry.FirstDose)
		} else if mode == "B" {
			value = int64(entry.SecondDose)
		}
		response.DataPoints[index] = [2]int64{value, timestamp}
	}
	return
}
