package cases

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson"
)

// Handler returns the number of new COVID-19 cases. Use Scope to report cases by province, region or age group.
type Handler struct {
	reporter.Reporter
	Scope
}

var _ simplejson.Handler = &Handler{}

type Scope int

const (
	ScopeAll = iota
	ScopeProvince
	ScopeRegion
	ScopeAge
)

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		TableQuery: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	var entries *datasets.Dataset
	switch handler.Scope {
	case ScopeAll:
		entries, err = handler.Reporter.GetCases()
	case ScopeProvince:
		entries, err = handler.Reporter.GetCasesByProvince()
	case ScopeRegion:
		entries, err = handler.Reporter.GetCasesByRegion()
	case ScopeAge:
		entries, err = handler.Reporter.GetCasesByAgeGroup()
	}

	if err != nil {
		return nil, fmt.Errorf("tableQuery for cases failed: %w", err)
	}

	return responder.GenerateTableQueryResponse(entries, args), nil
}
