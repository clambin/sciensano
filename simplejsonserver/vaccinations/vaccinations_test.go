package vaccinations_test

import (
	"context"
	"github.com/clambin/go-metrics/client"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
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

func TestHandler(t *testing.T) {
	f := &mocks.Fetcher{}
	f.On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).Return(vaccinationTestData, nil)

	r := reporter.NewWithOptions(time.Hour, client.Options{})
	r.Vaccinations.APIClient = f
	h := vaccinations.Handler{Reporter: r}

	ctx := context.Background()
	req := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{To: timestamp.Add(24 * time.Hour)}}}}

	response, err := h.Endpoints().Query(ctx, req)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{
		Columns: []query.Column{
			{Text: "time", Data: query.TimeColumn{time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)}},
			{Text: "booster", Data: query.NumberColumn{1, 6}},
			{Text: "full", Data: query.NumberColumn{3, 10}},
			{Text: "partial", Data: query.NumberColumn{3, 8}},
		},
	}, response)

	mock.AssertExpectationsForObjects(t, f)
}

var (
	timestamp = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	vaccinationTestData = []apiclient.APIResponse{
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "25-34", Dose: sciensano.TypeVaccinationSingle, Count: 1},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Flanders", AgeGroup: "35-44", Dose: sciensano.TypeVaccinationBooster, Count: 1},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "35-44", Dose: sciensano.TypeVaccinationFull, Count: 2},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "Brussels", AgeGroup: "25-34", Dose: sciensano.TypeVaccinationPartial, Count: 2},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp}, Region: "", AgeGroup: "", Dose: sciensano.TypeVaccinationPartial, Count: 1},

		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: sciensano.TypeVaccinationFull, Count: 3},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: sciensano.TypeVaccinationSingle, Count: 4},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: sciensano.TypeVaccinationPartial, Count: 5},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(24 * time.Hour)}, Region: "Brussels", AgeGroup: "25-34", Dose: sciensano.TypeVaccinationBooster, Count: 5},

		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Flanders", AgeGroup: "25-34", Dose: sciensano.TypeVaccinationPartial, Count: 9},
		&sciensano.APIVaccinationsResponse{TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(48 * time.Hour)}, Region: "Brussels", AgeGroup: "35-44", Dose: sciensano.TypeVaccinationBooster, Count: 9},
	}
)

func buildBigResponse() (bigResponse []apiclient.APIResponse) {
	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 0; i < 2*365; i++ {
		for _, region := range []string{"Brussels", "Flanders", "Wallonia"} {
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     i + 100,
			})
			bigResponse = append(bigResponse, &sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: ts},
				Region:    region,
				Dose:      sciensano.TypeVaccinationFull,
				Count:     i,
			})
		}

		ts = ts.Add(24 * time.Hour)
	}

	return
}
