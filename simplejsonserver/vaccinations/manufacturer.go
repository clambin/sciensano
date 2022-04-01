package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/dataset"
	"github.com/clambin/simplejson/v3/query"
)

// ManufacturerHandler returns COVID-19 vaccinations grouped by manufacturer
type ManufacturerHandler struct {
	reporter.Reporter
}

var _ simplejson.Handler = &ManufacturerHandler{}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler ManufacturerHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *ManufacturerHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var vaccinationData *dataset.Dataset
	vaccinationData, err = handler.Reporter.GetVaccinationsByManufacturer()

	if err != nil {
		return nil, fmt.Errorf("vaccinations by manufacturer failed: %w", err)
	}

	vaccinationData.Accumulate()
	return responder.GenerateTableQueryResponse(vaccinationData, req.Args), nil
}
