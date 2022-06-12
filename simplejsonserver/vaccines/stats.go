package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
	"time"
)

type StatsHandler struct {
	Reporter *reporter.Client
}

var _ simplejson.Handler = &StatsHandler{}

func (s StatsHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: s.tableQuery}
}

func (s *StatsHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	batches, err := s.Reporter.Vaccines.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}

	var vaccinations *data.Table
	vaccinations, err = s.Reporter.Vaccinations.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccinations call failed: %w", err)
	}

	d := calculateVaccineReserve(vaccinations, batches)

	return d.Filter(req.Args).CreateTableResponse(), nil
}

func calculateVaccineReserve(vaccinationsData, batches *data.Table) (output *data.Table) {
	// sum up the vaccinations and accumulate all vaccinations and batches
	output = sumVaccinations(vaccinationsData).Accumulate()
	batches = batches.Accumulate()

	// add the total received vaccines for each vaccinations entry
	addReserve(output, batches)

	return
}

func sumVaccinations(input *data.Table) (output *data.Table) {
	vaccinations := make([]float64, input.Frame.Rows())
	for i := 0; i < input.Frame.Rows(); i++ {
		var total float64
		for _, f := range input.Frame.Fields {
			if f.Name == "time" {
				continue
			}
			value := f.At(i).(float64)
			total += value
		}
		vaccinations[i] = total
	}

	return data.New(data.Column{Name: "time", Values: input.GetTimestamps()}, data.Column{Name: "vaccinations", Values: vaccinations})
}

func addReserve(vaccinationsData, batches *data.Table) {
	reserve := make([]float64, len(vaccinationsData.GetTimestamps()))

	batchIndex := -1
	var currentBatches float64
	var currentBatchTimestamp time.Time

	for r := 0; r < vaccinationsData.Frame.Rows(); r++ {
		v, _ := vaccinationsData.Frame.FloatAt(1, r)

		for batchIndex < batches.Frame.Rows()-1 && currentBatchTimestamp.Before(vaccinationsData.Frame.At(0, r).(time.Time)) {
			batchIndex++
			currentBatchTimestamp = batches.Frame.At(0, batchIndex).(time.Time)
			currentBatches = batches.Frame.At(1, batchIndex).(float64)
		}

		reserve[r] = currentBatches - v
	}

	vaccinationsData.Frame.Fields = append(vaccinationsData.Frame.Fields,
		grafanaData.NewField("reserve", nil, reserve),
	)
}
