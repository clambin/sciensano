package vaccinations

import (
	"context"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
	"strings"
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

		// FIXME: would be easier if data.Table had a "DeleteColumn" method
		columns := []data.Column{{Name: "time", Values: results.GetTimestamps()}}
		for _, c := range results.GetColumns() {
			if strings.HasPrefix(c, "booster") {
				values, _ := results.GetFloatValues(c)
				columns = append(columns, data.Column{Name: c, Values: values})
			}
		}

		return data.New(columns...).Accumulate().Filter(req.Args).CreateTableResponse(), nil
	}}
}
