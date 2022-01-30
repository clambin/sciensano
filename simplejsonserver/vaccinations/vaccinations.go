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
	vaccinationData, err = handler.Reporter.GetVaccinations()

	if err != nil {
		return nil, fmt.Errorf("vaccinations failed: %w", err)
	}

	vaccinationData.Accumulate()
	for index := range vaccinationData.Groups[1].Values {
		vaccinationData.Groups[1].Values[index] += vaccinationData.Groups[2].Values[index]
	}
	vaccinationData.Groups = append(vaccinationData.Groups[0:2], vaccinationData.Groups[3])

	return responder.GenerateTableQueryResponse(vaccinationData, req.Args), nil
}
