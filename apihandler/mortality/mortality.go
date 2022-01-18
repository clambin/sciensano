package mortality

import (
	"context"
	"fmt"
	"github.com/clambin/sciensano/apihandler/response"
	"github.com/clambin/sciensano/reporter"
	"github.com/clambin/sciensano/reporter/datasets"
	"github.com/clambin/simplejson"
	log "github.com/sirupsen/logrus"
	"time"
)

// Handler implements a grafana-json handler for COVID-19 cases
type Handler struct {
	Sciensano   reporter.Reporter
	targetTable simplejson.TargetTable
}

// New creates a new Handler
func New(client reporter.Reporter) (handler *Handler) {
	handler = &Handler{
		Sciensano: client,
	}

	handler.targetTable = simplejson.TargetTable{
		"mortality":        {TableQueryFunc: handler.buildCasesResponse},
		"mortality-region": {TableQueryFunc: handler.buildCasesResponse},
		"mortality-age":    {TableQueryFunc: handler.buildCasesResponse},
	}

	return
}

// Endpoints implements the grafana-json Endpoint function. It returns all supported endpoints
func (handler *Handler) Endpoints() simplejson.Endpoints {
	return simplejson.Endpoints{
		Search:     handler.Search,
		TableQuery: handler.TableQuery,
	}
}

// Search implements the grafana-json Search function. It returns all supported targets
func (handler *Handler) Search() (targets []string) {
	return handler.targetTable.Targets()
}

// TableQuery implements the grafana-json TableQuery function. It processes incoming TableQuery requests
func (handler *Handler) TableQuery(ctx context.Context, target string, args *simplejson.TableQueryArgs) (response *simplejson.TableQueryResponse, err error) {
	start := time.Now()
	response, err = handler.targetTable.RunTableQuery(ctx, target, args)
	if err != nil {
		return nil, fmt.Errorf("unable to build table query response for target '%s': %s", target, err.Error())
	}
	log.WithFields(log.Fields{"duration": time.Now().Sub(start), "target": target}).Info("TableQuery called")
	return
}

func (handler *Handler) buildCasesResponse(_ context.Context, target string, args *simplejson.TableQueryArgs) (output *simplejson.TableQueryResponse, err error) {
	var cases *datasets.Dataset
	switch target {
	case "mortality":
		cases, err = handler.Sciensano.GetMortality()
	case "mortality-region":
		cases, err = handler.Sciensano.GetMortalityByRegion()
	case "mortality-age":
		cases, err = handler.Sciensano.GetMortalityByAgeGroup()
	}

	if err == nil {
		output = response.GenerateTableQueryResponse(cases, args)
	}

	return
}
