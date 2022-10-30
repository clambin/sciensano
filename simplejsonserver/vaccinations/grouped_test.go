package vaccinations_test

import (
	"context"
	"errors"
	"github.com/clambin/httpclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
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

func TestGroupedHandler(t *testing.T) {
	var testCases = []struct {
		vaccinations.Scope
		Type       int
		Accumulate bool
		expected   *query.TableResponse
	}{
		{
			Scope:      vaccinations.ByRegion,
			Type:       vaccinations2.TypePartial,
			Accumulate: true,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: query.NumberColumn{1, 1}},
					{Text: "Brussels", Data: query.NumberColumn{2, 7}},
					{Text: "Flanders", Data: query.NumberColumn{0, 0}},
				},
			},
		},
		{
			Scope:      vaccinations.ByRegion,
			Type:       vaccinations2.TypeFull,
			Accumulate: true,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{2, 6}},
					{Text: "Flanders", Data: query.NumberColumn{1, 4}},
				},
			},
		},
		{
			Scope:      vaccinations.ByRegion,
			Type:       vaccinations2.TypeBooster,
			Accumulate: true,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 5}},
					{Text: "Flanders", Data: query.NumberColumn{1, 1}},
				},
			},
		},
		{
			Scope:      vaccinations.ByRegion,
			Type:       vaccinations2.TypeBooster,
			Accumulate: false,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "Brussels", Data: query.NumberColumn{0, 5}},
					{Text: "Flanders", Data: query.NumberColumn{1, 0}},
				},
			},
		},
		{
			Scope:      vaccinations.ByAge,
			Type:       vaccinations2.TypePartial,
			Accumulate: true,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "(unknown)", Data: query.NumberColumn{1, 1}},
					{Text: "25-34", Data: query.NumberColumn{2, 2}},
					{Text: "35-44", Data: query.NumberColumn{0, 5}},
				},
			},
		},
		{
			Scope:      vaccinations.ByAge,
			Type:       vaccinations2.TypeFull,
			Accumulate: true,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{1, 4}},
					{Text: "35-44", Data: query.NumberColumn{2, 6}},
				},
			},
		},
		{
			Scope:      vaccinations.ByAge,
			Type:       vaccinations2.TypeBooster,
			Accumulate: false,
			expected: &query.TableResponse{
				Columns: []query.Column{
					{Text: "time", Data: query.TimeColumn{timestamp, timestamp.Add(24 * time.Hour)}},
					{Text: "25-34", Data: query.NumberColumn{0, 5}},
					{Text: "35-44", Data: query.NumberColumn{1, 0}},
				},
			},
		},
	}

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	f := mocks.NewFetcher(t)
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	r := reporter.NewWithOptions(time.Hour, httpclient.Options{})
	r.Vaccinations.APIClient = f

	for index, testCase := range testCases {
		h := vaccinations.GroupedHandler{
			Reporter:   r,
			Type:       testCase.Type,
			Scope:      testCase.Scope,
			Accumulate: testCase.Accumulate,
		}

		response, err := h.Endpoints().Query(ctx, req)
		require.NoError(t, err, index)
		assert.Equal(t, testCase.expected, response, index)
	}
}

func TestGroupedHandler_Failure(t *testing.T) {
	f := mocks.NewFetcher(t)
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(nil, errors.New("fail"))

	r := reporter.NewWithOptions(time.Hour, httpclient.Options{})
	r.Vaccinations.APIClient = f

	h := vaccinations.GroupedHandler{
		Reporter: r,
		Type:     vaccinations2.TypeBooster,
		Scope:    vaccinations.ByAge,
	}

	_, err := h.Endpoints().Query(context.Background(), query.Request{})
	require.Error(t, err)
}

func BenchmarkVaccinationsGroupedHandler(b *testing.B) {
	content := buildBigResponse()
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(content, nil)

	r := reporter.NewWithOptions(time.Hour, httpclient.Options{})
	r.Vaccinations.APIClient = f

	h := vaccinations.GroupedHandler{
		Reporter: r,
		Type:     vaccinations2.TypePartial,
		Scope:    vaccinations.ByAge,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := h.Endpoints().Query(context.Background(), query.Request{})
		if err != nil {
			b.Fatal(err)
		}
	}
}
