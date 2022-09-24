package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
)

// GroupedHandler returns COVID-19 vaccinations grouped by region or age group, for a specific type (i.e. partial, full or booster vaccination)
type GroupedHandler struct {
	Reporter   *reporter.Client
	Type       int
	Accumulate bool
	Scope
}

var _ simplejson.Handler = &GroupedHandler{}

type Scope int

const (
	ScopeRegion = iota
	ScopeAge
)

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *GroupedHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{Query: handler.tableQuery}
}

func (handler *GroupedHandler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var vaccinationData *data.Table
	switch handler.Scope {
	case ScopeAge:
		vaccinationData, err = handler.Reporter.Vaccinations.GetByAgeGroup(handler.Type)
	case ScopeRegion:
		vaccinationData, err = handler.Reporter.Vaccinations.GetByRegion(handler.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("grouped vaccinations failed: %w", err)
	}

	if handler.Accumulate {
		vaccinationData = vaccinationData.Accumulate()
	}

	return vaccinationData.Filter(req.Args).CreateTableResponse(), nil
}
