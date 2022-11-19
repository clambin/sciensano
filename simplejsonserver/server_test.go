package simplejsonserver

import (
	"context"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	demographicsClient := mockDemographics.NewFetcher(t)
	demographicsClient.On("GetByRegion").Return(map[string]int{})
	demographicsClient.On("GetByAgeBracket", mock.AnythingOfType("bracket.Bracket")).Return(0)
	//demographicsClient.On("Run", mock.AnythingOfType("*context.emptyCtx")).Return()

	h, err := New(0, demographicsClient)
	require.NoError(t, err)

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
