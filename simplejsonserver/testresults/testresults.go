package testresults

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/sciensano/simplejsonserver/responder"
	"github.com/clambin/simplejson"
)

// Handler returns the COVID-19 test results
type Handler struct {
	reporter.Reporter
}

func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		TableQuery: handler.tableQuery,
	}
}

func (handler *Handler) tableQuery(_ context.Context, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	var tests *datasets.Dataset
	tests, err = handler.Reporter.GetTestResults()

	if err != nil {
		return nil, fmt.Errorf("testresults failed: %w", err)
	}

	return responder.GenerateTableQueryResponse(tests, args), nil
}
