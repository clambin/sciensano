package apihandler_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/simplejson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	h := apihandler.NewServer()

	assert.Len(t, h.GetHandlers(), 6)
}

func TestRun(t *testing.T) {
	h := apihandler.NewServer()

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

	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: time.Now()}}}

	wg := sync.WaitGroup{}
	for _, handler := range h.GetHandlers() {
		for _, target := range handler.Endpoints().Search() {
			wg.Add(1)
			go func(handler simplejson.Handler, target string) {
				_, err := handler.Endpoints().TableQuery(ctx, target, args)
				require.NoError(t, err, target)
				wg.Done()
			}(handler, target)
		}
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
	h := apihandler.NewServer()

	ctx := context.Background()
	args := &simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: time.Now()}}}

	h.Reporter.APICache.Refresh(context.Background())
	_ = h.Demographics.GetRegionFigures()

	b.ResetTimer()
	for i := 0; i < 1; i++ {
		for _, handler := range h.GetHandlers() {
			for _, target := range handler.Endpoints().Search() {
				_, err := handler.Endpoints().TableQuery(ctx, target, args)
				assert.NoError(b, err)
			}
		}
	}

}
