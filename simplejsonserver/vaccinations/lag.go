package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// LagHandler returns the time difference between partial and full COVID-19 vaccination
type LagHandler struct {
	reporter.Reporter
}

func (handler LagHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *LagHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var vaccinationsData *datasets.Dataset
	if vaccinationsData, err = handler.Reporter.GetVaccinations(); err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %w", err)
	}
	vaccinationsData.Accumulate()
	vaccinationsData.FilterByRange(req.Args.Range.From, req.Args.Range.To)
	timestamps, lag := buildLag(vaccinationsData)

	response = &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: timestamps},
			{Text: "lag", Data: lag},
		},
	}

	return
}

func buildLag(vaccinationsData *datasets.Dataset) (timestamps query.TimeColumn, lag query.NumberColumn) {
	vaccinationTimestamps := vaccinationsData.GetTimestamps()
	partial, _ := vaccinationsData.GetValues("partial")
	full, _ := vaccinationsData.GetValues("full")

	timestamps = make(query.TimeColumn, 0, len(vaccinationTimestamps))
	lag = make(query.NumberColumn, 0, len(vaccinationTimestamps))

	var firstDoseIndex int
	var lastSecondDose float64

	for index, value := range full {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		if value == 0 || value == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && partial[firstDoseIndex] < value {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinationTimestamps[index])
			lag = append(lag, vaccinationTimestamps[index].Sub(vaccinationTimestamps[firstDoseIndex]).Hours()/24)
		}

		lastSecondDose = value
	}

	return
}
