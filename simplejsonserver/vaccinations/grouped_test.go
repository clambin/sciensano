package vaccinations_test

import (
	"context"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson/v2/common"
	"github.com/clambin/simplejson/v2/query"
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
	args := query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}

	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	for index, testCase := range testCases {
		h := vaccinations.GroupedHandler{
			Reporter:        client,
			VaccinationType: testCase.VaccinationType,
			Scope:           testCase.Scope,
		}

		response, err := h.Endpoints().TableQuery(ctx, args)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.expected, response, index)
	}
}

func TestGroupedHandler_Failure(t *testing.T) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	cache.On("Get", "Vaccinations").Return(nil, false)

	h := vaccinations.GroupedHandler{
		Reporter:        client,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
	}

	_, err := h.Endpoints().TableQuery(context.Background(), query.Args{})
	require.Error(t, err)
}

func BenchmarkGroupedHandler(b *testing.B) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	content := buildBigResponse()
	cache.On("Get", "Vaccinations").Return(content, true)

	h := vaccinations.GroupedHandler{
		Reporter:        client,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
	}

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().TableQuery(context.Background(), query.Args{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
