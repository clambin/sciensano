package vaccines_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apiclient/fetcher/mocks"
	"github.com/clambin/sciensano/apiclient/sciensano"
	"github.com/clambin/sciensano/apiclient/vaccines"
	"github.com/clambin/sciensano/reporter"
	vaccinesHandler "github.com/clambin/sciensano/simplejsonserver/vaccines"
	"github.com/clambin/simplejson/v3/common"
	"github.com/clambin/simplejson/v3/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHandler_TableQuery_VaccinesTime(t *testing.T) {
	vaccineClient := &mocks.Fetcher{}
	timestamp := time.Date(2022, 1, 26, 0, 0, 0, 0, time.UTC)
	vaccineClient.
		On("Fetch", mock.AnythingOfType("*context.emptyCtx"), vaccines.TypeBatches).
		Return([]apiclient.APIResponse{
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, nil)

	vaccinationClient := &mocks.Fetcher{}
	vaccinationClient.
		On("Fetch", mock.AnythingOfType("*context.emptyCtx"), sciensano.TypeVaccinations).
		Return([]apiclient.APIResponse{
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-6 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     50,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-5 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     25,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-4 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-3 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     40,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-1 * 24 * time.Hour)},
				Dose:      sciensano.TypeVaccinationPartial,
				Count:     15,
			},
		}, nil)

	r := reporter.New(time.Hour)
	r.Vaccines.APIClient = vaccineClient
	r.Vaccinations.APIClient = vaccinationClient
	h := vaccinesHandler.DelayHandler{Reporter: r}

	request := query.Request{Args: query.Args{Args: common.Args{Range: common.Range{
		To: timestamp,
	}}}}

	response, err := h.Endpoints().Query(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, &query.TableResponse{Columns: []query.Column{
		{Text: "time", Data: query.TimeColumn{timestamp.Add(-3 * 24 * time.Hour), timestamp.Add(-2 * 24 * time.Hour), timestamp.Add(-24 * time.Hour)}},
		{Text: "delay", Data: query.NumberColumn{4, 5, 1}},
	}}, response)

	mock.AssertExpectationsForObjects(t, vaccineClient, vaccinationClient)
}
