package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type ManufacturerHandler struct {
	Reporter *reporter.Client
}

var _ simplejson.Handler = &ManufacturerHandler{}

func (m ManufacturerHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: m.tableQuery}
}

func (m *ManufacturerHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	batches, err := m.Reporter.Vaccines.GetByManufacturer()
	if err != nil {
		return nil, fmt.Errorf("vaccine manufacturer call failed: %w", err)
	}
	return batches.Accumulate().Filter(req.Args).CreateTableResponse(), nil
}
