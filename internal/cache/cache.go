package cache

import (
	"github.com/clambin/sciensano/pkg/sciensano"
	log "github.com/sirupsen/logrus"
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
	EndTime  time.Time
	Filter   string
	Value    string
	Response chan []sciensano.Vaccination
}

func New(duration time.Duration) *Cache {
	return &Cache{
		Tests:        make(chan TestsRequest),
		Vaccinations: make(chan VaccinationsRequest),
		API: &sciensano.Client{
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
			msg.Response <- cache.getVaccinations(msg.EndTime, msg.Filter, msg.Value)
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

func (cache *Cache) getVaccinations(end time.Time, filter, value string) (vaccinations []sciensano.Vaccination) {
	var err error
	switch filter {
	case "":
		vaccinations, err = cache.GetVaccinations(end)
	case "AgeGroup":
		vaccinations, err = cache.GetVaccinationsByAge(end, value)
	case "Region":
		vaccinations, err = cache.GetVaccinationsByRegion(end, value)
	default:
		log.WithField("filter", filter).Warning("ignoring unsupported filter")
	}
	if err != nil {
		log.WithField("err", err).Warning("failed to get vaccination stats")
	}
	return
}
