package vaccinations_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/measurement"
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

func TestHandler(t *testing.T) {
	cache := &mockCache.Holder{}
	client := reporter.New(time.Hour)
	client.APICache = cache

	h := vaccinations.Handler{Reporter: client}

	cache.On("Get", "Vaccinations").Return(nil, false).Once()

	ctx := context.Background()
	args := query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}

	_, err := h.Endpoints().TableQuery(ctx, args)
	assert.Error(t, err)

	cache.On("Get", "Vaccinations").Return(vaccinationTestData, true)

	response, err := h.Endpoints().TableQuery(ctx, args)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: query.TimeColumn{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)}},
			{Text: "partial", Data: query.NumberColumn{3, 8}},
			{Text: "full", Data: query.NumberColumn{3, 10}},
			{Text: "booster", Data: query.NumberColumn{1, 6}},
		},
	}, response)
}

var (
	timestamp = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	vaccinationTestData = []measurement.Measurement{
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "25-34", Dose: "C", Count: 1},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "35-44", Dose: "E", Count: 1},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "35-44", Dose: "B", Count: 2},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "25-34", Dose: "A", Count: 2},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "", AgeGroup: "", Dose: "A", Count: 1},

		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "B", Count: 3},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "C", Count: 4},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "A", Count: 5},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "25-34", Dose: "E", Count: 5},

		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: "A", Count: 9},
		&sciensano.APIVaccinationsResponseEntry{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: "E", Count: 9},
	}
)

func buildBigResponse() (bigResponse []measurement.Measurement) {
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "A",
				Count:     i + 100,
			})
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponseEntry{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      "B",
				Count:     i,
			})
		}

		ts = ts.Add(24 * time.Hour)
	}

	return
}
