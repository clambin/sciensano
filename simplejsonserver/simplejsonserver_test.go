package simplejsonserver_test

import (
	"context"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/reporter"
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
	demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)
	demographicsClient.On("Run", mock.AnythingOfType("*context.emptyCtx")).Return()

	h := simplejsonserver.Server{
		Server:       simplejson.Server{Name: "sciensano"},
		Reporter:     reporter.New(time.Hour),
		Demographics: demographicsClient,
	}
	ctx := context.Background()
	h.Initialize(ctx)

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
}
