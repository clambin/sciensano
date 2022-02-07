package cache

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	log "github.com/sirupsen/logrus"
	"time"
)

// Fetcher interface contains all functions an API interface needs to implement to be used by Cache
type Fetcher interface {
	Update(ctx context.Context, ch chan<- FetcherResponse)
}

// FetcherResponse contains one update from a Fetcher
type FetcherResponse struct {
	Name     string
	Response []apiclient.APIResponse
}

const maxRetries = 5
const retryDelay = 5 * time.Second

// Fetch is a convenience function. It calls an API and sends the data back to the cache.
// If the API fails, the call will be retried maxRetries number of times
func Fetch(ctx context.Context, ch chan<- FetcherResponse, name string, g func(context.Context) ([]apiclient.APIResponse, error)) {
	response := FetcherResponse{Name: name}
	var err error
	for i := 0; i < maxRetries; i++ {
		if response.Response, err = g(ctx); err == nil {
			ch <- response
			break
		}
		log.WithError(err).Warningf("API call for %s failed. retrying ...", name)
		time.Sleep(retryDelay)
	}
	if err != nil {
		log.WithError(err).Warningf("API call for %s failed after %d retries", name, maxRetries)
	}
}
