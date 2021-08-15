package sciensano_test

import (
	"context"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestVaccinationCache(t *testing.T) {
	testServer := mock.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewVaccinationsCache(time.Hour)
	cache.URL = apiServer.URL

	results, err := cache.GetVaccinations(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, results)

	results, err = cache.GetVaccinations(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, results)

	assert.Equal(t, 1, testServer.Count)
}

func TestVaccinationCache_Parallel(t *testing.T) {
	testServer := mock.Handler{Slow: true}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewVaccinationsCache(time.Hour)
	cache.URL = apiServer.URL

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		results, err := cache.GetVaccinations(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		wg.Done()
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			_, err := cache.GetVaccinations(ctx)
			assert.NoError(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, testServer.Count)
}

func TestVaccinationCacheByAge(t *testing.T) {
	testServer := mock.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewVaccinationsCache(time.Hour)
	cache.URL = apiServer.URL

	results, err := cache.GetVaccinationsByAge(context.Background())
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Contains(t, results, "35-44")
	assert.Contains(t, results, "45-54")
}

func TestVaccinationCacheByRegion(t *testing.T) {
	testServer := mock.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewVaccinationsCache(time.Hour)
	cache.URL = apiServer.URL

	results, err := cache.GetVaccinationsByRegion(context.Background())
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Contains(t, results, "Flanders")
	assert.Contains(t, results, "Brussels")
}

func TestTestResultsCache(t *testing.T) {
	testServer := mock.Handler{}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewTestResultsCache(time.Hour)
	cache.URL = apiServer.URL

	results, err := cache.GetTestResults(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, results)

	results, err = cache.GetTestResults(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, results)

	assert.Equal(t, 1, testServer.Count)
}

func TestTestResultCache_Parallel(t *testing.T) {
	testServer := mock.Handler{Slow: true}
	apiServer := httptest.NewServer(http.HandlerFunc(testServer.Handle))

	cache := sciensano.NewTestResultsCache(time.Hour)
	cache.URL = apiServer.URL

	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		results, err := cache.GetTestResults(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, results)
		wg.Done()
	}()

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			_, err := cache.GetTestResults(ctx)
			assert.NoError(t, err)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(t, 1, testServer.Count)
}
