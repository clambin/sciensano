package cache

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Cache struct {
	sciensano.API

	Tests        chan TestsRequest
	Vaccinations chan VaccinationsRequest
}

type TestsRequest struct {
	EndTime  time.Time
	Response chan []sciensano.Test
}

type VaccinationsRequest struct {
	EndTime         time.Time
	Filter          string
	Response        chan []sciensano.Vaccination
	GroupedResponse chan map[string][]sciensano.Vaccination
}

func New(duration time.Duration) *Cache {
	return &Cache{
		Tests:        make(chan TestsRequest),
		Vaccinations: make(chan VaccinationsRequest),
		API: &sciensano.Client{
			HTTPClient:    http.Client{Timeout: 15 * time.Second},
			CacheDuration: duration,
		},
	}
}

func (cache *Cache) Run() {
	log.Info("cache starting up")
loop:
	for {
		select {
		case msg, ok := <-cache.Tests:
			if ok == false {
				break loop
			}
			msg.Response <- cache.getTests(msg.EndTime)
		case msg, ok := <-cache.Vaccinations:
			if ok == false {
				break loop
			}
			switch msg.Filter {
			case "":
				result, _ := cache.GetVaccinations(msg.EndTime)
				msg.Response <- result
			case "AgeGroup":
				result, _ := cache.GetVaccinationsByAge(msg.EndTime)
				msg.GroupedResponse <- result
			case "Region":
				result, _ := cache.GetVaccinationsByRegion(msg.EndTime)
				msg.GroupedResponse <- result
			default:
				log.WithField("filter", msg.Filter).Warning("ignoring unsupported filter")
				msg.GroupedResponse <- map[string][]sciensano.Vaccination{}
			}
		}
	}
	log.Info("cache shutting down")
}

func (cache *Cache) Stop() {
	close(cache.Tests)
	close(cache.Vaccinations)
}

func (cache *Cache) getTests(end time.Time) (tests []sciensano.Test) {
	var err error
	if tests, err = cache.GetTests(end); err != nil {
		log.WithField("err", err).Warning("failed to get test results")
	}
	return
}
