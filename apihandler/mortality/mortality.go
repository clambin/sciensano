package mortality

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

// Handler implements a grafana-json handler for COVID-19 cases
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
		"mortality":        {TableQueryFunc: handler.buildCasesResponse},
		"mortality-region": {TableQueryFunc: handler.buildCasesResponse},
		"mortality-age":    {TableQueryFunc: handler.buildCasesResponse},
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

func (handler *Handler) buildCasesResponse(ctx context.Context, target string, args *grafanajson.TableQueryArgs) (output *grafanajson.TableQueryResponse, err error) {
	var cases *datasets.Dataset
	switch target {
	case "mortality":
		cases, err = handler.Sciensano.GetMortality(ctx)
	case "mortality-region":
		cases, err = handler.Sciensano.GetMortalityByRegion(ctx)
	case "mortality-age":
		cases, err = handler.Sciensano.GetMortalityByAgeGroup(ctx)
	}

	if err == nil {
		output = response.GenerateTableQueryResponse(cases, args)
	}

	return
}
