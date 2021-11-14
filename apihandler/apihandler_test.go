package apihandler_test

import (
	"context"
	"errors"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
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

	assert.Eventually(t, func() bool {
		response, err := http.Post("http://localhost:8080/search", "", nil)
		return err == nil && response.StatusCode == http.StatusOK
	}, 30*time.Second, 10*time.Millisecond)

	ctx := context.Background()
	args := &grafanajson.TableQueryArgs{
		CommonQueryArgs: grafanajson.CommonQueryArgs{
			Range: grafanajson.QueryRequestRange{To: time.Now()},
		},
	}

	for _, handler := range h.GetHandlers() {
		for _, target := range handler.Endpoints().Search() {
			_, err := handler.Endpoints().TableQuery(ctx, target, args)
			require.NoError(t, err)
		}
	}
}

func BenchmarkHandlers_Run(b *testing.B) {
	h := apihandler.NewServer()

	ctx := context.Background()
	args := &grafanajson.TableQueryArgs{
		CommonQueryArgs: grafanajson.CommonQueryArgs{
			Range: grafanajson.QueryRequestRange{To: time.Now()},
		},
	}

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
