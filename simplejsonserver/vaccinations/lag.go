package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	"time"
)

// LagHandler returns the time difference between partial and full COVID-19 vaccination
type LagHandler struct {
	Reporter *reporter.Client
}

func (handler LagHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *LagHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	vaccinationsData, err := handler.Reporter.Vaccinations.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to determine vaccination lag: %w", err)
	}

	vaccinationsData = vaccinationsData.Accumulate()
	lag := buildLag(vaccinationsData)

	return lag.Filter(req.Args).CreateTableResponse(), nil
}

func buildLag(input *data.Table) (output *data.Table) {
	vaccinationTimestamps := input.GetTimestamps()
	partial, _ := input.GetFloatValues("partial")
	full, _ := input.GetFloatValues("full")

	timestamps := make([]time.Time, 0, len(vaccinationTimestamps))
	lag := make([]float64, 0, len(vaccinationTimestamps))

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

	return data.New(data.Column{Name: "time", Values: timestamps}, data.Column{Name: "lag", Values: lag})
}
