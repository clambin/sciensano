package apihandler

import (
	"github.com/clambin/covid19/pkg/grafana/apiserver"
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"time"
)

// APIHandler implements a Grafana SimpleJson API that gets BE covid stats
type APIHandler struct {
	Cache *cache.Cache
}

func Create() (*APIHandler, error) {
	c := cache.New(15 * time.Minute)
	go c.Run()
	return &APIHandler{
		Cache: c,
	}, nil
}

// Search returns all supported targets
func (apiHandler *APIHandler) Search() []string {
	return allTargets()
}

// Query the DB and return the requested targets
func (apiHandler *APIHandler) Query(request *apiserver.APIQueryRequest) (response []apiserver.APIQueryResponse, err error) {
	for _, target := range request.Targets {
		var (
			group string
		)

		if group = findTargetGroup(target.Target); group == "" {
			log.WithField("target", target.Target).Warning("invalid target")
			continue
		}

		err = nil
		switch group {
		case "tests":
			responseChannel := make(chan []sciensano.Test)
			apiHandler.Cache.Tests <- cache.TestsRequest{
				EndTime:  request.Range.To,
				Response: responseChannel,
			}
			testStats := <-responseChannel
			response = append(response, buildTestPart(testStats, target.Target))
			close(responseChannel)

		case "vaccine":
			responseChannel := make(chan []sciensano.Vaccination)
			apiHandler.Cache.Vaccinations <- cache.VaccinationsRequest{
				EndTime:  request.Range.To,
				Response: responseChannel,
			}
			vaccineStats := <-responseChannel
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target.Target))
			close(responseChannel)

		case "vac-age":
			responseChannel := make(chan []sciensano.Vaccination)
			apiHandler.Cache.Vaccinations <- cache.VaccinationsRequest{
				EndTime:  request.Range.To,
				Filter:   "AgeGroup",
				Value:    getAgeGroupFromTarget(target.Target),
				Response: responseChannel,
			}
			vaccineStats := <-responseChannel
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target.Target))
			close(responseChannel)

		case "vac-reg":
			responseChannel := make(chan []sciensano.Vaccination)
			apiHandler.Cache.Vaccinations <- cache.VaccinationsRequest{
				EndTime:  request.Range.To,
				Filter:   "Region",
				Value:    getRegionFromTarget(target.Target),
				Response: responseChannel,
			}
			vaccineStats := <-responseChannel
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target.Target))
			close(responseChannel)
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
