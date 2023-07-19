package server

import (
	"context"
	"encoding/json"
	"github.com/clambin/go-common/cache"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/cache/sciensano"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	demographicsClient := mockDemographics.NewFetcher(t)
	demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)

	h := New(demographicsClient)
	h.apiCache.Cases.Fetcher = &fetcher[sciensano.Cases]{}
	h.apiCache.Hospitalisations.Fetcher = &fetcher[sciensano.Hospitalisations]{}
	h.apiCache.Mortalities.Fetcher = &fetcher[sciensano.Mortalities]{}
	h.apiCache.TestResults.Fetcher = &fetcher[sciensano.TestResults]{}
	h.apiCache.Vaccinations.Fetcher = &fetcher[sciensano.Vaccinations]{}

	ctx := context.Background()
	go h.apiCache.AutoRefresh(ctx, time.Second)
	// TODO: fix race condition
	time.Sleep(1 * time.Second)

	req := grafanaJSONServer.QueryRequest{Range: grafanaJSONServer.Range{To: time.Now()}}

	wg := sync.WaitGroup{}
	wg.Add(len(h.dataSources))
	var count int
	for target, handler := range h.dataSources {
		t.Run(target, func(t *testing.T) {
			count++
			//go func(handler simplejson.Handler, target string, counter int) {
			//t.Logf("%2d: %s started", count, target)
			resp, err := handler.Query(ctx, target, req)
			assert.NotZero(t, len(resp.(grafanaJSONServer.TableResponse).Columns[0].Data.(grafanaJSONServer.TimeColumn)))
			//t.Logf("%2d: %s done. err: %v", count, target, err)
			assert.NoError(t, err, target, target)
			wg.Done()
			//}(handler, target, count)
		})
	}
	wg.Wait()

	b := httptest.NewRecorder()
	h.Health(b, nil)
	assert.Equal(t, http.StatusOK, b.Code)
	assert.Contains(t, b.Body.String(), `"DataSources": `)
	assert.Contains(t, b.Body.String(), `"ReporterCache": `)
}

func BenchmarkVaccinations(b *testing.B) {
	demographicsClient := mockDemographics.Fetcher{}
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)

	h := New(&demographicsClient)
	h.apiCache.Cases.Fetcher = &fetcher[sciensano.Cases]{}
	h.apiCache.Hospitalisations.Fetcher = &fetcher[sciensano.Hospitalisations]{}
	h.apiCache.Mortalities.Fetcher = &fetcher[sciensano.Mortalities]{}
	h.apiCache.TestResults.Fetcher = &fetcher[sciensano.TestResults]{}
	h.apiCache.Vaccinations.Fetcher = &fetcher[sciensano.Vaccinations]{}

	req := grafanaJSONServer.QueryRequest{Range: grafanaJSONServer.Range{To: time.Now()}}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.dataSources["vacc-age-rate-full"].Query(ctx, "", req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type fetcher[T any] struct {
	cache *cache.Cache[string, T]
	lock  sync.RWMutex
}

func (f *fetcher[T]) GetLastModified(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}

func (f *fetcher[T]) Fetch(_ context.Context) (T, error) {
	f.lock.RLock()
	if f.cache == nil {
		f.cache = cache.New[string, T](time.Minute, 0)
	}
	e, ok := f.cache.Get(f.GetTarget())
	f.lock.RUnlock()
	if ok {
		return e, nil
	}
	var t T
	input, err := os.Open(filepath.Join("..", "cache", "sciensano", "input", f.GetTarget()+".json"))
	if err == nil {
		err = json.NewDecoder(input).Decode(&t)
		_ = input.Close()
	}
	if err == nil {
		f.lock.Lock()
		f.cache.Add(f.GetTarget(), t)
		f.lock.Unlock()
	}
	return t, err
}

func (f *fetcher[T]) GetTarget() (target string) {
	var t T
	switch interface{}(t).(type) {
	case sciensano.Cases:
		target = "cases"
	case sciensano.Hospitalisations:
		target = "hospitalisations"
	case sciensano.Mortalities:
		target = "mortalities"
	case sciensano.TestResults:
		target = "testResults"
	case sciensano.Vaccinations:
		target = "vaccinations"
	}
	return target
}

func TestFilterVaccinations(t *testing.T) {
	var f fetcher[sciensano.Vaccinations]
	vaccinations, _ := f.Fetch(context.Background())
	filtered := filterVaccinations(vaccinations, sciensano.Full)
	assert.Len(t, filtered, 720)
}

func BenchmarkFilterVaccinations(b *testing.B) {
	var f fetcher[sciensano.Vaccinations]
	vaccinations, _ := f.Fetch(context.Background())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filterVaccinations(vaccinations, sciensano.Full)
	}
}
