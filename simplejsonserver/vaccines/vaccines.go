package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type OverviewHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &OverviewHandler{}

func (o OverviewHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: o.tableQuery}
}

func (o *OverviewHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var batches *datasets.Dataset
	batches, err = o.Reporter.GetVaccines()
	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}
	batches = batches.Copy()
	batches.Accumulate()
	return responder.GenerateTableQueryResponse(batches, req.Args), nil
}

type ManufacturerHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &ManufacturerHandler{}

func (m ManufacturerHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: m.tableQuery}
}

func (m *ManufacturerHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var batches *datasets.Dataset
	batches, err = m.Reporter.GetVaccinesByManufacturer()
	if err != nil {
		return nil, fmt.Errorf("vaccine manufacturer call failed: %w", err)
	}
	batches.Accumulate()
	return responder.GenerateTableQueryResponse(batches, req.Args), nil
}

type StatsHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &StatsHandler{}

func (s StatsHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: s.tableQuery}
}

func (s *StatsHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var batches *datasets.Dataset
	if batches, err = s.Reporter.GetVaccines(); err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	var vaccinations *datasets.Dataset
	if vaccinations, err = s.Reporter.GetVaccinations(); err != nil {
		return nil, fmt.Errorf("vaccinations call failed: %w", err)
	}

	calculateVaccineReserve(vaccinations, batches)
	vaccinations.FilterByRange(req.Range.From, req.Range.To)

	vaccinationValues, _ := vaccinations.GetValues("sum")
	reserveValues, _ := vaccinations.GetValues("reserve")

	response = &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: query.TimeColumn(vaccinations.GetTimestamps())},
			{Text: "vaccinations", Data: query.NumberColumn(vaccinationValues)},
			{Text: "reserve", Data: query.NumberColumn(reserveValues)},
		},
	}
	return
}

func calculateVaccineReserve(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) {
	// summarize all vaccinations
	vaccinationsData.AddColumn("sum", func(values map[string]float64) (sum float64) {
		for _, value := range values {
			sum += value
		}
		return
	})

	// Add the batches to the vaccination dataset
	receivedTimestamps := batches.GetTimestamps()
	received, _ := batches.GetValues("total")
	for index, entry := range received {
		vaccinationsData.Add(receivedTimestamps[index], "batch", entry)
	}

	// accumulate the numbers
	vaccinationsData.Accumulate()

	// calculate reserve
	vaccinationsData.AddColumn("reserve", func(values map[string]float64) float64 {
		return values["batch"] - values["sum"]
	})
}

type DelayHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &DelayHandler{}

func (d DelayHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: d.tableQuery}
}

func (d *DelayHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var batches *datasets.Dataset
	batches, err = d.Reporter.GetVaccines()

	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	batches.Accumulate()

	var vaccinations *datasets.Dataset
	vaccinations, err = d.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("vaccination call failed: %w", err)
	}

	vaccinations.Accumulate()
	vaccinations.FilterByRange(req.Range.From, req.Range.To)

	timestampColumn, timeColumn := calculateVaccineDelay(vaccinations, batches)

	response = &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: timestampColumn},
			{Text: "time", Data: timeColumn},
		},
	}
	return
}

func calculateVaccineDelay(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) (timestamps query.TimeColumn, delays query.NumberColumn) {
	vaccinationsData.AddColumn("sum", func(values map[string]float64) (sum float64) {
		for _, value := range values {
			sum += value
		}
		return
	})

	batchCount := batches.Size()
	batchTimestamps := batches.GetTimestamps()
	batchData, _ := batches.GetValues("total")
	batchIndex := 0

	vaccinationTimestamps := vaccinationsData.GetTimestamps()
	vaccinations, _ := vaccinationsData.GetValues("sum")

	for index, sum := range vaccinations {
		// find when we reached that number of vaccines
		for batchIndex < batchCount && batchData[batchIndex] < sum {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			timestamps = append(timestamps, vaccinationTimestamps[index])
			delays = append(delays, vaccinationTimestamps[index].Sub(batchTimestamps[batchIndex-1]).Hours()/24)
		}
	}
	return
}
