package vaccines

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type OverviewHandler struct {
	Reporter *reporter.Client
}

var _ simplejson.Handler = &OverviewHandler{}

func (o OverviewHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: o.tableQuery}
}

func (o *OverviewHandler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	batches, err := o.Reporter.Vaccines.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccine call failed: %w", err)
	}
	return batches.Accumulate().Filter(req.Args).CreateTableResponse(), nil
}
