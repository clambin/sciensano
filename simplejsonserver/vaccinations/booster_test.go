package vaccinations_test

import (
	"context"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBoosterHandler(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	r := reporter.NewWithOptions(time.Hour, httpclient.Options{})
	r.Vaccinations.APIClient = f

	//req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}
	var req query.Request
	ctx := context.Background()

	h := vaccinations.BoosterHandler{Reporter: r}

	response, err := h.Endpoints().Query(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "time", Data: query.TimeColumn{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC)}},
			{Text: "booster", Data: query.NumberColumn{1, 1, 1}},
			{Text: "booster2", Data: query.NumberColumn{0, 5, 5}},
			{Text: "booster3", Data: query.NumberColumn{0, 0, 9}},
		},
	}, response)
}

func BenchmarkBoosterHandler(b *testing.B) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(buildBigResponse(), nil)

	r := reporter.NewWithOptions(time.Hour, httpclient.Options{})
	r.Vaccinations.APIClient = f

	//req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}
	var req query.Request
	ctx := context.Background()

	h := vaccinations.BoosterHandler{Reporter: r}

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}
