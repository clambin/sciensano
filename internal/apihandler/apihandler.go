package apihandler

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"time"
)

// APIHandler implements a Grafana SimpleJson API that gets BE covid stats
type APIHandler struct {
	Client sciensano.API
}

func Create() (*APIHandler, error) {
	client := sciensano.Client{
		CacheDuration: 15 * time.Minute,
	}
	return &APIHandler{Client: &client}, nil
}

// Search returns all supported targets
func (apiHandler *APIHandler) Search() []string {
	return allTargets()
}

// Query the DB and return the requested targets
func (apiHandler *APIHandler) Query(request *apiserver.APIQueryRequest) (response []apiserver.APIQueryResponse, err error) {
	for _, target := range request.Targets {
		var (
			testStats    sciensano.Tests
			vaccineStats []sciensano.Vaccination
			group        string
		)

		if group = findTargetGroup(target.Target); group == "" {
			log.WithField("target", target.Target).Warning("invalid target")
			continue
		}

		err = nil
		switch group {
		case "tests":
			if testStats, err = apiHandler.Client.GetTests(request.Range.To); err == nil {
				response = append(response, buildTestPart(testStats, target.Target))
			}
		case "vaccine":
			if vaccineStats, err = apiHandler.Client.GetVaccinations(request.Range.To); err == nil {
				vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
				response = append(response, buildVaccinePart(vaccineStats, target.Target))
			}
		case "vac-age":
			if vaccineStats, err = apiHandler.Client.GetVaccinationsByAge(request.Range.To, getAgeGroupFromTarget(target.Target)); err == nil {
				vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
				response = append(response, buildVaccinePart(vaccineStats, target.Target))
			}
		case "vac-reg":
			if vaccineStats, err = apiHandler.Client.GetVaccinationsByRegion(request.Range.To, getRegionFromTarget(target.Target)); err == nil {
				vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
				response = append(response, buildVaccinePart(vaccineStats, target.Target))
			}
		}

		if err != nil {
			log.WithField("err", err).Warning("unable to get statistics")
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

func buildVaccinePart(entries []sciensano.Vaccination, target string) (response apiserver.APIQueryResponse) {
	var timestamp, value int64

	response.Target = target
	response.DataPoints = make([][2]int64, len(entries))

	mode := getModeFromTarget(target)

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
