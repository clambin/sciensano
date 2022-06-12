package vaccines_test

import (
	"context"
	"github.com/clambin/sciensano/apiclient"
	mockCache "github.com/clambin/sciensano/apiclient/cache/mocks"
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
	vaccineCache := &mockCache.Holder{}
	timestamp := time.Date(2022, 1, 26, 0, 0, 0, 0, time.UTC)
	vaccineCache.
		On("Get", "Vaccines").
		Return([]apiclient.APIResponse{
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-7 * 24 * time.Hour)},
				Amount: 100,
			},
			&vaccines.APIBatchResponse{
				Date:   vaccines.Timestamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Amount: 50,
			},
		}, true)

	vaccinationCache := &mockCache.Holder{}
	vaccinationCache.
		On("Get", "Vaccinations").
		Return([]apiclient.APIResponse{
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-6 * 24 * time.Hour)},
				Dose:      "A",
				Count:     50,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-5 * 24 * time.Hour)},
				Dose:      "A",
				Count:     25,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-4 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-3 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-2 * 24 * time.Hour)},
				Dose:      "A",
				Count:     40,
			},
			&sciensano.APIVaccinationsResponse{
				TimeStamp: sciensano.TimeStamp{Time: timestamp.Add(-1 * 24 * time.Hour)},
				Dose:      "A",
				Count:     15,
			},
		}, true)

	r := reporter.New(time.Hour)
	r.Vaccines.APICache = vaccineCache
	r.Vaccinations.APICache = vaccinationCache
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

	mock.AssertExpectationsForObjects(t, vaccineCache, vaccinationCache)
}
