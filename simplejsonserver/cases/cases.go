package cases

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the number of new COVID-19 cases. Use Scope to report cases by province, region or age group.
type Handler struct {
	Reporter *reporter.Client
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
		Query: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var entries *data.Table
	switch handler.Scope {
	case ScopeAll:
		entries, err = handler.Reporter.Cases.Get()
	case ScopeProvince:
		entries, err = handler.Reporter.Cases.GetByProvince()
	case ScopeRegion:
		entries, err = handler.Reporter.Cases.GetByRegion()
	case ScopeAge:
		entries, err = handler.Reporter.Cases.GetByAgeGroup()
	}

	if err != nil {
		return nil, fmt.Errorf("tableQuery for cases failed: %w", err)
	}

	return entries.Filter(req.Args).CreateTableResponse(), nil
}
