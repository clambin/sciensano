package apihandler

import (
	"github.com/clambin/sciensano/internal/cache"
	"github.com/clambin/sciensano/pkg/grafana/apiserver"
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
func (apiHandler *APIHandler) Query(request *apiserver.QueryRequest) (response []apiserver.QueryResponse, err error) {
	for _, target := range request.Targets {
		var (
			group string
		)

		if group = findTargetGroup(target); group == "" {
			log.WithField("target", target).Warning("invalid target")
			continue
		}

		err = nil
		switch group {
		case "tests":
			req := cache.TestsRequest{
				EndTime:  request.To,
				Response: make(chan []sciensano.Test),
			}
			apiHandler.Cache.Tests <- req
			testStats := <-req.Response
			response = append(response, buildTestPart(testStats, target))

		case "vaccine":
			req := cache.VaccinationsRequest{
				EndTime:  request.To,
				Response: make(chan []sciensano.Vaccination),
			}
			apiHandler.Cache.Vaccinations <- req
			vaccineStats := <-req.Response
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target))

		case "vac-age":
			req := cache.VaccinationsRequest{
				EndTime:  request.To,
				Filter:   "AgeGroup",
				Value:    getAgeGroupFromTarget(target),
				Response: make(chan []sciensano.Vaccination),
			}
			apiHandler.Cache.Vaccinations <- req
			vaccineStats := <-req.Response
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target))

		case "vac-reg":
			req := cache.VaccinationsRequest{
				EndTime:  request.To,
				Filter:   "Region",
				Value:    getRegionFromTarget(target),
				Response: make(chan []sciensano.Vaccination),
			}
			apiHandler.Cache.Vaccinations <- req
			vaccineStats := <-req.Response
			vaccineStats = sciensano.AccumulateVaccinations(vaccineStats)
			response = append(response, buildVaccinePart(vaccineStats, target))
		}
	}
	return
}

func buildTestPart(entries sciensano.Tests, target string) (response apiserver.QueryResponse) {
	response.Target = target
	response.Data = make([]apiserver.QueryResponseData, len(entries))
	for index, entry := range entries {
		var value int64
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
		response.Data[index] = apiserver.QueryResponseData{Timestamp: entry.Timestamp, Value: value}
	}
	return
}

func buildVaccinePart(entries []sciensano.Vaccination, target string) (response apiserver.QueryResponse) {
	response.Target = target
	response.Data = make([]apiserver.QueryResponseData, len(entries))

	mode := getModeFromTarget(target)

	for index, entry := range entries {
		var value int64
		if mode == "A" {
			value = int64(entry.FirstDose)
		} else if mode == "B" {
			value = int64(entry.SecondDose)
		}
		response.Data[index] = apiserver.QueryResponseData{Timestamp: entry.Timestamp, Value: value}
	}
	return
}
