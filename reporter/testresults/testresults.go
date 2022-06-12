package testresults

import (
	"fmt"
	apiCache "github.com/clambin/sciensano/apiclient/cache"
	reportCache "github.com/clambin/sciensano/reporter/cache"
	"github.com/clambin/sciensano/reporter/table"
	"github.com/clambin/simplejson/v3/data"
	grafanaData "github.com/grafana/grafana-plugin-sdk-go/data"
)

type Reporter struct {
	ReportCache *reportCache.Cache
	APICache    apiCache.Holder
}

// Get returns all COVID-19 test results up to endTime
func (r *Reporter) Get() (results *data.Table, err error) {
	return r.ReportCache.MaybeGenerate("TestResults", func() (output *data.Table, err2 error) {
		if apiResult, found := r.APICache.Get("TestResults"); found {
			output = table.NewFromAPIResponse(apiResult)
			calculateRate(output)
		} else {
			err2 = fmt.Errorf("cache does not contain TestResults entries")
		}
		return
	})
}

func calculateRate(d *data.Table) {
	totalField, totalIdx := d.Frame.FieldByName("total")
	if totalIdx == -1 {
		panic("could not find `total` field in test results")
	}
	positiveField, positiveIdx := d.Frame.FieldByName("positive")
	if positiveIdx == -1 {
		panic("could not find `positive` field in test results")
	}

	var rates []float64

	for i := 0; i < totalField.Len(); i++ {
		var rate float64
		if totalField.At(i).(float64) != 0 {
			rate = positiveField.At(i).(float64) / totalField.At(i).(float64)
		}
		rates = append(rates, rate)
	}

	d.Frame.Fields = append(d.Frame.Fields, grafanaData.NewField("rate", nil, rates))
}
