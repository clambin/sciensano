package vaccines_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/reporter"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHandler_Failures(t *testing.T) {
	cache := &mockCache.Holder{}

	cache.On("Get", "Vaccinations").Return(nil, false)

	ctx := context.Background()
	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: time.Now(),
	}}}}

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	o := vaccinesHandler.OverviewHandler{Reporter: r}
	_, err := o.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	m := vaccinesHandler.ManufacturerHandler{Reporter: r}
	_, err = m.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	s := vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return(nil, false).Once()
	d := vaccinesHandler.DelayHandler{Reporter: r}
	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]apiclient.APIResponse{}, true).Once()
	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	cache.On("Get", "Vaccines").Return([]apiclient.APIResponse{}, true).Once()

	s = vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().Query(ctx, req)
	assert.Error(t, err)
}
