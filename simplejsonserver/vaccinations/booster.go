package vaccinations

import (
	"context"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

type BoosterHandler struct {
	Reporter *reporter.Client
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (b *BoosterHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: func(_ context.Context, req query.Request) (query.Response, error) {
		results, err := b.Reporter.Vaccinations.Get()
		if err != nil {
			return nil, err
		}

		return results.DeleteColumn("partial", "full", "singledose").Accumulate().Filter(req.Args).CreateTableResponse(), nil
	}}
}
