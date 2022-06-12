package testresults

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/simplejson/v3"
	"github.com/clambin/simplejson/v3/query"
)

// Handler returns the COVID-19 test results
type Handler struct {
	Reporter *reporter.Client
}

func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		Query: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, req query.Request) (query.Response, error) {
	tests, err := handler.Reporter.TestResults.Get()

	if err != nil {
		return nil, fmt.Errorf("testresults failed: %w", err)
	}

	return tests.Filter(req.Args).CreateTableResponse(), nil
}
