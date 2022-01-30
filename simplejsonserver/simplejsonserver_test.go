package simplejsonserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestTargets(t *testing.T) {
	h := simplejsonserver.NewServer()

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
		"vacc-age-booster",
		"vacc-age-full",
		"vacc-age-partial",
		"vacc-age-rate-booster",
		"vacc-age-rate-full",
		"vacc-age-rate-partial",
		"vacc-region-booster",
		"vacc-region-full",
		"vacc-region-partial",
		"vacc-region-rate-booster",
		"vacc-region-rate-full",
		"vacc-region-rate-partial",
		"vaccination-lag",
		"vaccinations",
		"vaccines",
		"vaccines-manufacturer",
		"vaccines-stats",
		"vaccines-time",
	}, h.Targets())
}

func TestRun(t *testing.T) {
	h := simplejsonserver.NewServer()

	go func() {
		err := h.Run(8080)
		require.True(t, errors.Is(err, http.ErrServerClosed))
	}()

	require.Eventually(t, func() bool {
		response, err := http.Get("http://127.0.0.1:8080/health")
		return err == nil && response.StatusCode == http.StatusOK
	}, 30*time.Second, 10*time.Millisecond)

	ctx := context.Background()
	h.Reporter.APICache.Refresh(ctx)

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: time.Now()}}}}

	wg := sync.WaitGroup{}
	for target, handler := range h.Handlers {
		wg.Add(1)
		go func(handler simplejson.Handler, target string) {
			_, err := handler.Endpoints().Query(ctx, req)
			require.NoError(t, err, target)
			wg.Done()
		}(handler, target)
	}
	wg.Wait()

	response, err := http.Get("http://127.0.0.1:8080/health")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	_ = response.Body.Close()
	var result interface{}
	err = json.Unmarshal(body, &result)

	require.NoError(t, err)
	assert.Contains(t, result, "Handlers")
	assert.Contains(t, result, "APICache")
	assert.Contains(t, result, "ReporterCache")
	assert.Contains(t, result, "Demographics")
}

func BenchmarkHandlers_Run(b *testing.B) {
	h := simplejsonserver.NewServer()

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: time.Now()}}}}

	h.Reporter.APICache.Refresh(context.Background())
	_ = h.Demographics.GetRegionFigures()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, handler := range h.Handlers {
			_, err := handler.Endpoints().Query(ctx, req)
			assert.NoError(b, err)
		}
	}
}
