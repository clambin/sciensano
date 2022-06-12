package mortality

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the number of deaths. Use Scope to report by Region or Age
type Handler struct {
	Reporter *reporter.Client
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
		Query: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var entries *data.Table
	switch handler.Scope {
	case ScopeAll:
		entries, err = handler.Reporter.Mortality.Get()
	case ScopeRegion:
		entries, err = handler.Reporter.Mortality.GetByRegion()
	case ScopeAge:
		entries, err = handler.Reporter.Mortality.GetByAgeGroup()
	}

	if err != nil {
		return nil, fmt.Errorf("mortality failed: %w", err)
	}

	return entries.Filter(req.Args).CreateTableResponse(), nil
}
