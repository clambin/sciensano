package simplejsonserver

import (
	"context"
	"encoding/json"
	"github.com/clambin/sciensano/cache/sciensano"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

	h, err := New(0, demographicsClient)
	require.NoError(t, err)
	h.apiCache.Cases.Fetcher = &fetcher[sciensano.Cases]{}
	h.apiCache.Hospitalisations.Fetcher = &fetcher[sciensano.Hospitalisations]{}
	h.apiCache.Mortalities.Fetcher = &fetcher[sciensano.Mortalities]{}
	h.apiCache.TestResults.Fetcher = &fetcher[sciensano.TestResults]{}
	h.apiCache.Vaccinations.Fetcher = &fetcher[sciensano.Vaccinations]{}

	ctx := context.Background()

	assert.Equal(t, []string{
		"cases",
		"cases-age",
		"cases-province",
		"cases-region",
		"hospitalisations",
		"hospitalisations-province",
		"hospitalisations-region",
		"mortality",
		"mortality-age",
		"mortality-region",
		"tests",
		"vacc-age",
		"vacc-age-rate",
		"vacc-age-rate-full",
		"vacc-age-rate-partial",
		"vacc-manufacturer",
		"vacc-region",
		"vacc-region-rate",
		"vacc-region-rate-full",
		"vacc-region-rate-partial",
		"vaccinations",
	}, h.server.Targets())

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: time.Now()}}}}

	wg := sync.WaitGroup{}
	wg.Add(len(h.server.Handlers))
	var count int
	for target, handler := range h.server.Handlers {
		t.Run(target, func(t *testing.T) {
			count++
			go func(handler simplejson.Handler, target string, counter int) {
				//t.Logf("%2d: %s started", count, target)
				var resp query.Response
				resp, err = handler.Endpoints().Query(ctx, req)
				assert.NotZero(t, len(resp.(query.TableResponse).Columns[0].Data.(query.TimeColumn)))
				//t.Logf("%2d: %s done. err: %v", count, target, err)
				assert.NoError(t, err, target, target)
				wg.Done()
			}(handler, target, count)
		})
	}
	wg.Wait()

	b := httptest.NewRecorder()
	h.Health(b, nil)
	assert.Equal(t, http.StatusOK, b.Code)
	assert.Contains(t, b.Body.String(), `"Handlers": `)
	assert.Contains(t, b.Body.String(), `"ReporterCache": `)
}

type fetcher[T any] struct {
}

func (f fetcher[T]) GetLastModified(_ context.Context) (time.Time, error) {
	return time.Now(), nil
}

func (f fetcher[T]) Fetch(_ context.Context) (T, error) {
	var t T
	input, err := os.Open(filepath.Join("..", "cache", "sciensano", "input", f.GetTarget()+".json"))
	if err == nil {
		err = json.NewDecoder(input).Decode(&t)
		_ = input.Close()
	}
	return t, err
}

func (f fetcher[T]) GetTarget() (target string) {
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
