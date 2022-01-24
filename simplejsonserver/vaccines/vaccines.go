package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v2"
	"github.com/clambin/simplejson/v2/query"
)

type OverviewHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &OverviewHandler{}

func (o OverviewHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: o.tableQuery}
}

func (o *OverviewHandler) tableQuery(_ context.Context, args query.Args) (response *query.TableResponse, err error) {
	var batches *datasets.Dataset
	batches, err = o.Reporter.GetVaccines()
	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}
	batches.Accumulate()
	return responder.GenerateTableQueryResponse(batches, args), nil
}

type ManufacturerHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &ManufacturerHandler{}

func (m ManufacturerHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: m.tableQuery}
}

func (m *ManufacturerHandler) tableQuery(_ context.Context, args query.Args) (response *query.TableResponse, err error) {
	var batches *datasets.Dataset
	batches, err = m.Reporter.GetVaccinesByManufacturer()
	if err != nil {
		return nil, fmt.Errorf("vaccine manufacturer call failed: %w", err)
	}
	batches.Accumulate()
	return responder.GenerateTableQueryResponse(batches, args), nil
}

type StatsHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &StatsHandler{}

func (s StatsHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: s.tableQuery}
}

func (s *StatsHandler) tableQuery(_ context.Context, args query.Args) (response *query.TableResponse, err error) {
	var batches *datasets.Dataset
	batches, err = s.Reporter.GetVaccines()

	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	batches.Accumulate()

	var vaccinations *datasets.Dataset
	vaccinations, err = s.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("vaccinations call failed: %w", err)
	}

	summarizeVaccinations(vaccinations)
	vaccinations.Accumulate()

	vaccinations.Groups = append(vaccinations.Groups, datasets.DatasetGroup{
		Name:   "reserve",
		Values: calculateVaccineReserve(vaccinations, batches),
	})

	batches.ApplyRange(args.Range.From, args.Range.To)
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

	response = &query.TableResponse{
		Columns: []query.Column{
			{Text: "timestamp", Data: query.TimeColumn(vaccinations.Timestamps)},
			{Text: "vaccinations", Data: query.NumberColumn(vaccinations.Groups[0].Values)},
			{Text: "reserve", Data: query.NumberColumn(vaccinations.Groups[1].Values)},
		},
	}
	return
}

func summarizeVaccinations(vaccinationsData *datasets.Dataset) {
	for index := range vaccinationsData.Timestamps {
		total := 0.0
		for _, group := range vaccinationsData.Groups {
			total += group.Values[index]
		}
		vaccinationsData.Groups[0].Values[index] = total
	}
	vaccinationsData.Groups = vaccinationsData.Groups[:1]
}

func calculateVaccineReserve(vaccinationsData *datasets.Dataset, batches *datasets.Dataset) (reserve []float64) {
	batchIndex := 0
	lastBatch := 0.0

	for index, timestamp := range vaccinationsData.Timestamps {
		// find the last time we received vaccines
		for batchIndex < len(batches.Timestamps) && batches.Timestamps[batchIndex].After(timestamp) == false {
			// how many vaccines have we received so far?
			lastBatch = batches.Groups[0].Values[batchIndex]
			batchIndex++
		}

		// add it to the list
		reserve = append(reserve, lastBatch-vaccinationsData.Groups[0].Values[index])
	}

	return
}

type DelayHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &DelayHandler{}

func (d DelayHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: d.tableQuery}
}

func (d *DelayHandler) tableQuery(_ context.Context, args query.Args) (response *query.TableResponse, err error) {
	var batches *datasets.Dataset
	batches, err = d.Reporter.GetVaccines()

	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	batches.Accumulate()
	batches.ApplyRange(args.Range.From, args.Range.To)

	var vaccinations *datasets.Dataset
	vaccinations, err = d.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("vaccination call failed: %w", err)
	}

	vaccinations.Accumulate()
	vaccinations.ApplyRange(args.Range.From, args.Range.To)

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
	batchIndex := 0
	batchCount := len(batches.Timestamps)

	for index, timestamp := range vaccinationsData.Timestamps {
		// how many vaccines did we need to perform this many vaccinations?
		vaccinesNeeded := vaccinationsData.Groups[0].Values[index] + vaccinationsData.Groups[1].Values[index]

		// find when we reached that number of vaccines
		for batchIndex < batchCount && batches.Groups[0].Values[batchIndex] < vaccinesNeeded {
			batchIndex++
		}

		// we depleted the *previous* batch. report the time difference between now and when we received that batch
		if batchIndex > 0 {
			timestamps = append(timestamps, timestamp)
			delays = append(delays, timestamp.Sub(batches.Timestamps[batchIndex-1]).Hours()/24)
		}
	}
	return
}
