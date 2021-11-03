package testresults

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	"github.com/clambin/sciensano/sciensano/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 test results
type Handler struct {
	Sciensano   sciensano.APIClient
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(client sciensano.APIClient) (handler *Handler) {
	handler = &Handler{
		Sciensano: client,
	}

	handler.targetTable = grafanajson.TargetTable{
		"tests": {TableQueryFunc: handler.buildTestTableResponse},
	}

	return
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() grafanajson.Endpoints {
	return grafanajson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

// Search implements the grafana-json Search function. It returns all supported targets
func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

// TableQuery implements the grafana-json TableQuery function. It processes incoming TableQuery requests
func (handler *Handler) TableQuery(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Debug("TableQuery called")
	return
}

func (handler *Handler) buildTestTableResponse(ctx context.Context, _ string, args *grafanajson.TableQueryArgs) (response *grafanajson.TableQueryResponse, err error) {
	var tests *datasets.Dataset
	tests, err = handler.Sciensano.GetTestResults(ctx)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve test results: %s", err.Error())
	}

	tests.ApplyRange(args.Range.From, args.Range.To)
	rows := len(tests.Timestamps)
	timestamps := make(grafanajson.TableQueryResponseTimeColumn, rows)
	allTests := make(grafanajson.TableQueryResponseNumberColumn, rows)
	positiveTests := make(grafanajson.TableQueryResponseNumberColumn, rows)
	positiveRate := make(grafanajson.TableQueryResponseNumberColumn, rows)

	for index, timestamp := range tests.Timestamps {
		timestamps[index] = timestamp
		entry := tests.Groups[0].Values[index].(*sciensano.TestResult)
		allTests[index] = float64(entry.Total)
		positiveTests[index] = float64(entry.Positive)
		positiveRate[index] = entry.Ratio()
	}

	response = new(grafanajson.TableQueryResponse)
	response.Columns = []grafanajson.TableQueryResponseColumn{
		{Text: "timestamp", Data: timestamps},
		{Text: "total", Data: allTests},
		{Text: "positive", Data: positiveTests},
		{Text: "rate", Data: positiveRate},
	}

	return
}
