package apiclient

import (
	"context"
	"sync"
	"time"
)

type Cache struct {
	APIClient
	Retention         time.Duration
	lock              sync.Mutex
	testResults       []*APITestResultsResponse
	testExpiry        time.Time
	testOnce          *sync.Once
	vaccinations      []*APIVaccinationsResponse
	vaccinationExpiry time.Time
	vaccinationOnce   *sync.Once
}

func (cache *Cache) GetTestResults(ctx context.Context) (results []*APITestResultsResponse, err error) {
	cache.lock.Lock()
	if cache.testOnce == nil || time.Now().After(cache.testExpiry) {
		metricCacheMiss.WithLabelValues("tests").Add(1.0)
		cache.testOnce = &sync.Once{}
		cache.testExpiry = time.Now().Add(cache.Retention)
	} else {
		metricCacheHit.WithLabelValues("tests").Add(1.0)
	}
	cache.lock.Unlock()

	cache.testOnce.Do(func() {
		if results, err = cache.APIClient.GetTestResults(ctx); err == nil {
			cache.testResults = results
		}
	})

	return cache.testResults, err
}

func (cache *Cache) GetVaccinations(ctx context.Context) (results []*APIVaccinationsResponse, err error) {
	cache.lock.Lock()
	if cache.vaccinationOnce == nil || time.Now().After(cache.vaccinationExpiry) {
		metricCacheMiss.WithLabelValues("vaccinations").Add(1.0)
		cache.vaccinationOnce = &sync.Once{}
		cache.vaccinationExpiry = time.Now().Add(cache.Retention)
	} else {
		metricCacheHit.WithLabelValues("vaccinations").Add(1.0)
	}
	cache.lock.Unlock()

	cache.vaccinationOnce.Do(func() {
		if results, err = cache.APIClient.GetVaccinations(ctx); err == nil {
			cache.vaccinations = results
		}
	})

	return cache.vaccinations, err
}
