package hospitalisations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/data"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the number of hospitalisations. Use Scope to report by Region or Province.
type Handler struct {
	Reporter *reporter.Client
	Scope
}

type Scope int

const (
	ScopeAll = iota
	ScopeRegion
	ScopeProvince
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
		entries, err = handler.Reporter.Hospitalisations.Get()
	case ScopeRegion:
		entries, err = handler.Reporter.Hospitalisations.GetByRegion()
	case ScopeProvince:
		entries, err = handler.Reporter.Hospitalisations.GetByProvince()
	}

	if err != nil {
		return nil, fmt.Errorf("hospitalisations failed: %w", err)
	}

	return entries.Filter(req.Args).CreateTableResponse(), nil
}
