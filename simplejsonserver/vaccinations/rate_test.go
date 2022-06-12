package vaccinations_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/demographics/bracket"
	mockDemographics "github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/reporter"
	vaccinations2 "github.com/clambin/sciensano/reporter/vaccinations"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRateHandler(t *testing.T) {
	var testCases = []struct {
		vaccinations.Scope
		Type     int
		expected *query.TableResponse
	}{
		{
			Scope: vaccinations.ScopeRegion,
			Type:  vaccinations2.TypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
					{Text: "Flanders", Data: query.NumberColumn{0, 0}},
				},
			},
		},
		{
			Scope: vaccinations.ScopeRegion,
			Type:  vaccinations2.TypeFull,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
					{Text: "Flanders", Data: query.NumberColumn{0.01, 0.04}},
				},
			},
		},
		{
			Scope: vaccinations.ScopeRegion,
			Type:  vaccinations2.TypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
					{Text: "Flanders", Data: query.NumberColumn{0.01, 0.01}},
				},
			},
		},
		{
			Scope: vaccinations.ScopeAge,
			Type:  vaccinations2.TypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0.02, 0.02}},
					{Text: "35-44", Data: query.NumberColumn{0, 0.5}},
				},
			},
		},
		{
			Scope: vaccinations.ScopeAge,
			Type:  vaccinations2.TypeFull,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0.01, 0.04}},
					{Text: "35-44", Data: query.NumberColumn{0.2, 0.6}},
				},
			},
		},
		{
			Scope: vaccinations.ScopeAge,
			Type:  vaccinations2.TypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0, 0.05}},
					{Text: "35-44", Data: query.NumberColumn{0.1, 0.1}},
				},
			},
		},
	}

	cache := &mockCache.Holder{}
	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APICache = cache

	demographicsClient := &mockDemographics.Fetcher{}
	demographicsClient.
		On("GetByRegion").
		Return(map[string]int{
			"Flanders": 100,
			//"Brussels": 10,
		})

	demographicsClient.
		On("GetByAgeBracket", bracket.Bracket{Low: 25, High: 34}).
		Return(100)

	demographicsClient.
		On("GetByAgeBracket", bracket.Bracket{Low: 35, High: 44}).
		Return(10)

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	for index, testCase := range testCases {
		h := vaccinations.RateHandler{
			Reporter: r,
			Type:     testCase.Type,
			Scope:    testCase.Scope,
			Fetcher:  demographicsClient,
		}

		response, err := h.Endpoints().Query(ctx, req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.expected, response, index)
	}

	mock.AssertExpectationsForObjects(t, cache, demographicsClient)
}

func TestRateHandler_Failure(t *testing.T) {
	cache := &mockCache.Holder{}
	cache.On("Get", "Vaccinations").Return(nil, false).Once()

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APICache = cache

	demographicsClient := &mockDemographics.Fetcher{}

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}
	h := vaccinations.RateHandler{
		Reporter: r,
		Type:     vaccinations2.TypeBooster,
		Scope:    vaccinations.ScopeAge,
		Fetcher:  demographicsClient,
	}

	_, err := h.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache, demographicsClient)
}

func BenchmarkVaccinationsRateHandler(b *testing.B) {
	b.StopTimer()
	cache := &mockCache.Holder{}
	content := buildBigResponse()
	cache.On("Get", "Vaccinations").Return(content, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APICache = cache

	demographicsClient := &mockDemographics.Fetcher{}
	demographicsClient.On("GetByRegion").Return(map[string]int{"Brussels": 1, "Flanders": 6, "Wallonia": 4})

	h := vaccinations.RateHandler{
		Reporter: r,
		Type:     vaccinations2.TypeBooster,
		Scope:    vaccinations.ScopeRegion,
		Fetcher:  demographicsClient,
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), query.Request{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
