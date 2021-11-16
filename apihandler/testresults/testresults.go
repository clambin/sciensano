package testresults

import (
	"context"
	"fmt"
	grafanajson "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/apihandler/response"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 test results
type Handler struct {
	Sciensano   reporter.Reporter
	targetTable grafanajson.TargetTable
}

// New creates a new Handler
func New(client reporter.Reporter) (handler *Handler) {
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
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Info("TableQuery called")
	return
}

func (handler *Handler) buildTestTableResponse(_ context.Context, _ string, args *grafanajson.TableQueryArgs) (output *grafanajson.TableQueryResponse, err error) {
	var tests *datasets.Dataset
	tests, err = handler.Sciensano.GetTestResults()

	if err == nil {
		output = response.GenerateTableQueryResponse(tests, args)

		// TODO: calculate ratio in reporter?
		positiveRate := make(grafanajson.TableQueryResponseNumberColumn, len(tests.Timestamps))

		for index := range tests.Timestamps {
			positiveRate[index] = output.Columns[2].Data.(grafanajson.TableQueryResponseNumberColumn)[index] / output.Columns[1].Data.(grafanajson.TableQueryResponseNumberColumn)[index]
		}

		output.Columns = append(output.Columns, grafanajson.TableQueryResponseColumn{
			Text: "rate",
			Data: positiveRate,
		})
	}

	return
}
