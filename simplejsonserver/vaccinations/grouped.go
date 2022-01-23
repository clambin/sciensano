package vaccinations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson"
)

// GroupedHandler returns COVID-19 vaccinations grouped by region or age group, for a specific type (i.e. partial, full or booster vaccination)
type GroupedHandler struct {
	reporter.Reporter
	reporter.VaccinationType
	Scope
}

var _ simplejson.Handler = &GroupedHandler{}

type Scope int

const (
	ScopeRegion = iota
	ScopeAge
)

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler GroupedHandler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{TableQuery: handler.tableQuery}
}

func (handler *GroupedHandler) tableQuery(_ context.Context, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	var vaccinationData *datasets.Dataset
	switch handler.Scope {
	case ScopeAge:
		vaccinationData, err = handler.Reporter.GetVaccinationsByAgeGroup(handler.VaccinationType)
	case ScopeRegion:
		vaccinationData, err = handler.Reporter.GetVaccinationsByRegion(handler.VaccinationType)
	}

	if err != nil {
		return nil, fmt.Errorf("grouped vaccinations failed: %w", err)
	}

	vaccinationData.Accumulate()
	return responder.GenerateTableQueryResponse(vaccinationData, args), nil
}
