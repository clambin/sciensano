package mortality

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson"
)

// Handler returns the number of deaths. Use Scope to report by Region or Age
type Handler struct {
	reporter.Reporter
	Scope
}

type Scope int

const (
	ScopeAll = iota
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
		entries, err = handler.Reporter.GetMortality()
	case ScopeRegion:
		entries, err = handler.Reporter.GetMortalityByRegion()
	case ScopeAge:
		entries, err = handler.Reporter.GetMortalityByAgeGroup()
	}

	if err != nil {
		return nil, fmt.Errorf("mortality failed: %w", err)
	}

	return responder.GenerateTableQueryResponse(entries, args), nil
}
