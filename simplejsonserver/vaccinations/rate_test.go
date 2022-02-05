package vaccinations_test

import (
	"context"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
	"github.com/clambin/sciensano/demographics/mocks"
	"github.com/clambin/sciensano/reporter"
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
		reporter.VaccinationType
		expected *query.TableResponse
	}{
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
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
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
					{Text: "Flanders", Data: query.NumberColumn{0.01, 0.04}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 0}},
					{Text: "Flanders", Data: query.NumberColumn{0.01, 0.01}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0.02, 0.02}},
					{Text: "35-44", Data: query.NumberColumn{0, 0.5}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeFull,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0.01, 0.04}},
					{Text: "35-44", Data: query.NumberColumn{0.2, 0.6}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "timestamp", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0, 0.05}},
					{Text: "35-44", Data: query.NumberColumn{0.1, 0.1}},
				},
			},
		},
	}

	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache
	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	demographics := &mocks.Demographics{}
	demographics.
		On("GetRegionFigures").
		Return(map[string]int{
			"Flanders": 100,
			//"Brussels": 10,
		})

	demographics.
		On("GetAgeGroupFigures").
		Return(map[string]int{
			"25-34": 100,
			"35-44": 10,
		})

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	for index, testCase := range testCases {
		h := vaccinations.RateHandler{
			Reporter:        client,
			VaccinationType: testCase.VaccinationType,
			Scope:           testCase.Scope,
			Demographics:    demographics,
		}

		response, err := h.Endpoints().Query(ctx, req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.expected, response, index)
	}

	mock.AssertExpectationsForObjects(t, cache, demographics)
}

func TestRateHandler_Failure(t *testing.T) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	demographics := &mocks.Demographics{}

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}
	h := vaccinations.RateHandler{
		Reporter:        client,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
		Demographics:    demographics,
	}

	cache.On("Get", "Vaccinations").Return(nil, false).Once()
	_, err := h.Endpoints().Query(ctx, req)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache, demographics)
}

func BenchmarkVaccinationsRateHandler(b *testing.B) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	demographics := &mocks.Demographics{}
	demographics.On("GetRegionFigures").Return(map[string]int{"Brussels": 1, "Flanders": 6, "Wallonia": 4})

	content := buildBigResponse()
	cache.On("Get", "Vaccinations").Return(content, true)

	h := vaccinations.RateHandler{
		Reporter:        client,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeRegion,
		Demographics:    demographics,
	}

	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), query.Request{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
