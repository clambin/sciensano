package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson/v2"
	"github.com/clambin/simplejson/v2/query"
)

// LagHandler returns the time difference between partial and full COVID-19 vaccination
type LagHandler struct {
	reporter.Reporter
}

func (handler LagHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: handler.tableQuery}
}

func (handler *LagHandler) tableQuery(_ context.Context, _ query.Args) (response *query.TableResponse, err error) {
	var vaccinationsData *datasets.Dataset

	vaccinationsData, err = handler.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %w", err)
	}

	vaccinationsData.Accumulate()
	// vaccinationsData.ApplyRange(args.Range.From, args.Range.To)
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
	var firstDoseIndex, lastSecondDose int

	/*
		if len(vaccinationsData.Groups) == 0 {
			log.Warning("no vaccination data to calculate lag")
			return
		}
	*/

	// run through all vaccinations
	for index := range vaccinationsData.Timestamps {
		// we only measure lag when there is actually a second dose
		// we don't report when the 2nd dose doesn't change
		full := int(vaccinationsData.Groups[1].Values[index])
		if full == 0 || full == lastSecondDose {
			continue
		}

		// find the time when we reached the number of first Doses that equals (or higher) the current Second Dose number
		for firstDoseIndex <= index && int(vaccinationsData.Groups[0].Values[firstDoseIndex]) < full {
			firstDoseIndex++
		}

		// if we found it, add it to the columns
		if firstDoseIndex <= index {
			timestamps = append(timestamps, vaccinationsData.Timestamps[index])
			lag = append(lag, vaccinationsData.Timestamps[index].Sub(vaccinationsData.Timestamps[firstDoseIndex]).Hours()/24)
		}

		lastSecondDose = full
	}

	return
}
