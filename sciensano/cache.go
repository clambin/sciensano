package sciensano

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"
)

type VaccinationsCache struct {
	URL          string
	HTTPClient   *http.Client
	retention    time.Duration
	expiry       time.Time
	lock         sync.Mutex
	once         *sync.Once
	result       []*apiVaccinationsResponse
	onceByAge    *sync.Once
	byAge        map[string][]*apiVaccinationsResponse
	onceByRegion *sync.Once
	byRegion     map[string][]*apiVaccinationsResponse
}

// NewVaccinationsCache creates a new cache for vaccination results, with the provided retention time.
func NewVaccinationsCache(retention time.Duration) *VaccinationsCache {
	return &VaccinationsCache{
		URL:        baseURL,
		HTTPClient: &http.Client{},
		retention:  retention,
	}
}

func (cache *VaccinationsCache) GetVaccinations(ctx context.Context) (result []*apiVaccinationsResponse, err error) {
	cache.lock.Lock()
	if cache.once == nil || time.Now().After(cache.expiry) {
		cache.once = &sync.Once{}
		cache.onceByAge = nil
		cache.onceByRegion = nil
		cache.expiry = time.Now().Add(cache.retention)
	}
	cache.lock.Unlock()

	cache.once.Do(func() {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, cache.URL+"/Data/COVID19BE_VACC.json", nil)

		var resp *http.Response
		resp, err = cache.HTTPClient.Do(req)

		if err != nil {
			return
		}

		defer func(body io.ReadCloser) {
			_ = body.Close()
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			err = errors.New(resp.Status)
			return
		}

		var body []byte
		if body, err = io.ReadAll(resp.Body); err == nil {
			var stats []*apiVaccinationsResponse
			if err = json.Unmarshal(body, &stats); err == nil {
				cache.result = stats
			}
		}
	})

	return cache.result, err
}

func (cache *VaccinationsCache) GetVaccinationsByAge(ctx context.Context) (result map[string][]*apiVaccinationsResponse, err error) {
	var vaccinations []*apiVaccinationsResponse
	vaccinations, err = cache.GetVaccinations(ctx)

	if err != nil {
		return
	}

	cache.lock.Lock()
	if cache.onceByAge == nil {
		cache.onceByAge = &sync.Once{}
	}
	cache.lock.Unlock()

	cache.onceByAge.Do(func() {
		output := make(map[string][]*apiVaccinationsResponse)
		for _, entry := range vaccinations {
			output[entry.AgeGroup] = append(output[entry.AgeGroup], entry)
		}
		cache.byAge = output
	})

	return cache.byAge, err
}

func (cache *VaccinationsCache) GetVaccinationsByRegion(ctx context.Context) (result map[string][]*apiVaccinationsResponse, err error) {
	var vaccinations []*apiVaccinationsResponse
	vaccinations, err = cache.GetVaccinations(ctx)

	if err != nil {
		return
	}

	cache.lock.Lock()
	if cache.onceByRegion == nil {
		cache.onceByRegion = &sync.Once{}
	}
	cache.lock.Unlock()

	cache.onceByRegion.Do(func() {
		output := make(map[string][]*apiVaccinationsResponse)
		for _, entry := range vaccinations {
			output[entry.Region] = append(output[entry.Region], entry)
		}
		cache.byRegion = output
	})

	return cache.byRegion, err
}

// TestResultsCache retrieves test results and caches them for a configured amount of time
type TestResultsCache struct {
	URL        string
	HTTPClient *http.Client
	retention  time.Duration
	expiry     time.Time
	lock       sync.Mutex
	once       *sync.Once
	result     []apiTestResponse
}

// NewTestResultsCache creates a new cache for test results, with the provided retention time.
func NewTestResultsCache(retention time.Duration) *TestResultsCache {
	return &TestResultsCache{
		URL:        baseURL,
		HTTPClient: &http.Client{},
		retention:  retention,
	}
}

func (cache *TestResultsCache) GetTestResults(ctx context.Context) (result []apiTestResponse, err error) {
	cache.lock.Lock()
	if cache.once == nil || time.Now().After(cache.expiry) {
		cache.once = &sync.Once{}
		cache.expiry = time.Now().Add(cache.retention)
	}
	cache.lock.Unlock()

	cache.once.Do(func() {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, cache.URL+"/Data/COVID19BE_tests.json", nil)

		var resp *http.Response
		resp, err = cache.HTTPClient.Do(req)

		if err == nil {
			if resp.StatusCode == http.StatusOK {
				var body []byte
				body, err = io.ReadAll(resp.Body)

				if err == nil {
					var stats []apiTestResponse
					err = json.Unmarshal(body, &stats)

					if err == nil {
						cache.result = stats
					}
				}
			} else {
				err = errors.New(resp.Status)
			}
			_ = resp.Body.Close()
		}
	})

	return cache.result, err
}
