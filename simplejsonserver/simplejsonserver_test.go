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
	"net/http"
	"net/http/httptest"
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
		"boosters",
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
		"vacc-age-rate-full",
		"vacc-age-rate-partial",
		"vacc-manufacturer",
		"vacc-region-booster",
		"vacc-region-rate-full",
		"vacc-region-rate-partial",
		"vaccinations",
	}, h.Server.Targets())

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: time.Now()}}}}

	wg := sync.WaitGroup{}
	wg.Add(len(h.Server.Handlers))
	var count int
	for target, handler := range h.Server.Handlers {
		count++
		go func(handler simplejson.Handler, target string, counter int) {
			t.Logf("%2d: %s started", counter, target)
			_, err := handler.Endpoints().Query(ctx, req)
			t.Logf("%2d: %s done. err: %v", counter, target, err)
			wg.Done()
			assert.NoError(t, err, target, target)
		}(handler, target, count)
	}
	wg.Wait()

	b := httptest.NewRecorder()
	h.Health(b, nil)
	assert.Equal(t, http.StatusOK, b.Code)
	assert.Contains(t, b.Body.String(), `"Handlers": `)
	assert.Contains(t, b.Body.String(), `"ReporterCache": `)
}
