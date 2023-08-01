package server

import (
	"context"
	"encoding/json"
	grafanaJSONServer "github.com/clambin/grafana-json-server"
	"github.com/clambin/sciensano/cache/sciensano"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	demographicsClient := mockDemographics.NewFetcher(t)
	//demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)

	h := New(demographicsClient, "")
	h.apiCache.Cases.Fetcher = &fetcher[sciensano.Cases]{}
	h.apiCache.Hospitalisations.Fetcher = &fetcher[sciensano.Hospitalisations]{}
	h.apiCache.Mortalities.Fetcher = &fetcher[sciensano.Mortalities]{}
	h.apiCache.TestResults.Fetcher = &fetcher[sciensano.TestResults]{}
	h.apiCache.Vaccinations.Fetcher = &fetcher[sciensano.Vaccinations]{}

	ctx := context.Background()
	go h.apiCache.AutoRefresh(ctx, time.Second)

	assert.Eventually(t, func() bool {
		if resp := h.apiCache.Cases.Get(ctx); resp == nil {
			return false
		}
		if resp := h.apiCache.Hospitalisations.Get(ctx); resp == nil {
			return false
		}
		if resp := h.apiCache.Mortalities.Get(ctx); resp == nil {
			return false
		}
		if resp := h.apiCache.TestResults.Get(ctx); resp == nil {
			return false
		}
		if resp := h.apiCache.Vaccinations.Get(ctx); resp == nil {
			return false
		}
		return true
	}, time.Second, 100*time.Millisecond)

	for target, handler := range h.handlers {
		t.Run(target, func(t *testing.T) {
			payload := `{"summary":"Total"}`
			if strings.HasPrefix(target, "vaccinations-rate") {
				payload = `{"summary":"ByAgeGroup"}`
			}

			req := grafanaJSONServer.QueryRequest{Targets: []grafanaJSONServer.QueryRequestTarget{
				{Target: target, Payload: []byte(payload)},
			}, Range: grafanaJSONServer.Range{To: time.Now()}}

			resp, err := handler.Query(ctx, target, req)
			require.NoError(t, err)
			assert.NotZero(t, len(resp.(grafanaJSONServer.TableResponse).Columns[0].Data.(grafanaJSONServer.TimeColumn)))
		})
	}
	//wg.Wait()

	b := httptest.NewRecorder()
	h.Health(b, nil)
	assert.Equal(t, http.StatusOK, b.Code)
	assert.Contains(t, b.Body.String(), `"DataSources": `)
	assert.Contains(t, b.Body.String(), `"ReporterCache": `)
}

func BenchmarkVaccinations(b *testing.B) {
	demographicsClient := mockDemographics.Fetcher{}
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)

	h := New(&demographicsClient, "")
	h.apiCache.Cases.Fetcher = &fetcher[sciensano.Cases]{}
	h.apiCache.Hospitalisations.Fetcher = &fetcher[sciensano.Hospitalisations]{}
	h.apiCache.Mortalities.Fetcher = &fetcher[sciensano.Mortalities]{}
	h.apiCache.TestResults.Fetcher = &fetcher[sciensano.TestResults]{}
	h.apiCache.Vaccinations.Fetcher = &fetcher[sciensano.Vaccinations]{}

	req := grafanaJSONServer.QueryRequest{
		Targets: []grafanaJSONServer.QueryRequestTarget{
			{Target: "vaccinations-rate-full", Payload: []byte(`{"summary":"ByAgeGroup"}`)},
		},
		Range: grafanaJSONServer.Range{To: time.Now()},
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.handlers["vaccinations-rate-full"].Query(ctx, "vaccinations-rate-full", req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

type fetcher[T any] struct {
	cache  T
	loaded bool
	lock   sync.Mutex
}

func (f *fetcher[T]) GetLastModified(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}

func (f *fetcher[T]) Fetch(_ context.Context) (T, error) {
	f.lock.Lock()
	defer f.lock.Unlock()
	if f.loaded {
		return f.cache, nil
	}

	input, err := os.Open(filepath.Join("..", "cache", "sciensano", "input", f.GetTarget()+".json"))
	if err != nil {
		return f.cache, err
	}

	defer func() {
		_ = input.Close()
	}()

	if err = json.NewDecoder(input).Decode(&f.cache); err == nil {
		f.loaded = true
	}
	return f.cache, err
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
