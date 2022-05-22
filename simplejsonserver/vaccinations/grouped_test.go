package vaccinations_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGroupedHandler(t *testing.T) {
	var testCases = []struct {
		vaccinations.Scope
		reporter.VaccinationType
		expected *query.TableResponse
	}{
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: query.NumberColumn{1, 1}},
					{Text: "Brussels", Data: query.NumberColumn{2, 7}},
					{Text: "Flanders", Data: query.NumberColumn{0, 0}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypeFull,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{2, 6}},
					{Text: "Flanders", Data: query.NumberColumn{1, 4}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 5}},
					{Text: "Flanders", Data: query.NumberColumn{1, 1}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: query.NumberColumn{1, 1}},
					{Text: "25-34", Data: query.NumberColumn{2, 2}},
					{Text: "35-44", Data: query.NumberColumn{0, 5}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeFull,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{1, 4}},
					{Text: "35-44", Data: query.NumberColumn{2, 6}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0, 5}},
					{Text: "35-44", Data: query.NumberColumn{1, 1}},
				},
			},
		},
	}

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	cache := &mockCache.Holder{}
	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	for index, testCase := range testCases {
		h := vaccinations.GroupedHandler{
			Reporter:        r,
			VaccinationType: testCase.VaccinationType,
			Scope:           testCase.Scope,
		}

		response, err := h.Endpoints().Query(ctx, req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.expected, response, index)
	}
}

func TestGroupedHandler_Failure(t *testing.T) {
	cache := &mockCache.Holder{}
	cache.On("Get", "Vaccinations").Return(nil, false)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	h := vaccinations.GroupedHandler{
		Reporter:        r,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
	}

	_, err := h.Endpoints().Query(context.Background(), query.Request{})
	require.Error(t, err)
}

func BenchmarkVaccinationsGroupedHandler(b *testing.B) {
	cache := &mockCache.Holder{}
	content := buildBigResponse()
	cache.On("Get", "Vaccinations").Return(content, true)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.APICache = cache

	h := vaccinations.GroupedHandler{
		Reporter:        r,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
	}

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), query.Request{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
