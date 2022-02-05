package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the overall COVID-19 vaccinations
type Handler struct {
	reporter.Reporter
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (output query.Response, err error) {
	var vaccinationData *datasets.Dataset
	if vaccinationData, err = handler.Reporter.GetVaccinations(); err != nil {
		return nil, fmt.Errorf("vaccinations failed: %w", err)
	}

	timestamps := vaccinationData.GetTimestamps()
	partial, _ := vaccinationData.GetValues("partial")
	full, _ := vaccinationData.GetValues("full")
	singledose, _ := vaccinationData.GetValues("singledose")
	booster, _ := vaccinationData.GetValues("booster")

	d := datasets.New()
	for index, timestamp := range timestamps {
		d.Add(timestamp, "partial", partial[index])
		d.Add(timestamp, "full", full[index]+singledose[index])
		d.Add(timestamp, "booster", booster[index])
	}
	d.Accumulate()

	return responder.GenerateTableQueryResponse(d, req.Args), nil
}
