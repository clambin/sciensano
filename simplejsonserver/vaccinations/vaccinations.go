package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the overall COVID-19 vaccinations
type Handler struct {
	Reporter *reporter.Client
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	vaccinationData, err := handler.Reporter.Vaccinations.Get()
	if err != nil {
		return nil, fmt.Errorf("vaccinations failed: %w", err)
	}

	partials, _ := vaccinationData.GetFloatValues("partial")
	full, _ := vaccinationData.GetFloatValues("full")
	singledose, _ := vaccinationData.GetFloatValues("singledose")

	if len(full) != len(singledose) {
		panic("data for second dose & single full-dose should be the same")
	}

	for i := range full {
		full[i] += singledose[i]
	}

	d := data.New(
		data.Column{Name: "time", Values: vaccinationData.GetTimestamps()},
		data.Column{Name: "full", Values: full},
		data.Column{Name: "partial", Values: partials})

	return d.Accumulate().Filter(req.Args).CreateTableResponse(), nil
}
