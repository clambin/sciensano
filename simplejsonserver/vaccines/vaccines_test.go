package vaccines_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestHandler_Failures(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines.TypeBatches).Return(nil, errors.New("fail"))

	ctx := context.Background()
	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccines.APIClient = f

	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: time.Now(),
	}}}}

	o := vaccinesHandler.OverviewHandler{Reporter: r}
	_, err := o.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	m := vaccinesHandler.ManufacturerHandler{Reporter: r}
	_, err = m.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	s := vaccinesHandler.StatsHandler{Reporter: r}
	_, err = s.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	d := vaccinesHandler.DelayHandler{Reporter: r}
	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	_, err = d.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, f)
}
