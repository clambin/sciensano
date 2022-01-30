package hospitalisations

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the number of hospitalisations. Use Scope to report by Region or Province.
type Handler struct {
	reporter.Reporter
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
	var entries *datasets.Dataset
	switch handler.Scope {
	case ScopeAll:
		entries, err = handler.Reporter.GetHospitalisations()
	case ScopeRegion:
		entries, err = handler.Reporter.GetHospitalisationsByRegion()
	case ScopeProvince:
		entries, err = handler.Reporter.GetHospitalisationsByProvince()
	}

	if err != nil {
		return nil, fmt.Errorf("hospitalisations failed: %w", err)
	}

	return responder.GenerateTableQueryResponse(entries, req.Args), nil
}
