package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	"time"
)

type DelayHandler struct {
	Reporter *reporter.Client
}

var _ simplejson.Handler = &DelayHandler{}

func (d DelayHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: d.tableQuery}
}

func (d *DelayHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	batches, err := d.Reporter.Vaccines.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	var vaccinations *data.Table
	vaccinations, err = d.Reporter.Vaccinations.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccination call failed: %w", err)
	}

	return calculateDelay(vaccinations, batches).Filter(req.Args).CreateTableResponse(), nil
}

func calculateDelay(vaccinationsData, batches *data.Table) *data.Table {
	vaccinationsData = sumVaccinations(vaccinationsData).Accumulate()
	batches = batches.Accumulate()

	var timestamps []time.Time
	var delays []float64
	var batchIndex int

	for r := 0; r < vaccinationsData.Frame.Rows(); r++ {
		// find when we reached that number of vaccines
		vaccinations := vaccinationsData.Frame.At(1, r).(float64)
		for batchIndex < batches.Frame.Rows() && batches.Frame.At(1, batchIndex).(float64) < vaccinations {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			vaccTimestamp := vaccinationsData.Frame.At(0, r).(time.Time)
			batchTimestamp := batches.Frame.At(0, batchIndex-1).(time.Time)

			timestamps = append(timestamps, vaccTimestamp)
			delays = append(delays, vaccTimestamp.Sub(batchTimestamp).Hours()/24)
		}
	}

	return data.New(data.Column{Name: "time", Values: timestamps}, data.Column{Name: "delay", Values: delays})
}
