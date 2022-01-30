package testresults

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the COVID-19 test results
type Handler struct {
	reporter.Reporter
}

func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		Query: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (response query.Response, err error) {
	var tests *datasets.Dataset
	tests, err = handler.Reporter.GetTestResults()

	if err != nil {
		return nil, fmt.Errorf("testresults failed: %w", err)
	}

	return responder.GenerateTableQueryResponse(tests, req.Args), nil
}
