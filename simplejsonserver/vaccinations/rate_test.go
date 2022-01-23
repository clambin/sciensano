package vaccinations_test

import (
	"context"
	"github.com/clambin/sciensano/demographics/mocks"
	mockCache "github.com/clambin/sciensano/measurement/mocks"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/vaccinations"
	"github.com/clambin/simplejson"
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
		expected *simplejson.TableQueryResponse
	}{
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: simplejson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Flanders", Data: simplejson.TableQueryResponseNumberColumn{0, 0}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypeFull,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: simplejson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Flanders", Data: simplejson.TableQueryResponseNumberColumn{0.01, 0.04}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeRegion,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: simplejson.TableQueryResponseNumberColumn{0, 0}},
					{Text: "Flanders", Data: simplejson.TableQueryResponseNumberColumn{0.01, 0.01}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypePartial,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: simplejson.TableQueryResponseNumberColumn{0.02, 0.02}},
					{Text: "35-44", Data: simplejson.TableQueryResponseNumberColumn{0, 0.5}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeFull,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: simplejson.TableQueryResponseNumberColumn{0.01, 0.04}},
					{Text: "35-44", Data: simplejson.TableQueryResponseNumberColumn{0.2, 0.6}},
				},
			},
		},
		{
			Scope:           vaccinations.ScopeAge,
			VaccinationType: reporter.VaccinationTypeBooster,
			expected: &simplejson.TableQueryResponse{
				Columns: []simplejson.TableQueryResponseColumn{
					{Text: "timestamp", Data: simplejson.TableQueryResponseTimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: simplejson.TableQueryResponseNumberColumn{0, 0.05}},
					{Text: "35-44", Data: simplejson.TableQueryResponseNumberColumn{0.1, 0.1}},
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
	args := simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: timestamp.Add(24 * time.Hour)}}}

	for index, testCase := range testCases {
		h := vaccinations.RateHandler{
			Reporter:        client,
			VaccinationType: testCase.VaccinationType,
			Scope:           testCase.Scope,
			Demographics:    demographics,
		}

		response, err := h.Endpoints().TableQuery(ctx, &args)
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
	args := simplejson.TableQueryArgs{Args: simplejson.Args{Range: simplejson.Range{To: timestamp.Add(24 * time.Hour)}}}
	h := vaccinations.RateHandler{
		Reporter:        client,
		VaccinationType: reporter.VaccinationTypeBooster,
		Scope:           vaccinations.ScopeAge,
		Demographics:    demographics,
	}

	cache.On("Get", "Vaccinations").Return(nil, false).Once()
	_, err := h.Endpoints().TableQuery(ctx, &args)
	assert.Error(t, err)

	mock.AssertExpectationsForObjects(t, cache, demographics)
}

func BenchmarkRateHandler(b *testing.B) {
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
		_, err := h.Endpoints().TableQuery(context.Background(), &simplejson.TableQueryArgs{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
