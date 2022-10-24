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

	timestamps := vaccinationData.GetTimestamps()
	partial, _ := vaccinationData.GetFloatValues("partial")
	second, _ := vaccinationData.GetFloatValues("full")
	singledose, _ := vaccinationData.GetFloatValues("singledose")
	//booster, _ := vaccinationData.GetFloatValues("booster")

	if len(second) != len(singledose) {
		panic("data for second dose & single full-dose should be the same")
	}

	full := make([]float64, len(second))
	for i := range second {
		full[i] = second[i] + singledose[i]
	}

	d := data.New(
		data.Column{Name: "time", Values: timestamps},
		//data.Column{Name: "booster", Values: booster},
		data.Column{Name: "full", Values: full},
		data.Column{Name: "partial", Values: partial})

	return d.Accumulate().Filter(req.Args).CreateTableResponse(), nil
}
