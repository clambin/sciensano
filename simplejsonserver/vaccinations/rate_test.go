package vaccinations_test

import (
	"context"
	"errors"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
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

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f

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

	mock.AssertExpectationsForObjects(t, f, demographicsClient)
}

func TestRateHandler_Failure(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f

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

	mock.AssertExpectationsForObjects(t, f, demographicsClient)
}

func BenchmarkVaccinationsRateHandler(b *testing.B) {
	content := buildBigResponse()

	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(content, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f

	demographicsClient := &mockDemographics.Fetcher{}
	demographicsClient.On("GetByRegion").Return(map[string]int{"Brussels": 1, "Flanders": 6, "Wallonia": 4})

	h := vaccinations.RateHandler{
		Reporter: r,
		Type:     vaccinations2.TypeBooster,
		Scope:    vaccinations.ScopeRegion,
		Fetcher:  demographicsClient,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), query.Request{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
