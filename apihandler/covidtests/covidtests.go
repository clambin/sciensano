package covidtests

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/sciensano"
	log "github.com/sirupsen/logrus"
	"time"
)

type Handler struct {
	Sciensano   sciensano.APIClient
	targetTable grafanajson.TargetTable
}

func New(client sciensano.APIClient) (handler *Handler) {
	handler = &Handler{
		Sciensano: client,
	}

	handler.targetTable = grafanajson.TargetTable{
		"tests": {TableQueryFunc: handler.buildTestTableResponse},
	}

	return
}

func (handler *Handler) Endpoints() grafanajson.Endpoints {
	return grafanajson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

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
	var tests []sciensano.TestResult
	tests, err = handler.Sciensano.GetTests(ctx, args.Range.To)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve test results: %s", err.Error())
	}

	rows := len(tests)
	timestamps := make(grafanajson.TableQueryResponseTimeColumn, rows)
	allTests := make(grafanajson.TableQueryResponseNumberColumn, rows)
	positiveTests := make(grafanajson.TableQueryResponseNumberColumn, rows)
	positiveRate := make(grafanajson.TableQueryResponseNumberColumn, rows)

	for index, test := range tests {
		timestamps[index] = test.Timestamp
		allTests[index] = float64(test.Total)
		positiveTests[index] = float64(test.Positive)
		positiveRate[index] = float64(test.Positive) / float64(test.Total)
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
