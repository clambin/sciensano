package apihandler_test

import (
	"context"
	"fmt"
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apiclient"
	"github.com/clambin/sciensano/apihandler"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/vaccines"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestAPIHandler_Vaccines(t *testing.T) {
	endDate := time.Date(2021, 06, 01, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.vaccinesClient.On("GetBatches", mock.Anything).Return([]*vaccines.Batch{
		{
			Date:   vaccines.Time{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 6, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 7, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
	}, nil)

	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vaccines", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafana_json.TableQueryResponseTimeColumn:
			require.Equal(t, "timestamp", column.Text)
			require.NotEmpty(t, data)
			lastDate := data[len(data)-1]
			assert.Equal(t, 2021, lastDate.Year())
			assert.Equal(t, time.Month(6), lastDate.Month())
			assert.Equal(t, 1, lastDate.Day())
		case grafana_json.TableQueryResponseNumberColumn:
			switch column.Text {
			case "vaccines":
				require.NotZero(t, len(data))
				assert.Equal(t, 600.0, data[len(data)-1])
			}
		}
	}

	stack.vaccinesClient.AssertExpectations(t)
}

func TestAPIHandler_Vaccines_Stats(t *testing.T) {
	endDate := time.Date(2021, 05, 1, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.vaccinesClient.On("GetBatches", mock.Anything).Return([]*vaccines.Batch{
		{
			Date:   vaccines.Time{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
	}, nil)

	stack.sciensanoClient.On("GetVaccinations", mock.Anything).Return([]*apiclient.APIVaccinationsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     200,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     50,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     300,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     320,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     155,
		},
	}, nil)

	response, err := stack.apiHandler.Endpoints().TableQuery(context.Background(), "vaccines-stats", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafana_json.TableQueryResponseTimeColumn:
			require.Equal(t, "timestamp", column.Text)
			require.NotZero(t, len(data))
			lastDate := data[len(data)-1]
			assert.Equal(t, 2021, lastDate.Year())
			assert.Equal(t, time.Month(4), lastDate.Month())
			assert.Equal(t, 1, lastDate.Day())
		case grafana_json.TableQueryResponseNumberColumn:
			switch column.Text {
			case "vaccinations":
				require.NotZero(t, len(data))
				assert.Equal(t, 1175.0, data[len(data)-1])
			case "reserve":
				require.NotZero(t, len(data))
				assert.Equal(t, -375.0, data[len(data)-1])
			}
		}
	}
	mock.AssertExpectationsForObjects(t, &stack.vaccinesClient, &stack.sciensanoClient)
}

func TestAPIHandler_Vaccines_Time(t *testing.T) {
	endDate := time.Date(2021, 06, 1, 0, 0, 0, 0, time.UTC)
	request := &grafana_json.TableQueryArgs{
		CommonQueryArgs: grafana_json.CommonQueryArgs{
			Range: grafana_json.QueryRequestRange{To: endDate},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	stack := createStack(ctx)
	defer stack.Destroy()

	stack.vaccinesClient.On("GetBatches", mock.Anything).Return([]*vaccines.Batch{
		{
			Date:   vaccines.Time{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
		{
			Date:   vaccines.Time{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Amount: 200,
		},
	}, nil)

	stack.sciensanoClient.On("GetVaccinations", mock.Anything).Return([]*apiclient.APIVaccinationsResponse{
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     200,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 16, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     50,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     300,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 3, 31, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     150,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "A",
			Count:     300,
		},
		{
			TimeStamp: apiclient.TimeStamp{Time: time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)},
			Dose:      "B",
			Count:     150,
		},
	}, nil)

	response, err := stack.apiHandler.Endpoints().TableQuery(ctx, "vaccines-time", request)
	require.NoError(t, err)

	for _, column := range response.Columns {
		switch data := column.Data.(type) {
		case grafana_json.TableQueryResponseTimeColumn:
			assert.Equal(t, "timestamp", column.Text)
			require.NotZero(t, len(data))
			lastDate := data[len(data)-1]
			assert.Equal(t, 2021, lastDate.Year())
			assert.Equal(t, time.Month(4), lastDate.Month())
			assert.Equal(t, 1, lastDate.Day())
		case grafana_json.TableQueryResponseNumberColumn:
			switch column.Text {
			case "time":
				require.NotZero(t, len(data))
				assert.Equal(t, 74.0, data[0])
			}
		}
	}

	mock.AssertExpectationsForObjects(t, &stack.vaccinesClient, &stack.sciensanoClient)
}

func TestVaccineDelay(t *testing.T) {
	vaccinations := []sciensano.Vaccination{{
		Timestamp: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
		Partial:   10,
		Full:      0,
	}, {
		Timestamp: time.Date(2021, 01, 15, 0, 0, 0, 0, time.UTC),
		Partial:   15,
		Full:      1,
	}, {
		Timestamp: time.Date(2021, 02, 1, 0, 0, 0, 0, time.UTC),
		Partial:   15,
		Full:      4,
	}, {
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Partial:   25,
		Full:      5,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Partial:   35,
		Full:      10,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Partial:   35,
		Full:      15,
	}}

	batches := []*vaccines.Batch{{
		Date:   vaccines.Time{Time: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 20,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 02, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 40,
	}, {
		Date:   vaccines.Time{Time: time.Date(2021, 03, 01, 0, 0, 0, 0, time.UTC)},
		Amount: 50,
	}}

	expected := []struct {
		Timestamp time.Time
		Value     float64
	}{{
		Timestamp: time.Date(2021, 02, 15, 0, 0, 0, 0, time.UTC),
		Value:     45,
	}, {
		Timestamp: time.Date(2021, 03, 1, 0, 0, 0, 0, time.UTC),
		Value:     28,
	}, {
		Timestamp: time.Date(2021, 03, 15, 0, 0, 0, 0, time.UTC),
		Value:     42,
	}}

	timestamps, delays := apihandler.CalculateVaccineDelay(vaccinations, batches)

	if assert.Equal(t, len(expected), len(timestamps)) && assert.Equal(t, len(expected), len(delays)) {
		for index, entry := range expected {
			assert.Equal(t, entry.Timestamp, timestamps[index], fmt.Sprintf("index: %d", index))
			assert.Equal(t, entry.Value, delays[index], fmt.Sprintf("index: %d", index))
		}

	}
}
