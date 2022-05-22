package simplejsonserver_test

import (
	"context"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/simplejsonserver"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	demographicsClient := &mockDemographics.Fetcher{}
	h := simplejsonserver.NewServerWithDemographicsClient(demographicsClient)

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
		"vacc-manufacturer",
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

	demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)
	demographicsClient.On("Run", mock.AnythingOfType("*context.emptyCtx")).Return()

	ctx := context.Background()
	h.RunBackgroundTasks(ctx)

	require.Eventually(t, func() bool {
		h := h.Reporter.APICache.Stats()
		return len(h) == 6
	}, time.Minute, time.Second)

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

	for name, count := range h.Reporter.APICache.Stats() {
		assert.NotZero(t, count, name)
	}
}

func BenchmarkHandlers_Run(b *testing.B) {
	demographicsClient := &mockDemographics.Fetcher{}
	h := simplejsonserver.NewServerWithDemographicsClient(demographicsClient)

	demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: time.Now()}}}}

	go h.Reporter.APICache.Run(context.Background(), time.Minute)
	require.Eventually(b, func() bool {
		health := h.Reporter.APICache.Stats()
		//b.Logf("health: %v", health)
		return len(health) == 6
	}, 10*time.Second, 100*time.Millisecond)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, handler := range h.Handlers {
			_, err := handler.Endpoints().Query(ctx, req)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
