package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// ManufacturerHandler returns COVID-19 vaccinations grouped by manufacturer
type ManufacturerHandler struct {
	Reporter *reporter.Client
}

var _ simplejson.Handler = &ManufacturerHandler{}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *ManufacturerHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *ManufacturerHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	vaccinationData, err := handler.Reporter.Vaccinations.GetByManufacturer()

	if err != nil {
		return nil, fmt.Errorf("vaccinations by manufacturer failed: %w", err)
	}

	return vaccinationData.Accumulate().Filter(req.Args).CreateTableResponse(), nil
}
